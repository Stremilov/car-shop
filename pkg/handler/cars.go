package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Car struct {
	CarID int    `json:"car_id"`
	Name  string `json:"name"`
	Power string `json:"power"`
	Type  string `json:"type"`
	Year  int    `json:"year"`
}

type CarUpdate struct {
	Name  string `json:"name"`
	Power string `json:"power"`
	Type  string `json:"type"`
	Year  int    `json:"year"`
}

// @Summary      Add new car
// @Description  add new car
// @Tags         cars
// @Accept       json
// @Produce      json
// @Success      201  {object}  Car
// @Router       /api/car/ [post]
func (h *Handler) addCar(c *gin.Context) {
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

// @Summary      Get all cars
// @Description  get all cars
// @Tags         cars
// @Accept       json
// @Produce      json
// @Success      200  {object} Car
// @Router       /api/car/get-all [get]
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

// @Summary      Get car by id
// @Description  get car by id
// @Tags         cars
// @Accept       json
// @Produce      json
// @Param        carID path string true "car ID"
// @Success      200  {object}  Car
// @Router       /api/car/{carID} [get]
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

// @Summary      Delete car by id
// @Description  delete car by user id
// @Tags         cars
// @Accept       json
// @Produce      json
// @Param        carID path string true "Car ID"
// @Success      200  {object}  Car
// @Router       /api/car/{carID} [delete]
func (h *Handler) deleteCarByID(ctx *gin.Context) {
	carID := ctx.Param("carID")

	query := `DELETE FROM cars WHERE id = $1`

	result, err := db.Exec(query, carID)
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
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Car not found"})
		return
	}

	ctx.Status(http.StatusOK)

}

// @Summary      Update car info by id
// @Description  update car info by user id
// @Tags         cars
// @Accept       json
// @Produce      json
// @Param        carID path string true "Car ID"
// @Param request body Car true "body"
// @Success      201  {object}  Car
// @Router       /api/car/{userID} [patch]
func (h *Handler) updateCarInfoByID(ctx *gin.Context) {
	carID := ctx.Param("carID")

	var carUpdate CarUpdate
	if err := ctx.ShouldBindJSON(&carUpdate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	values := []interface{}{}
	setClauses := []string{}

	if carUpdate.Name != "" {
		setClauses = append(setClauses, "name = $1")
		values = append(values, carUpdate.Name)
	}

	if carUpdate.Power != "" {
		setClauses = append(setClauses, "power = $"+strconv.Itoa(len(values)+1))
		values = append(values, carUpdate.Power)
	}

	if carUpdate.Type != "" {
		setClauses = append(setClauses, "type = $"+strconv.Itoa(len(values)+1))
		values = append(values, carUpdate.Type)
	}

	if carUpdate.Year != 0 {
		setClauses = append(setClauses, "year = $"+strconv.Itoa(len(values)+1))
		values = append(values, carUpdate.Year)
	}

	if len(setClauses) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	query := fmt.Sprintf("UPDATE cars SET %s WHERE id = $%d;",
		strings.Join(setClauses, ", "), len(values)+1)

	values = append(values, carID)

	_, err := db.Exec(query, values...)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update data"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Data updated successfully"})
}
