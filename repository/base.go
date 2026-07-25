package repository

import "gorm.io/gorm"

type repo struct {
	DB *gorm.DB
}

func NewRepository(gormDB *gorm.DB) Repository {
	return &repo{DB: gormDB}
}
