package main

import (
	"log"

	goapi "github.com/Stremilov/car-shop"
	"github.com/Stremilov/car-shop/pkg/handler"
)

func main() {
	handlers := new(handler.Handler)

	server := new(goapi.Server)
	if err := server.Run("8080", handlers.InitRoutesAndDB()); err != nil {
		log.Fatalf("Error running server: %v", err)
	}
}
