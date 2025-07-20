package controllers

import (
	"Ecom/config"
	"Ecom/connections"
	"Ecom/helpers"
	"Ecom/models"
	"Ecom/postgres"
	"Ecom/utils"
	"context"
	"database/sql"
	"errors"
	"net/http"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"reflect"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func LoginHandler(c *gin.Context) {
	var creds models.Creds
	err := c.ShouldBindJSON(&creds)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	// Input validation
	if !isValidEmail(creds.Email) {
		utils.ErrorResponse(c, http.StatusBadRequest, errors.New("invalid email format"))
		return
	}
	if len(creds.Password) < 8 {
		utils.ErrorResponse(c, http.StatusBadRequest, errors.New("password must be at least 8 characters"))
		return
	}

	query := postgres.New(connections.DB)

	user, err := query.GetUserByEmail(context.Background(), creds.Email)
	if err == sql.ErrNoRows {
		utils.ErrorResponse(c, http.StatusUnauthorized, errors.New("user not found"))
		return
	} else if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	if helpers.IsPasswordCorrect(creds.Password, user.Password) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"email": creds.Email,
			"id":    user.ID.String(),
			"role":  user.Role,
			"exp":   time.Now().Add(1 * time.Hour).Unix(),
		})

		tokenString, err := token.SignedString([]byte(config.JWTSecret))
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, errors.New("invalid email or password"))
			return
		}

		// Set cookie with security flags
		c.SetSameSite(http.SameSiteNoneMode)
		secure := true 
		c.SetCookie("token", tokenString, 3600, "", "", secure, true)
		utils.SuccessResponseWithMessage(c, "Successfully Logged in")
		return
	}

	utils.ErrorResponse(c, http.StatusUnauthorized, errors.New("incorrect password"))
}

func SignupHandler(c *gin.Context) {
	var user models.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	// Input validation
	if !isValidEmail(user.Email) {
		utils.ErrorResponse(c, http.StatusBadRequest, errors.New("invalid email format"))
		return
	}
	if len(user.Password) < 8 {
		utils.ErrorResponse(c, http.StatusBadRequest, errors.New("password must be at least 8 characters"))
		return
	}

	user.Password, err = helpers.Encrypt(user.Password)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	newUser := postgres.AddUserParams{
		Fullname: user.Fullname,
		Email:    user.Email,
		Password: user.Password,
		Role:     models.RoleUser,
	}

	query := postgres.New(connections.DB)
	err = query.AddUser(context.Background(), newUser)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	utils.SuccessResponseWithMessage(c, "Signed Up Successfully")
}

func LogOutHandler(c *gin.Context) {
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("token", "", -1, "/", "", true, true)
	utils.SuccessResponseWithMessage(c, "Successfully Logged Out")
}

func MeHandler(c *gin.Context) {
	claims, err := parseJWTToken(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, err)
		return
	}
	email, emailOk := claims["email"].(string)
	if !emailOk {
		utils.ErrorResponse(c, http.StatusUnauthorized, errors.New("invalid jwt claims"))
		return
	}
	query := postgres.New(connections.DB)
	user, err := query.GetUserByEmail(context.Background(), email)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	type userResponse struct {
		Fullname string      `json:"fullname"`
		Email    string      `json:"email"`
		Role     string      `json:"role"`
		ID       interface{} `json:"id,omitempty"`
	}
	resp := userResponse{
		Fullname: user.Fullname,
		Email:    user.Email,
		Role:     user.Role,
	}
	if idField := getField(user, "ID"); idField != nil {
		resp.ID = idField
	}
	utils.SuccessResponse(c, resp)
}

func GetAllUsersHandler(c *gin.Context) {
	query := postgres.New(connections.DB)
	users, err := query.GetUsers(context.Background())
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	type userResponse struct {
		Fullname string      `json:"fullname"`
		Email    string      `json:"email"`
		Role     string      `json:"role"`
		ID       interface{} `json:"id,omitempty"`
	}
	resp := make([]userResponse, 0, len(users))
	for _, user := range users {
		item := userResponse{
			Fullname: user.Fullname,
			Email:    user.Email,
			Role:     user.Role,
		}
		if idField := getField(user, "ID"); idField != nil {
			item.ID = idField
		}
		resp = append(resp, item)
	}
	utils.SuccessResponse(c, resp)
}

// getField is a helper to get the ID field if present (for future-proofing)
func getField(user interface{}, field string) interface{} {
	// Use reflection to get the field if it exists
	// This is a safe fallback if the user struct has an ID field
	v := reflect.ValueOf(user)
	if v.Kind() == reflect.Struct {
		f := v.FieldByName(field)
		if f.IsValid() {
			return f.Interface()
		}
	}
	return nil
}

// Email validation helper
func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
