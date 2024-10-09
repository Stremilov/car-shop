package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type User struct {
	UserID    int    `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
}

type UserUpdate struct {
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Age       *int    `json:"age,omitempty"`
}

// @Summary      Add new user
// @Description  add user to the database
// @Tags         users
// @Accept       json
// @Produce      json
// @Param request body User true "body"
// @Success      201  {object}  User
// @Router       /api/user/ [post]
func (h *Handler) addUser(ctx *gin.Context) {
	var p User

	if err := ctx.ShouldBindJSON(&p); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	query := "INSERT INTO people (first_name, last_name, age) VALUES ($1, $2, $3)"
	_, err := db.Exec(query, p.FirstName, p.LastName, p.Age)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert data into database"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Person added successfully"})

}

// @Summary      Get all users
// @Description  get all users from database
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {object}  User
// @Router       /api/user/get-all [get]
func (h *Handler) getAllUsers(ctx *gin.Context) {
	rows, err := db.Query("SELECT first_name, last_name, age FROM people")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query database"})
		return
	}
	defer rows.Close()

	var people []User
	for rows.Next() {
		var p User
		if err := rows.Scan(&p.FirstName, &p.LastName, &p.Age); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row"})
			return
		}
		people = append(people, p)
	}

	ctx.JSON(http.StatusOK, people)
}

// @Summary      Update user info by id
// @Description  update user info by user id
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        userID path string true "User ID"
// @Param request body User true "body"
// @Success      201  {object}  User
// @Router       /api/user/{userID} [patch]
func (h *Handler) updateUserInfoByID(ctx *gin.Context) {
	userID := ctx.Param("userID")

	var userUpdate UserUpdate

	if err := ctx.ShouldBindJSON(&userUpdate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	values := []interface{}{}
	setClauses := []string{}

	if userUpdate.FirstName != nil {
		setClauses = append(setClauses, "first_name = $1")
		values = append(values, *userUpdate.FirstName)
	}

	if userUpdate.LastName != nil {
		setClauses = append(setClauses, "last_name = $"+strconv.Itoa(len(values)+1))
		values = append(values, *userUpdate.LastName)
	}

	if userUpdate.Age != nil {
		setClauses = append(setClauses, "age = $"+strconv.Itoa(len(values)+1))
		values = append(values, *userUpdate.Age)
	}

	if len(setClauses) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	query := fmt.Sprintf("UPDATE people SET %s WHERE id = $%d;",
		strings.Join(setClauses, ", "), len(values)+1)

	values = append(values, userID)

	_, err := db.Exec(query, values...)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid request payload"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Data updated successfully"})

}

// @Summary      Get user info by id
// @Description  get user info by user id
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        userID path string true "User ID"
// @Success      200  {object}  User
// @Router       /api/user/{userID} [get]
func (h *Handler) getUserByID(ctx *gin.Context) {
	userID := ctx.Param("userID")

	query := `
	SELECT 
		people.first_name,
		people.last_name,
		people.age
	FROM 
		people
	WHERE 
		people.id = $1
	`

	row := db.QueryRow(query, userID)

	var user User
	if err := row.Scan(&user.FirstName, &user.LastName, &user.Age); err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid request"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row"})

		return
	}

	ctx.JSON(http.StatusOK, user)
}

// @Summary      Delete user info by id
// @Description  delete user info by user id
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        userID path string true "User ID"
// @Success      200  {object}  User
// @Router       /api/user/{userID} [delete]
func (h *Handler) deleteUserByID(ctx *gin.Context) {
	userID := ctx.Param("userID")

	query := `DELETE FROM people WHERE id = $1`

	result, err := db.Exec(query, userID)
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
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	ctx.Status(http.StatusOK)

}
