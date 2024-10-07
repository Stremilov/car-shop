package handler

import "github.com/gin-gonic/gin"

type Handler struct {
}

func (h *Handler) InitRoutesAndDB() *gin.Engine {
	initDB()
	router := gin.New()

	api := router.Group("/api")
	{
		users := api.Group("/user")
		{
			users.POST("/", h.addUser)
			users.GET("/get-all", h.getAllUsers)
			users.GET("/:userID", h.getUserByID)
			users.PATCH("/:userID", h.updateUserInfoByID)
		}

		cars := api.Group("/car")
		{
			cars.POST("/", h.addCar)
			cars.GET("/:carID", h.getCarByID)
			cars.GET("/get-all", h.getAllCars)
		}

		orders := api.Group("/orders")
		{
			orders.POST("/", h.createOrder)
			orders.GET("/get-all", h.getAllOrders)
		}
	}

	return router
}
