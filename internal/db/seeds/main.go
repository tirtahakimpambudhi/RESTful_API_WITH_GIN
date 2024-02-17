package main

import (
	"fmt"
	"go_gin/internal/config"
	"go_gin/internal/db"
	"go_gin/internal/db/seeds/seeder"
	"go_gin/internal/domain/model"
	"go_gin/internal/handler"
	"go_gin/internal/repository"
	"log"
)

func main() {
	seedMethod, fill := handler.HandleArgs()
	dbs, err := db.NewPGStore(config.Database).Connect()
	if err != nil {
		fmt.Println(err.Error())
	}
	userSeed := seeder.NewUserSeed(repository.NewUsersRepository())
	seeds := seeder.NewSeeder(dbs, []model.CustomSeeds{userSeed})
	errExec := seeder.Execute(seedMethod, fill, seeds)
	if errExec != nil {
		log.Fatal(errExec.Error())
	}

}
