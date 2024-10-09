package service

import "github.com/Stremilov/car-shop/pkg/repository"

type User interface {

}

type Service struct {
	User
}

func NewService(repos *repository.Repository) *Service{
	return &Service{}
}