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
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// @Summary Login
// @Description Logs in a user and returns an authentication token.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.Creds true "User credentials"
// @Success 200 {object} utils.TypeSuccessResponse
// @Failure 400 {object} utils.TypeErrorResponse
// @Failure 401 {object} utils.TypeErrorResponse
// @Failure 500 {object} utils.TypeErrorResponse
// @Router /login [post]
func LoginHandler(c *gin.Context) {
	var creds models.Creds
	err := c.ShouldBindJSON(&creds)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err)
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

		//Create Token Object
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user":       creds.Email,
			"id":         user.ID,
			"expires at": time.Now().Add(1 * time.Hour).Unix(),
		})

		//Sign Token
		tokenString, err := token.SignedString([]byte(config.JWTSecret))
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, errors.New("invalid email or password"))
		}

		//Set Cookie
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("token", tokenString, int(time.Now().Add(1*time.Hour).Unix()), "", "", false, true)
		models.UserID = user.ID
		utils.SuccessResponse(c, gin.H{"message": "Successfully Logged in"})
		return
	}

	utils.ErrorResponse(c, http.StatusUnauthorized, errors.New("incorrect password"))
}

// @Summary Signup
// @Description Registers a new user.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.User true "User information"
// @Success 200 {object} utils.TypeSuccessResponse
// @Failure 400 {object} utils.TypeErrorResponse
// @Failure 500 {object} utils.TypeErrorResponse
// @Router /signup [post]
func SignupHandler(c *gin.Context) {
	var user models.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err)
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
	utils.SuccessResponse(c, gin.H{"message": "Signed Up Successfully"})
}

// @Summary Logout
// @Description Logs out the currently authenticated user.
// @Tags Authentication
// @Produce json
// @Success 200 {object} utils.TypeSuccessResponse
// @Router /logout [post]
func LogOutHandler(c *gin.Context) {
	models.UserID = uuid.Nil
	c.SetCookie("token", "", -1, "/", "", false, true)
	utils.SuccessResponse(c, gin.H{"message": "Successflly Logged Out"})
}
