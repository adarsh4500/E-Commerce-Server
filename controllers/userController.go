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
		// Create Token Object
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"email": creds.Email,
			"id":    user.ID.String(),
			"exp":   time.Now().Add(1 * time.Hour).Unix(),
		})

		// Sign Token
		tokenString, err := token.SignedString([]byte(config.JWTSecret))
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, errors.New("invalid email or password"))
			return
		}

		// Set Cookie (secure, httpOnly)
		c.SetSameSite(http.SameSiteLaxMode)
		secure := false
		if c.Request.TLS != nil {
			secure = true
		}
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
	c.SetCookie("token", "", -1, "/", "", false, true)
	utils.SuccessResponseWithMessage(c, "Successfully Logged Out")
}

// Email validation helper
func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
