package routes

import (
	"Ecom/controllers"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func SetupRouter() *gin.Engine {

	router := gin.Default()

	router.POST("/signup", controllers.SignupHandler)
	router.POST("/login", controllers.LoginHandler)
	router.GET("/logout", controllers.Authenticate, controllers.LogOutHandler)
	router.GET("/products/:id", controllers.Authenticate, controllers.GetProductHandler)
	router.GET("/products", controllers.Authenticate, controllers.GetAllProductsHandler)
	router.POST("/products/new", controllers.Authenticate, controllers.NewProductHandler)
	router.POST("/products/update/:id", controllers.Authenticate, controllers.UpdateProductHandler)
	router.POST("/products/delete/:id", controllers.Authenticate, controllers.DeleteProductHandler)

	return router
}
