package database

import (
	"golang-structure-template-with-di/app/models/entity"
)

func Migrate() {
	db := PostgreDB
	err := db.AutoMigrate(entity.ExampleProject{})
	if err != nil {
		panic("ERROR MIGRATE example project : " + err.Error())
	}

}
