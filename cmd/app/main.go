package main

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	_ "go_gin/cmd/docs"
	"go_gin/internal/config"
	"go_gin/internal/controller"
	"go_gin/internal/db"
	"go_gin/internal/domain/model"
	"go_gin/internal/middleware"
	"go_gin/internal/repository"
	"go_gin/internal/routes"
	"go_gin/internal/service"
	"go_gin/pkg/shutdown"
	"log"
	"net/http"
	"time"
)

//	@title 			Users Dan TodoList RESTful API
//	@version 		1.0.1
//	@description 	RESTful API Users Dan TodoList Menggunakan Auth JWT (JSON WEB TOKEN)

// @host localhost:3500
// @basePath /api
func main() {
	pgstore := db.NewPGStore(config.Database)
	dbs, _ := pgstore.Connect()
	validation := validator.New()
	repositoryTodolist := repository.NewTodolistRepository()
	repositoryUser := repository.NewUsersRepository()
	serviceUser := service.NewUsersService(dbs, repositoryUser, validation)
	serviceTodolist := service.NewTodoListService(dbs, validation, repositoryTodolist)
	controllerUser := controller.NewUsersController(serviceUser)
	controllerTodolist := controller.NewTodoListController(serviceTodolist)
	router := routes.Routes{Controller: controllerUser, Middleware: &middleware.Middleware{Repository: repositoryUser, DB: dbs}, TodoList: controllerTodolist}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Server.Port),
		Handler: router.Run(), //type gin.RouterGroup
	}
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	wait := shutdown.ShutDown(context.Background(), 2*time.Second, map[string]model.Operation{
		"http-server": func(ctx context.Context) error {
			return srv.Shutdown(context.Background())
		},
		"database": func(ctx context.Context) error {
			sql, _ := dbs.DB()
			return sql.Close()
		},
	})
	<-wait
}
