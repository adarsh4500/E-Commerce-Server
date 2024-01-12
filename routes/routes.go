package routes

import (
	"Ecom/controllers"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func SetupRouter() *gin.Engine {

	router := gin.Default()

	authMiddleware := controllers.Authenticate

	router.POST("/signup", controllers.SignupHandler)
	router.POST("/login", controllers.LoginHandler)

	productGroup := router.Group("/products", authMiddleware)
	{
		productGroup.GET("/:id", controllers.GetProductHandler)
		productGroup.GET("", controllers.GetAllProductsHandler)
		productGroup.POST("/new", controllers.NewProductHandler)
		productGroup.POST("/update/:id", controllers.UpdateProductHandler)
		productGroup.POST("/delete/:id", controllers.DeleteProductHandler)
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
		orderGroup.POST("/updatestatus", controllers.UpdateOrderStatusHandler)
	}

	return router
}
