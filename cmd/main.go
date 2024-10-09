package main

import (
	"log"

	goapi "github.com/Stremilov/car-shop"
	"github.com/Stremilov/car-shop/pkg/handler"
	"github.com/Stremilov/car-shop/pkg/repository"
	"github.com/Stremilov/car-shop/pkg/service"
)

// @title GoAPI test project
// @version 1.0
// @description API documentation for test project

func main() {
	repos := repository.NewRepository()
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	server := new(goapi.Server)
	if err := server.Run("8080", handlers.InitRoutesAndDB()); err != nil {
		log.Fatalf("Error running server: %v", err)
	}
}
