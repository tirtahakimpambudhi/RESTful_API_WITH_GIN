package model

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go_gin/internal/domain/model/web"
	"gorm.io/gorm"
)

type CustomSeeds interface {
	Run(ctx context.Context, DB *gorm.DB, fill int) error
}

type UsersRepository interface {
	GetUsers(ctx context.Context, DB *gorm.DB, query web.GetAllValue) Users
	GetUsersBySearch(ctx context.Context, DB *gorm.DB, query web.SearchValue) Users
	GetUserByID(ctx context.Context, DB *gorm.DB, ID uuid.UUID) (User, error)
	GetUserByEmail(ctx context.Context, DB *gorm.DB, Email string) (User, error)
	CreateUser(ctx context.Context, DB *gorm.DB, user User) error
	CreateUsers(ctx context.Context, DB *gorm.DB, users Users) error
	UpdateUserID(ctx context.Context, DB *gorm.DB, user User, ID uuid.UUID)
	DeleteUserByID(ctx context.Context, DB *gorm.DB, ID uuid.UUID)
	DeleteUsersByIDs(ctx context.Context, DB *gorm.DB, IDs []uuid.UUID)
	RestoreUserByID(ctx context.Context, DB *gorm.DB, ID uuid.UUID)
	RestoreUsersByIDs(ctx context.Context, DB *gorm.DB, IDs []uuid.UUID)
	UsersExistByID(ctx context.Context, DB *gorm.DB, ID uuid.UUID) bool
	UsersExistByIDs(ctx context.Context, DB *gorm.DB, IDs []uuid.UUID) bool
}

type TodoListRepository interface {
	GetTodoListsSearch(ctx context.Context, DB *gorm.DB, query web.SearchValue, userID uuid.UUID) TodoLists
	GetTodoLists(ctx context.Context, DB *gorm.DB, query web.GetAllValue, userID uuid.UUID) TodoLists
	GetTodoListByID(ctx context.Context, DB *gorm.DB, ID int, userID uuid.UUID) (TodoList, error)
	CreateTodoList(ctx context.Context, DB *gorm.DB, todolist TodoList) error
	CreateTodoLists(ctx context.Context, DB *gorm.DB, todolists TodoLists) error
	UpdateTodoListByID(ctx context.Context, DB *gorm.DB, todolist TodoList, ID int, userID uuid.UUID)
	DeleteTodoListByID(ctx context.Context, DB *gorm.DB, ID int, userID uuid.UUID)
	DeleteTodoListsByIDs(ctx context.Context, DB *gorm.DB, IDs []int, userID uuid.UUID)
	TodoListExistByID(ctx context.Context, DB *gorm.DB, ID int, userID uuid.UUID) bool
	TodoListsExistByIDs(ctx context.Context, DB *gorm.DB, IDs []int, userID uuid.UUID) bool
}

type UsersService interface {
	FindUsersBySearch(ctx context.Context, params web.SearchQuery) (UsersResponses, web.Pagination, error)
	FindUsers(ctx context.Context, params web.GetAllQuery) (UsersResponses, web.Pagination, error)
	FindUserByID(ctx context.Context, ID uuid.UUID) (UserResponse, error)
	CreateUser(ctx context.Context, user UserRequest) error
	CreateUsers(ctx context.Context, users UsersRequests) error
	UpdateUserID(ctx context.Context, user UserLoginUpdateRequest, ID uuid.UUID) error
	DeleteUserByID(ctx context.Context, ID uuid.UUID) error
	DeleteUsersByIDs(ctx context.Context, IDs []uuid.UUID) error
	RestoreUserByID(ctx context.Context, ID uuid.UUID) error
	RestoreUsersByIDs(ctx context.Context, IDs []uuid.UUID) error
	LoginUsers(ctx context.Context, users UserLoginUpdateRequest) (string, string, error)
	LogoutUsers(ctx context.Context, refreshToken string) error
	RefreshTokenUser(ctx context.Context, refreshToken string) (string, error)
}

type TodoListService interface {
	FindTodoListsBySearch(ctx context.Context, params web.Params) (responses TodoListResponses, pagination web.Pagination, errService error)
	FindTodoLists(ctx context.Context, params web.Params) (responses TodoListResponses, pagination web.Pagination, errService error)
	FindTodoListByID(ctx context.Context, params web.Params) (response TodoListResponse, errService error)
	CreateTodoList(ctx context.Context, request TodoListRequest, params web.Params) (errService error)
	CreatesTodoLists(ctx context.Context, requests TodoListRequests, params web.Params) (errService error)
	UpdateTodoList(ctx context.Context, request TodoListRequest, params web.Params) (errService error)
	DeleteTodoList(ctx context.Context, params web.Params) (errService error)
	DeletesTodoLists(ctx context.Context, params web.Params) (errService error)
}

type UsersController interface {
	GetAll(c *gin.Context)
	GetBySearch(c *gin.Context)
	GetByID(c *gin.Context)
	CreateUser(c *gin.Context)
	CreateUsers(c *gin.Context)
	LoginUser(c *gin.Context)
	RefreshTokenUser(c *gin.Context)
	UpdateUserID(c *gin.Context)
	DeleteUserByID(c *gin.Context)
	DeleteUsersByIDs(c *gin.Context)
	RestoreUserByID(c *gin.Context)
	RestoreUsersByIDs(c *gin.Context)
	LogoutUser(c *gin.Context)
}

type TodoListController interface {
	GetTodoListSearch(c *gin.Context)
	GetTodoListAll(c *gin.Context)
	GetTodoListByID(c *gin.Context)
	CreateTodoList(c *gin.Context)
	CreatesTodoLists(c *gin.Context)
	UpdateTodoList(c *gin.Context)
	DeleteTodoList(c *gin.Context)
	DeleteTodoLists(c *gin.Context)
}
