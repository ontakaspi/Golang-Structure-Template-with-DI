package repository

import (
	"golang-structure-template-with-di/app/models/entity"
	"golang-structure-template-with-di/libraries/httpResponse"
	"gorm.io/gorm"
)

// ExampleProjectRepository as @Service
type ExampleProjectRepository struct {
	//db *gorm.DB
	db *gorm.DB `di.inject:"dbPostgres"`
}

func (r ExampleProjectRepository) Create(ExampleProject *entity.ExampleProject) *entity.ExampleProject {
	if err := r.db.Create(ExampleProject).Error; err != nil {
		httpResponse.InternalServerError(err)
	}
	return ExampleProject
}
