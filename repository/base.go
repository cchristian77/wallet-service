package repository

import "gorm.io/gorm"

type repo struct {
	DB *gorm.DB
}

// NewRepository initializes a new Repository instance, managing database interactions.
func NewRepository(gormDB *gorm.DB) Repository {
	return &repo{DB: gormDB}
}
