package service

import (
	"golang-structure-template-with-di/app/models/entity"
	"golang-structure-template-with-di/app/models/request"
	"golang-structure-template-with-di/app/repository"
	"strings"
)

type ExampleProjectService struct {
	ExampleProjectRepository *repository.ExampleProjectRepository `di.inject:"ExampleProjectRepository"`
}

// ExampleProject http request
func (s *ExampleProjectService) ExampleProject(requestContainerProject request.CreateExampleProject) *entity.ExampleProject {

	// Create a new entity.ExampleProject
	//lowercase name request
	projectKey := strings.ToLower(requestContainerProject.Name)
	projectKey = strings.ReplaceAll(projectKey, " ", "-")
	ExampleProject := entity.ExampleProject{
		Key:  projectKey,
		Name: requestContainerProject.Name,
	}
	return s.ExampleProjectRepository.Create(&ExampleProject)
}
