package handler

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func initDB() {
	connStr := "user=levstremilov password=postgres dbname=testdb sslmode=disable"

	var err error

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	query := `
	CREATE TABLE IF NOT EXISTS people (
		id SERIAL PRIMARY KEY,
		first_name VARCHAR(50),
		last_name VARCHAR(50),
		age INTEGER
	);
	
	CREATE TABLE IF NOT EXISTS cars (
		id SERIAL PRIMARY KEY,
		name VARCHAR(50),
		power INTEGER,
		type VARCHAR(10),
		year INTEGER
	);
	
	CREATE TABLE IF NOT EXISTS orders (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES people(id),
		car_id INTEGER REFERENCES cars(id),
		order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err = db.Exec(query)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	fmt.Println("Connected to the database successfully and ensured table exists!")
}
