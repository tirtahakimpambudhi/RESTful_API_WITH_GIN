package routes

import (
	"github.com/gin-gonic/gin"
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
	//Admin And Moderator
	router.GET("/api/admin/users", r.Middleware.Authentication, r.Middleware.AuthorizationRoleAdmin, r.Controller.GetAll)
	router.POST("/api/admin/registers", r.Middleware.Authentication, r.Middleware.AuthorizationRoleAdmin, r.Controller.CreateUsers)
	router.GET("/api/admin/user/:id", r.Middleware.Authentication, r.Middleware.AuthorizationRoleAdmin, r.Controller.GetByID)
	router.DELETE("/api/admin/user/:id", r.Middleware.Authentication, r.Middleware.AuthorizationRoleAdmin, r.Controller.DeleteUserByID)
	router.PATCH("/api/admin/user/:id", r.Middleware.Authentication, r.Middleware.AuthorizationRoleAdmin, r.Controller.RestoreUserByID)
	router.DELETE("/api/admin/users", r.Middleware.Authentication, r.Middleware.AuthorizationRoleAdmin, r.Controller.DeleteUsersByIDs)
	router.PATCH("/api/admin/users", r.Middleware.Authentication, r.Middleware.AuthorizationRoleAdmin, r.Controller.RestoreUsersByIDs)

	router.GET("/api/moderator/users", r.Middleware.Authentication, r.Middleware.AuthorizationRoleModerator, r.Middleware.IsLogin, r.Controller.GetAll)
	router.POST("/api/moderator/registers", r.Middleware.Authentication, r.Middleware.AuthorizationRoleModerator, r.Middleware.IsLogin, r.Controller.CreateUsers)
	router.GET("/api/moderator/user/:id", r.Middleware.Authentication, r.Middleware.AuthorizationRoleModerator, r.Middleware.IsLogin, r.Controller.GetByID)
	//All Role

	//todolist
	router.GET("/api/user/:id/todolist", r.Middleware.IsLogin, r.TodoList.GetTodoListAll)
	router.GET("/api/user/:id/todolist/s", r.Middleware.IsLogin, r.TodoList.GetTodoListSearch)
	router.POST("/api/user/:id/todolist", r.Middleware.IsLogin, r.TodoList.CreateTodoList)
	router.PUT("/api/user/:id/todolist", r.Middleware.IsLogin, r.TodoList.UpdateTodoList)
	router.POST("/api/user/:id/todolists", r.Middleware.IsLogin, r.TodoList.CreatesTodoLists)
	router.DELETE("/api/user/:id/todolist", r.Middleware.IsLogin, r.TodoList.DeleteTodoList)
	router.DELETE("/api/user/:id/todolists", r.Middleware.IsLogin, r.TodoList.DeleteTodoLists)

	//users
	router.GET("/api/authentication/", r.Middleware.Authentication, r.Middleware.AuthorizationAllRole)
	router.PUT("/api/user/:id", r.Middleware.IsLogin, r.Controller.UpdateUserID)
	router.POST("/api/register", r.Controller.CreateUser)
	router.POST("/api/login", r.Controller.LoginUser)
	router.DELETE("/api/logout", r.Controller.LogoutUser)
	router.GET("/api/refresh", r.Controller.RefreshTokenUser)

	return router
}
