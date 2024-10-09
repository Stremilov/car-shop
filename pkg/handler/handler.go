package handler

import (
	_ "github.com/Stremilov/car-shop/docs"
	"github.com/Stremilov/car-shop/pkg/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutesAndDB() *gin.Engine {
	initDB()
	router := gin.New()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	{
		users := api.Group("/user")
		{
			users.POST("/", h.addUser)
			users.GET("/get-all", h.getAllUsers)
			users.GET("/:userID", h.getUserByID)
			users.PATCH("/:userID", h.updateUserInfoByID)
			users.DELETE("/:userID", h.deleteUserByID)
		}

		cars := api.Group("/car")
		{
			cars.POST("/", h.addCar)
			cars.GET("/:carID", h.getCarByID)
			cars.GET("/get-all", h.getAllCars)
			cars.PATCH(":carID", h.updateCarInfoByID)
			cars.DELETE("/:carID", h.deleteCarByID)
		}

		orders := api.Group("/orders")
		{
			orders.POST("/", h.createOrder)
			orders.GET("/get-all", h.getAllOrders)
			orders.GET("/:userID", h.getOrdersByUserID)
			orders.DELETE("/:orderID", h.deleteOrderByID)
		}
	}

	return router
}
