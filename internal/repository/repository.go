package repository

import "backend/internal/models"

type DatabaseRepository interface {
	GetMovies() ([]*models.Movie, error)
	GetAllMovies() ([]*models.Movie, error)
	GetMovie() (*models.Movie, error)
	PutMovies() ([]*models.Movie, error)
	PutMovie() (*models.Movie, error)
}
