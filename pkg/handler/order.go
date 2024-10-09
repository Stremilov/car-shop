package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Order struct {
	OrderID   int    `json:"order_id"`
	OrderDate string `json:"order_date"`
	User      User   `json:"user"`
	Car       Car    `json:"car"`
}

type OrderByUserID struct {
	OrderID   int    `json:"order_id"`
	OrderDate string `json:"order_date"`
	User      User   `json:"user"`
	Car       []Car  `json:"car"`
}

type OrderInput struct {
	UserID int `json:"user_id"`
	CarID  int `json:"car_id"`
}

// @Summary      Create order
// @Description  create order
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param request body OrderInput true "body"
// @Success      201  {object}  OrderInput
// @Router       /api/orders/ [post]
func (h *Handler) createOrder(ctx *gin.Context) {
	var orderInput OrderInput

	if err := ctx.ShouldBindJSON(&orderInput); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "Invalid request"})
		return
	}

	query := `INSERT INTO orders (user_id, car_id) VALUES ($1, $2)`
	_, err := db.Exec(query, orderInput.UserID, orderInput.CarID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "Failed to insert into database"})

		return
	}

	ctx.JSON(http.StatusOK, orderInput)

}

// @Summary      Get all orders
// @Description  get all orders
// @Tags         orders
// @Accept       json
// @Produce      json
// @Success      200  {object}  Order
// @Router       /api/orders/get-all [get]
func (h *Handler) getAllOrders(ctx *gin.Context) {
	query := `
    SELECT 
        orders.id AS order_id,
        orders.order_date,
        people.id AS user_id,
        people.first_name,
        people.last_name,
        people.age,
        cars.id AS car_id,
        cars.name AS car_name,
        cars.power,
        cars.type,
        cars.year
    FROM 
        orders
    JOIN 
        people ON orders.user_id = people.id
    JOIN 
        cars ON orders.car_id = cars.id;
    `

	rows, err := db.Query(query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "Failed to query database"})
		return
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var o Order
		if err := rows.Scan(
			&o.OrderID, &o.OrderDate,
			&o.User.UserID, &o.User.FirstName, &o.User.LastName, &o.User.Age,
			&o.Car.CarID, &o.Car.Name, &o.Car.Power, &o.Car.Type, &o.Car.Year,
		); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"err": "Failed to scan row"})

			return
		}
		orders = append(orders, o)
	}

	ctx.JSON(http.StatusOK, orders)

}

// @Summary      Delete order by id
// @Description  delete order by user id
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        orderID path string true "Order ID"
// @Success      200  {object}  Order
// @Router       /api/orders/{orderID} [delete]
func (h *Handler) deleteOrderByID(ctx *gin.Context) {
	orderID := ctx.Param("orderID")

	query := `DELETE FROM orders WHERE id = $1`

	result, err := db.Exec(query, orderID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if rowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	ctx.Status(http.StatusOK)

}

// @Summary      Get order by user id
// @Description  get order by user id
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        userID path string true "User ID"
// @Success      200  {object}  Order
// @Router       /api/orders/{userID} [get]
func (h *Handler) getOrdersByUserID(ctx *gin.Context) {
	userID := ctx.Param("userID")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "userID is required"})
		return
	}

	query := `
    SELECT 
        orders.id AS order_id,
        orders.order_date,
        people.id AS user_id,
        people.first_name,
        people.last_name,
        people.age,
        cars.id AS car_id,
        cars.name AS car_name,
        cars.power,
        cars.type,
        cars.year
    FROM 
        orders
    JOIN 
        people ON orders.user_id = people.id
    JOIN 
        cars ON orders.car_id = cars.id
    WHERE 
        orders.user_id = $1;
    `

	rows, err := db.Query(query, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		if err := rows.Scan(&order.OrderID, &order.OrderDate,
			&order.User.UserID, &order.User.FirstName, &order.User.LastName, &order.User.Age,
			&order.Car.CarID, &order.Car.Name, &order.Car.Power, &order.Car.Type, &order.Car.Year); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to scan row"})
			return
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred during row iteration"})
		return
	}

	if len(orders) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No orders found for this user"})
		return
	}

	ctx.JSON(http.StatusOK, orders)
}
