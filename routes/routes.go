package routes

import (
	"Ecom/config"
	"Ecom/controllers"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func SetupRouter() *gin.Engine {

	router := gin.Default()

	// Parse allowed origins from environment variable
	allowedOrigins := []string{"http://localhost:5173"}
	if config.AllowedOrigins != "" {
		allowedOrigins = strings.Split(config.AllowedOrigins, ",")
		for i := range allowedOrigins {
			allowedOrigins[i] = strings.TrimSpace(allowedOrigins[i])
		}
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60,
	}))

	authMiddleware := controllers.Authenticate
	adminMiddleware := controllers.RequireAdmin

	router.POST("/signup", controllers.SignupHandler)
	router.POST("/login", controllers.LoginHandler)
	router.POST("/logout", controllers.LogOutHandler)

	// Public product routes
	router.GET("/products/count", controllers.GetProductsCountHandler)
	router.GET("/products", controllers.GetAllProductsHandler)
	router.GET("/products/:id", controllers.GetProductHandler)

	// Admin product routes
	productGroup := router.Group("/products", authMiddleware)
	{
		productGroup.POST("/new", adminMiddleware, controllers.NewProductHandler)
		productGroup.POST("/update/:id", adminMiddleware, controllers.UpdateProductHandler)
		productGroup.POST("/delete/:id", adminMiddleware, controllers.DeleteProductHandler)
	}

	cartGroup := router.Group("/cart", authMiddleware)
	{
		cartGroup.POST("/new", controllers.AddToCartHandler)
		cartGroup.POST("/update/:product_id", controllers.UpdateCartItemHandler)
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

	// User endpoints
	userGroup := router.Group("/user", authMiddleware)
	{
		userGroup.GET("/me", controllers.MeHandler)
		userGroup.GET("/cart/count", controllers.CartCountHandler)
		userGroup.GET("/orders/history", controllers.OrderHistoryHandler)
	}

	// Admin endpoints
	adminGroup := router.Group("/admin", adminMiddleware)
	{
		adminGroup.GET("/users", controllers.GetAllUsersHandler)
		adminGroup.GET("/orders", controllers.AdminAllOrdersHandler)
	}

	return router
}
