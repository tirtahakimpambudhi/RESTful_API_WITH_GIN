package model

import (
	"github.com/google/uuid"
	"time"
)

type TodoList struct {
	TaskID      int        `json:"task_id" gorm:"primaryKey;column:task_id"`
	UserID      uuid.UUID  `json:"user_id" gorm:"column:user_id"`
	TaskName    string     `json:"task_name" gorm:"column:task_name"`
	Description string     `json:"description" gorm:"column:description"`
	DueDate     *time.Time `json:"due_date" gorm:"column:due_date"`
	Priority    int        `json:"priority" gorm:"column:priority"`
	Completed   bool       `json:"completed" gorm:"column:completed"`
	CreatedAt   time.Time  `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"column:updated_at"`
	User        User       `gorm:"foreignKey:user_id;references:id" json:"user"`
}

func (t *TodoList) TableName() string {
	return "todolist"
}

type TodoListResponse struct {
	TaskID      int        `json:"task_id" gorm:"primaryKey;column:task_id"`
	TaskName    string     `json:"task_name" gorm:"column:task_name"`
	Description string     `json:"description" gorm:"column:description"`
	DueDate     *time.Time `json:"due_date" gorm:"column:due_date"`
	Priority    int        `json:"priority" gorm:"column:priority"`
	Completed   bool       `json:"completed" gorm:"column:completed"`
	CreatedAt   time.Time  `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"column:updated_at"`
}
type TodoListRequest struct {
	TaskName    string `json:"task_name" gorm:"column:task_name" validate:"required"`
	Description string `json:"description" gorm:"column:description" validate:"required"`
	DueDate     *Date  `json:"due_date" gorm:"column:due_date" validate:"required"`
	Priority    int    `json:"priority" gorm:"column:priority" validate:"min=1"`
	Completed   bool   `json:"completed" gorm:"column:completed" validate:"eq=true|eq=false"`
}

func (t *TodoListRequest) ToTodoList(user_id uuid.UUID) *TodoList {
	dateTime := time.Date(t.DueDate.Year, time.Month(t.DueDate.Month), t.DueDate.Day, 0, 0, 0, 0, time.UTC)
	return &TodoList{
		UserID:      user_id,
		TaskName:    t.TaskName,
		Description: t.Description,
		DueDate:     &dateTime,
		Priority:    t.Priority,
		Completed:   t.Completed,
	}
}

func (t *TodoList) ToTodoListResponse() *TodoListResponse {
	return &TodoListResponse{
		TaskID:      t.TaskID,
		TaskName:    t.TaskName,
		Description: t.Description,
		DueDate:     t.DueDate,
		Priority:    t.Priority,
		Completed:   t.Completed,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}
