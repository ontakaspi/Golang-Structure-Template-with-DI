package entity

import "gorm.io/gorm"

type ExampleProject struct {
	gorm.Model
	Key  string `json:"key"`
	Name string `json:"name"`
}
