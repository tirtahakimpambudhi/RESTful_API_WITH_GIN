package repository

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"go_gin/internal/config"
	"go_gin/internal/domain/model"
	"go_gin/internal/domain/model/web"
	"go_gin/pkg/helper"
	"gorm.io/gorm"
	"strings"
)

type TodolistRepository struct {
}

func NewTodolistRepository() *TodolistRepository {
	return &TodolistRepository{}
}

// perkecil argumen
func (t *TodolistRepository) scanTodoList(rows *sql.Rows) (model.TodoList, error) {
	var todolist model.TodoList
	err := rows.Scan(
		&todolist.TaskID,
		&todolist.UserID,
		&todolist.TaskName,
		&todolist.Description,
		&todolist.DueDate,
		&todolist.Priority,
		&todolist.Completed,
		&todolist.CreatedAt,
		&todolist.UpdatedAt,
	)
	return todolist, err
}
func (t *TodolistRepository) GetTodoListsSearch(ctx context.Context, DB *gorm.DB, query web.SearchValue, userId uuid.UUID) model.TodoLists {
	valueSearch := []string{"%", query.Search, "%"}
	key := strings.Join(valueSearch, "")
	var todolists model.TodoLists
	rows, err := DB.WithContext(ctx).Model(&model.TodoList{}).Where("user_id = ?", userId).Where("task_name LIKE ?", key).Or("description LIKE ?", key).Offset(int(query.Offset)).Limit(config.Other.Limit).Rows()
	helper.Panic(err)
	defer rows.Close()
	for rows.Next() {
		todolist, err := t.scanTodoList(rows)
		if err != nil {
			helper.Panic(err)
			return nil
		}
		todolists = append(todolists, todolist)
	}
	return todolists
}

func (t *TodolistRepository) GetTodoLists(ctx context.Context, DB *gorm.DB, query web.GetAllValue, userId uuid.UUID) model.TodoLists {
	var todolists model.TodoLists
	rows, err := DB.WithContext(ctx).Model(&model.TodoList{}).Where("user_id = ?", userId).Offset(int(query.Offset)).Limit(config.Other.Limit).Rows()
	helper.Panic(err)
	defer rows.Close()
	for rows.Next() {
		todolist, err := t.scanTodoList(rows)
		if err != nil {
			helper.Panic(err)
			return nil
		}
		todolists = append(todolists, todolist)
	}
	return todolists
}

func (t *TodolistRepository) GetTodoListByID(ctx context.Context, DB *gorm.DB, ID int, userId uuid.UUID) (model.TodoList, error) {
	var todolist model.TodoList
	err := DB.WithContext(ctx).Model(&model.TodoList{}).Where("user_id = ?", userId).Where("id = ?", ID).Take(&todolist).Error
	if err != nil {
		return model.TodoList{}, err
	}
	return todolist, nil
}

func (t *TodolistRepository) CreateTodoList(ctx context.Context, DB *gorm.DB, todolist model.TodoList) error {
	err := DB.WithContext(ctx).Model(&model.TodoList{}).Create(&todolist).Error
	return err
}

func (t *TodolistRepository) CreateTodoLists(ctx context.Context, DB *gorm.DB, todolists model.TodoLists) error {
	var err error
	if len(todolists) > config.Other.LimitInsert {
		err = DB.WithContext(ctx).CreateInBatches(&todolists, config.Other.LimitInsert).Error
	} else {
		err = DB.WithContext(ctx).Create(&todolists).Error
	}
	return err
}

func (t *TodolistRepository) UpdateTodoListByID(ctx context.Context, DB *gorm.DB, todolist model.TodoList, ID int, userId uuid.UUID) {
	err := DB.WithContext(ctx).Model(&model.TodoList{}).Where("user_id = ?", userId).Where("id = ?", ID).Updates(&todolist).Error
	helper.Panic(err)
}

func (t *TodolistRepository) DeleteTodoListByID(ctx context.Context, DB *gorm.DB, ID int, userId uuid.UUID) {
	err := DB.WithContext(ctx).Model(&model.TodoList{}).Where("user_id = ?", userId).Where("id = ?", ID).Delete(&model.TodoList{}).Error
	helper.Panic(err)
}

func (t *TodolistRepository) DeleteTodoListsByIDs(ctx context.Context, DB *gorm.DB, IDs []int, userId uuid.UUID) {
	err := DB.WithContext(ctx).Model(&model.TodoList{}).Where("user_id = ?", userId).Where("id IN ?", IDs).Delete(&model.TodoList{}).Error
	helper.Panic(err)
}

func (t *TodolistRepository) TodoListExistByID(ctx context.Context, DB *gorm.DB, ID int, userId uuid.UUID) bool {
	var count int64
	err := DB.WithContext(ctx).Where("user_id = ?", userId).Where("id = ?", ID).Count(&count).Error
	helper.Panic(err)
	return count == 1
}

func (t *TodolistRepository) TodoListsExistByIDs(ctx context.Context, DB *gorm.DB, IDs []int, userId uuid.UUID) bool {
	var count int64
	err := DB.WithContext(ctx).Where("user_id = ?", userId).Where("id IN ?", IDs).Count(&count).Error
	helper.Panic(err)
	return count == int64(len(IDs))
}
