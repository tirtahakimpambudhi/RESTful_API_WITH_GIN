package main

import (
	"fmt"
	"go_gin/internal/config"
	"go_gin/internal/db"
	"go_gin/internal/domain/model"
)

// just for emergency
func migration() {
	pgstore := db.NewPGStore(config.Database)
	dbs, err := pgstore.Connect()
	if err != nil {
		panic(err.Error())
	}
	err = dbs.AutoMigrate(&model.User{})
	if err != nil {
		panic(fmt.Errorf("error migrating users %s", err.Error()))
	}

	err = dbs.AutoMigrate(&model.TodoList{})
	if err != nil {
		panic(fmt.Errorf("error migrating users %s", err.Error()))
	}
}

func main() {
	migration()
}
