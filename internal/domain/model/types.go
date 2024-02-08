package model

import (
	"context"
	"github.com/google/uuid"
)

type Operation func(ctx context.Context) error
type Date struct {
	Year  int `json:"year" validate:"required,gte=1900,lte=2100"`
	Month int `json:"month" validate:"required,gte=1,lte=12"`
	Day   int `json:"day" validate:"required,gte=1,lte=31"`
}
type Users []User
type UsersResponses []UserResponse
type UsersRequests []UserRequest
type TodoLists []TodoList
type TodoListRequests []TodoListRequest
type TodoListResponses []TodoListResponse

type TodoListRequestsValidation struct {
	TodoList []TodoListRequest `json:"todo_list" validate:"required,dive"`
}
type UsersRequestsValidation struct {
	Users []UserRequest `json:"users" validate:"required,dive"`
}

func (u Users) ToUsersResponses() UsersResponses {
	var userResponse UsersResponses
	for _, user := range u {
		userResponse = append(userResponse, *user.ToUserResponse())
	}
	return userResponse
}
func (u UsersRequests) ToUsers() Users {
	var users Users
	for _, user := range u {
		users = append(users, *user.ToUser())
	}
	return users
}
func (t TodoListRequests) ToTodoLists(user_id uuid.UUID) TodoLists {
	var todolists TodoLists
	for _, todolistRequest := range t {
		todolists = append(todolists, *todolistRequest.ToTodoList(user_id))
	}
	return todolists
}
func (t TodoLists) ToTodoListResponses() TodoListResponses {
	var todolists TodoListResponses
	for _, todolistResponse := range t {
		todolists = append(todolists, *todolistResponse.ToTodoListResponse())
	}
	return todolists
}
