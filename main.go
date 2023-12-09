package main

import (
	"Ecom/postgres"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var secret_Key = "hsgvhgy337gh37bhd37GD787rn4HF7348rG3HB784V34dskjkvk3tgIA78BBbgF74GGF"

const (
	Success      = 1
	Failed       = 0
	UserNotExist = 2
)

func helloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, "Hello World")
}

func loginHandler(c *gin.Context) {
	var creds postgres.ValidateCredsParams
	err := c.BindJSON(&creds)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	conn, err := sql.Open("postgres", "user=postgres password=alpha dbname=Ecom sslmode=disable")
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	db := postgres.New(conn)

	authResult, err := db.ValidateCreds(context.Background(), creds)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	if authResult == Failed {
		c.JSON(http.StatusUnauthorized, "Wrong Password")
		return
	}
	if authResult == UserNotExist {
		c.JSON(http.StatusUnauthorized, "User Not Found")
		return
	}

	//Create Token Object
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":       creds.Email,
		"expires at": time.Now().Add(1 * time.Hour).Unix(),
	})

	//Sign Token
	tokenString, err := token.SignedString([]byte(secret_Key))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Email or Password"})
	}

	//Set Cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("token", tokenString, int(time.Now().Add(1*time.Hour).Unix()), "", "", false, true)
	c.JSON(http.StatusOK, "Successfully Logged in")
}

func authenticate(c *gin.Context) {
	//Get Token off Cookie
	tokenString, err := c.Cookie("token")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	//Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		return []byte(secret_Key), nil
	})
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user := claims["email"]
		expiresAt := claims["expires at"].(float64)
		if float64(time.Now().Unix()) < expiresAt {
			c.Set("user", user)
			c.Next()
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}

func logOutHandler(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, "Logged Out")
}

func GetAllProductsHandler(c *gin.Context) {
	conn, err := sql.Open("postgres", "user=postgres password=alpha dbname=Ecom sslmode=disable")
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	db := postgres.New(conn)

	products, err := db.GetProducts(context.Background())
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, products)
}

func GetProductHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	conn, err := sql.Open("postgres", "user=postgres password=alpha dbname=Ecom sslmode=disable")
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	db := postgres.New(conn)
	product, err := db.GetProductById(context.Background(), int32(id))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, product)
}

func AddProductHandler(c *gin.Context) {
	var newProduct postgres.AddProductParams
	err := c.BindJSON(&newProduct)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	conn, err := sql.Open("postgres", "user=postgres password=alpha dbname=Ecom sslmode=disable")
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	db := postgres.New(conn)

	product, err := db.AddProduct(context.Background(), newProduct)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, product)
}

func main() {

	router := gin.Default()

	router.POST("/login", loginHandler)
	router.GET("/hello", authenticate, helloHandler)
	router.GET("/logout", authenticate, logOutHandler)
	router.GET("/products/:id", authenticate, GetProductHandler)
	router.GET("/products", authenticate, GetAllProductsHandler)

	fmt.Println("Starting")
	router.Run(":8080")

}
