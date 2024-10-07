package handler

import (
	"net/http"
	"database/sql"

	"github.com/gin-gonic/gin"
)

type Car struct {
	CarID int    `json:"car_id"`
	Name  string `json:"name"`
	Power string `json:"power"`
	Type  string `json:"type"`
	Year  int    `json:"year"`
}

func (h *Handler) addCar(c*gin.Context) {
	var car Car

	if err := c.ShouldBindJSON(&car); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	query := `INSERT INTO cars (name, power, type, year) VALUES ($1, $2, $3, $4)`
	_, err := db.Exec(query, car.Name, car.Power, car.Type, car.Year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert data into database"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Car added successfully"})
}

func (h *Handler) getAllCars(ctx *gin.Context) {
	rows, err := db.Query("SELECT name, power, type, year FROM cars")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query database"})
		return
	}
	defer rows.Close()

	var cars []Car
	for rows.Next() {
		var car Car
		if err := rows.Scan(&car.Name, &car.Power, &car.Type, &car.Year); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row"})
			return
		}
		cars = append(cars, car)
	}

	ctx.JSON(http.StatusOK, cars)
}

func (h *Handler) getCarByID(ctx *gin.Context) {
	carID := ctx.Param("carID")

	query := `
	SELECT 
		cars.name,
		cars.power,
		cars.type,
		cars.year
	FROM
		cars
	WHERE
		cars.id = $1
	`

	row := db.QueryRow(query, carID)

	var car Car
	if err := row.Scan(&car.Name, &car.Power, &car.Type, &car.Year); err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Car not found"})

			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to find car"})
		return
	}

	ctx.JSON(http.StatusOK, car)

}