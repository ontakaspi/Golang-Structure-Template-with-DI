package main

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/goioc/di"
	"golang-structure-template-with-di/app/controller"
	"golang-structure-template-with-di/app/repository"
	"golang-structure-template-with-di/app/service"
	"golang-structure-template-with-di/config/database"
	"reflect"
)

func initializeBean() {

	/* GLOBAL - inject the all global */
	_, _ = di.RegisterBeanFactory("dbPostgres", di.Prototype, func(ctx context.Context) (interface{}, error) {
		db := database.PostgreDB
		return db, nil
	})
	_, _ = di.RegisterBeanFactory("validator", di.Singleton, func(ctx context.Context) (interface{}, error) {
		validate := validator.New()
		return validate, nil
	})

	/* CONTROLLER - inject the all controller */
	_, _ = di.RegisterBean("ExampleProjectController", reflect.TypeOf((*controller.ExampleProjectController)(nil)))

	/* SERVICE - inject the all service */
	_, _ = di.RegisterBean("AuthService", reflect.TypeOf((*service.AuthService)(nil)))
	_, _ = di.RegisterBean("ExampleProjectService", reflect.TypeOf((*service.ExampleProjectService)(nil)))

	/* REPOSITORY - inject the all repository */
	_, _ = di.RegisterBean("ExampleProjectRepository", reflect.TypeOf((*repository.ExampleProjectRepository)(nil)))

	/* ----------------- */
	_ = di.InitializeContainer()
}
