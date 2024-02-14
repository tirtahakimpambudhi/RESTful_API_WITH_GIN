package routes

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go_gin/internal/controller"
	"go_gin/internal/domain/model"
	"go_gin/internal/middleware"
)

type Routes struct {
	Controller model.UsersController
	Middleware *middleware.Middleware
	TodoList   *controller.TodoListController
}

func (r *Routes) Run() *gin.Engine {
	router := gin.Default()

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	admin := api.Group("/admin")
	moderator := api.Group("/moderator")
	//Admin And Moderator
	admin.GET("/users", r.Middleware.Authentication, r.Middleware.AuthorizationRoleAdmin, r.Controller.GetAll)
	admin.GET("/users/search", r.Middleware.Authentication, r.Middleware.AuthorizationRoleAdmin, r.Controller.GetBySearch)
	admin.GET("/user/:id", r.Middleware.Authentication, r.Middleware.AuthorizationRoleAdmin, r.Controller.GetByID)
	admin.POST("/registers", r.Middleware.Authentication, r.Middleware.AuthorizationRoleAdmin, r.Controller.CreateUsers)
	admin.PATCH("/user/:id", r.Middleware.Authentication, r.Middleware.AuthorizationRoleAdmin, r.Controller.RestoreUserByID)
	admin.PATCH("/users", r.Middleware.Authentication, r.Middleware.AuthorizationRoleAdmin, r.Controller.RestoreUsersByIDs)
	admin.DELETE("/user/:id", r.Middleware.Authentication, r.Middleware.AuthorizationRoleAdmin, r.Controller.DeleteUserByID)
	admin.DELETE("/users", r.Middleware.Authentication, r.Middleware.AuthorizationRoleAdmin, r.Controller.DeleteUsersByIDs)

	moderator.GET("/users", r.Middleware.Authentication, r.Middleware.AuthorizationRoleModerator, r.Middleware.IsLogin, r.Controller.GetAll)
	moderator.POST("/registers", r.Middleware.Authentication, r.Middleware.AuthorizationRoleModerator, r.Middleware.IsLogin, r.Controller.CreateUsers)
	moderator.GET("/user/:id", r.Middleware.Authentication, r.Middleware.AuthorizationRoleModerator, r.Middleware.IsLogin, r.Controller.GetByID)
	//All Role

	//todolist
	api.GET("/user/:id/todolists", r.Middleware.IsLogin, r.TodoList.GetTodoListAll)
	api.GET("/user/:id/todolists/s", r.Middleware.IsLogin, r.TodoList.GetTodoListSearch)
	api.GET("/user/:id/todolist", r.Middleware.IsLogin, r.TodoList.GetTodoListByID)
	api.POST("/user/:id/todolist", r.Middleware.IsLogin, r.TodoList.CreateTodoList)
	api.POST("/user/:id/todolists", r.Middleware.IsLogin, r.TodoList.CreatesTodoLists)
	api.PUT("/user/:id/todolist", r.Middleware.IsLogin, r.TodoList.UpdateTodoList)
	api.DELETE("/user/:id/todolist", r.Middleware.IsLogin, r.TodoList.DeleteTodoList)
	api.DELETE("/user/:id/todolists", r.Middleware.IsLogin, r.TodoList.DeleteTodoLists)

	//users
	api.GET("/auth", r.Middleware.Authentication, r.Middleware.AuthorizationAllRole)
	api.GET("/refresh", r.Controller.RefreshTokenUser)
	api.POST("/register", r.Controller.CreateUser)
	api.POST("/login", r.Controller.LoginUser)
	api.PUT("/user/:id", r.Middleware.IsLogin, r.Controller.UpdateUserID)
	api.DELETE("/logout", r.Controller.LogoutUser)

	return router
}
