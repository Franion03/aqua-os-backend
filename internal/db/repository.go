package db

import "github.com/franion03/aqua-os-backend/internal/models"

// Repository defines the data access contract.
// Implementations: sqlite.go (dev), dynamo.go (prod).
type Repository interface {
	GetAllLevels() ([]models.Level, error)
	GetLevel(id int) (*models.Level, error)
	GetExercises(levelID *int) ([]models.Exercise, error)
	AddExercise(req models.ExerciseCreate) (*models.Exercise, error)
	DeleteExercise(id int) error
	Close() error
}
