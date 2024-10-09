package repository

type User interface {

}

type Repository struct {
	User
}

func NewRepository() *Repository{
	return &Repository{}
}