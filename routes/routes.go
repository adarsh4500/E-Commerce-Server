package routes

import (
	"Ecom/controllers"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func SetupRouter() *gin.Engine {

	router := gin.Default()

	authMiddleware := controllers.Authenticate
	adminMiddleware := controllers.RequireAdmin

	router.POST("/signup", controllers.SignupHandler)
	router.POST("/login", controllers.LoginHandler)
	router.POST("/logout", controllers.LogOutHandler)

	productGroup := router.Group("/products", authMiddleware)
	{
		productGroup.GET("/:id", controllers.GetProductHandler)
		productGroup.GET("", controllers.GetAllProductsHandler)
		productGroup.POST("/new", adminMiddleware, controllers.NewProductHandler)
		productGroup.POST("/update/:id", adminMiddleware, controllers.UpdateProductHandler)
		productGroup.POST("/delete/:id", adminMiddleware, controllers.DeleteProductHandler)
	}

	cartGroup := router.Group("/cart", authMiddleware)
	{
		cartGroup.POST("/new", controllers.AddToCartHandler)
		cartGroup.POST("/remove/:id", controllers.RemoveFromCartHandler)
		cartGroup.GET("", controllers.ViewCartHandler)
		cartGroup.POST("/clear", controllers.ClearCartHandler)
		cartGroup.POST("/checkout", controllers.PlaceOrderHandler)
	}

	orderGroup := router.Group("/orders", authMiddleware)
	{
		orderGroup.GET("", controllers.ViewOrderHandler)
		orderGroup.GET("/:id", controllers.ViewOrderItemsHandler)
		orderGroup.POST("/updatestatus", adminMiddleware, controllers.UpdateOrderStatusHandler)
	}

	return router
}
