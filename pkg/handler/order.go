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
