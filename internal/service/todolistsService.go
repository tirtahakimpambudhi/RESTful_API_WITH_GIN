package service

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"go_gin/internal/config"
	"go_gin/internal/domain/model"
	"go_gin/internal/domain/model/web"
	"go_gin/internal/exception"
	"go_gin/internal/repository"
	"gorm.io/gorm"
	"math"
	"strings"
)

type TodoListService struct {
	DB         *gorm.DB
	Validator  *validator.Validate
	Repository *repository.TodolistRepository
}

func NewTodoListService(DB *gorm.DB, validator *validator.Validate, repository *repository.TodolistRepository) *TodoListService {
	return &TodoListService{DB: DB, Validator: validator, Repository: repository}
}

func (t *TodoListService) CreateTodoList(ctx context.Context, request model.TodoListRequest, params web.Params) (errService error) {
	tx := t.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err := r.(error)
			errService = exception.NewError(err, exception.ErrorInternalServer)
		}
	}()
	userID, validUUID := params.UserID.ToUUID()
	if validUUID != nil {
		tx.Rollback()
		errService = exception.NewError(validUUID, exception.ErrorBadRequest)
		return
	}
	badRequest := t.Validator.Struct(request)
	if badRequest != nil {
		tx.Rollback()
		errService = exception.NewError(badRequest, exception.ErrorBadRequest)
		return
	}
	errConflict := t.Repository.CreateTodoList(ctx, tx, *request.ToTodoList(userID))
	if errConflict != nil {
		tx.Rollback()
		errService = exception.NewError(errConflict, exception.ErrorConflict)
		return
	}
	tx.Commit()
	return
}

func (t *TodoListService) CreatesTodoLists(ctx context.Context, requests model.TodoListRequests, params web.Params) (errService error) {
	tx := t.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err := r.(error)
			errService = exception.NewError(err, exception.ErrorInternalServer)
		}
	}()
	userID, validUUID := params.UserID.ToUUID()
	if validUUID != nil {
		tx.Rollback()
		errService = exception.NewError(validUUID, exception.ErrorBadRequest)
		return
	}
	requestValidation := model.TodoListRequestsValidation{TodoList: requests}
	badRequest := t.Validator.Struct(requestValidation)
	if badRequest != nil {
		tx.Rollback()
		errService = exception.NewError(badRequest, exception.ErrorBadRequest)
		return
	}
	errConflict := t.Repository.CreateTodoLists(ctx, tx, requests.ToTodoLists(userID))
	if errConflict != nil {
		tx.Rollback()
		errService = exception.NewError(errConflict, exception.ErrorConflict)
		return
	}
	tx.Commit()
	return
}

func (t *TodoListService) UpdateTodoList(ctx context.Context, request model.TodoListRequest, params web.Params) (errService error) {
	tx := t.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err := r.(error)
			errService = exception.NewError(err, exception.ErrorInternalServer)
		}
	}()
	if params.Query == nil {
		tx.Rollback()
		errService = exception.NewError(fmt.Errorf("error query params is nil"), exception.ErrorInternalServer)
		return
	}
	badRequest := t.Validator.Struct(request)
	if badRequest != nil {
		tx.Rollback()
		errService = exception.NewError(badRequest, exception.ErrorBadRequest)
		return
	}
	userID, validUUID := params.UserID.ToUUID()
	queryParams, ok := params.Query.(web.TodoListByIDQuery)
	if !ok {
		tx.Rollback()
		errService = exception.NewError(fmt.Errorf("error parsing get by id query params"), exception.ErrorInternalServer)
		return
	}

	if validUUID != nil {
		tx.Rollback()
		errService = exception.NewError(validUUID, exception.ErrorBadRequest)
		return
	}

	validationQueryParams := t.Validator.Struct(queryParams)
	if validationQueryParams != nil {
		tx.Rollback()
		errService = exception.NewError(validationQueryParams, exception.ErrorBadRequest)
		return
	}
	value, errParsing := queryParams.ToValue()
	if errParsing != nil {
		tx.Rollback()
		errService = exception.NewError(errParsing, exception.ErrorInternalServer)
		return
	}
	exist := t.Repository.TodoListExistByID(ctx, tx, value.ID, userID)
	if !exist {
		tx.Rollback()
		errService = exception.NewError(fmt.Errorf("todolist with id %v not found", value.ID), exception.ErrorNotFound)
		return
	}
	t.Repository.UpdateTodoListByID(ctx, tx, *request.ToTodoList(userID), value.ID, userID)
	tx.Commit()
	return
}

func (t *TodoListService) DeleteTodoList(ctx context.Context, params web.Params) (errService error) {
	tx := t.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err := r.(error)
			errService = exception.NewError(err, exception.ErrorInternalServer)
		}
	}()
	if params.Query == nil {
		tx.Rollback()
		errService = exception.NewError(fmt.Errorf("error query params is nil"), exception.ErrorInternalServer)
		return
	}
	userID, validUUID := params.UserID.ToUUID()
	queryParams, ok := params.Query.(web.TodoListByIDQuery)
	if !ok {
		tx.Rollback()
		errService = exception.NewError(fmt.Errorf("error parsing get by id query params"), exception.ErrorInternalServer)
		return
	}

	if validUUID != nil {
		tx.Rollback()
		errService = exception.NewError(validUUID, exception.ErrorBadRequest)
		return
	}

	validationQueryParams := t.Validator.Struct(queryParams)
	if validationQueryParams != nil {
		tx.Rollback()
		errService = exception.NewError(validationQueryParams, exception.ErrorBadRequest)
		return
	}
	value, errParsing := queryParams.ToValue()
	if errParsing != nil {
		tx.Rollback()
		errService = exception.NewError(errParsing, exception.ErrorInternalServer)
		return
	}
	exist := t.Repository.TodoListExistByID(ctx, tx, value.ID, userID)
	if !exist {
		tx.Rollback()
		errService = exception.NewError(fmt.Errorf("todolist with id %v not found", value.ID), exception.ErrorNotFound)
		return
	}
	t.Repository.DeleteTodoListByID(ctx, tx, value.ID, userID)
	tx.Commit()
	return
}

func (t *TodoListService) DeletesTodoLists(ctx context.Context, params web.Params) (errService error) {
	tx := t.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err := r.(error)
			errService = exception.NewError(err, exception.ErrorInternalServer)
		}
	}()
	if params.Query == nil {
		tx.Rollback()
		errService = exception.NewError(fmt.Errorf("error query params is nil"), exception.ErrorInternalServer)
		return
	}
	userID, validUUID := params.UserID.ToUUID()
	queryParams, ok := params.Query.(web.TodoListByIDsQuery)
	if !ok {
		tx.Rollback()
		errService = exception.NewError(fmt.Errorf("error parsing get by id query params"), exception.ErrorInternalServer)
		return
	}

	if validUUID != nil {
		tx.Rollback()
		errService = exception.NewError(validUUID, exception.ErrorBadRequest)
		return
	}

	validationQueryParams := t.Validator.Struct(queryParams)
	if validationQueryParams != nil {
		tx.Rollback()
		errService = exception.NewError(validationQueryParams, exception.ErrorBadRequest)
		return
	}
	value, errParsing := queryParams.ToValue()
	if errParsing != nil {
		tx.Rollback()
		errService = exception.NewError(errParsing, exception.ErrorInternalServer)
		return
	}
	exist := t.Repository.TodoListsExistByIDs(ctx, tx, value.IDs, userID)
	if !exist {
		tx.Rollback()
		errService = exception.NewError(fmt.Errorf("todolist with id %v not found", value.IDs), exception.ErrorNotFound)
		return
	}
	t.Repository.DeleteTodoListsByIDs(ctx, tx, value.IDs, userID)
	tx.Commit()
	return
}

func (t *TodoListService) FindTodoListsBySearch(ctx context.Context, params web.Params) (responses model.TodoListResponses, pagination web.Pagination, errService error) {
	tx := t.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err := r.(error)
			errService = exception.NewError(err, exception.ErrorInternalServer)
		}
	}()
	if params.Query == nil {
		tx.Rollback()
		errService = exception.NewError(fmt.Errorf("error query params is nil"), exception.ErrorInternalServer)
		return
	}
	userID, validUUID := params.UserID.ToUUID()
	queryParams, ok := params.Query.(web.SearchQuery)
	if !ok {
		tx.Rollback()
		errService = exception.NewError(fmt.Errorf("error parsing search query"), exception.ErrorInternalServer)
		return
	}
	if validUUID != nil {
		tx.Rollback()
		errService = exception.NewError(validUUID, exception.ErrorBadRequest)
		return
	}
	validationQueryParams := t.Validator.Struct(queryParams)
	if validationQueryParams != nil {
		tx.Rollback()
		errService = exception.NewError(validationQueryParams, exception.ErrorBadRequest)
		return
	}

	value, errParsing := queryParams.ToValue()
	if errParsing != nil {
		tx.Rollback()
		errService = exception.NewError(errParsing, exception.ErrorInternalServer)
		return
	}
	responses = t.Repository.GetTodoListsSearch(ctx, tx, *value, userID).ToTodoListResponses()
	var totalData int64
	valueSearch := []string{"%", value.Search, "%"}
	key := strings.Join(valueSearch, "")
	errCount := tx.WithContext(ctx).Model(&model.TodoList{}).Where("user_id = ?", userID).Where("task_name LIKE ?", key).Or("description LIKE ?", key).Count(&totalData).Error
	if errCount != nil {
		tx.Rollback()
		errService = exception.NewError(errCount, exception.ErrorInternalServer)
		return
	}
	tx.Commit()
	totalPage := int(math.Ceil(float64(totalData) / float64(config.Other.Limit)))
	pagination = web.Pagination{
		Next:      value.Page + 1,
		Current:   value.Page,
		Previous:  value.Page - 1,
		TotalPage: totalPage,
		Data:      int(totalData),
	}
	return
}

func (t *TodoListService) FindTodoLists(ctx context.Context, params web.Params) (responses model.TodoListResponses, pagination web.Pagination, errService error) {
	tx := t.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err := r.(error)
			errService = exception.NewError(err, exception.ErrorInternalServer)
		}
	}()
	if params.Query == nil {
		tx.Rollback()
		errService = exception.NewError(fmt.Errorf("error query params is nil"), exception.ErrorInternalServer)
		return
	}
	userID, validUUID := params.UserID.ToUUID()
	queryParams, ok := params.Query.(web.GetAllQuery)
	if !ok {
		tx.Rollback()
		errService = exception.NewError(fmt.Errorf("error parsing get all query params"), exception.ErrorInternalServer)
		return
	}

	if validUUID != nil {
		tx.Rollback()
		errService = exception.NewError(validUUID, exception.ErrorBadRequest)
		return
	}

	validationQueryParams := t.Validator.Struct(queryParams)
	if validationQueryParams != nil {
		tx.Rollback()
		errService = exception.NewError(validationQueryParams, exception.ErrorBadRequest)
		return
	}

	value, errParsing := queryParams.ToValue()
	if errParsing != nil {
		tx.Rollback()
		errService = exception.NewError(errParsing, exception.ErrorInternalServer)
		return
	}
	responses = t.Repository.GetTodoLists(ctx, tx, *value, userID).ToTodoListResponses()
	var totalData int64
	errCount := tx.WithContext(ctx).Model(&model.TodoList{}).Where("user_id = ?", userID).Count(&totalData).Error
	if errCount != nil {
		tx.Rollback()
		errService = exception.NewError(errCount, exception.ErrorInternalServer)
		return
	}
	tx.Commit()
	totalPage := int(math.Ceil(float64(totalData) / float64(config.Other.Limit)))
	pagination = web.Pagination{
		Next:      value.Page + 1,
		Current:   value.Page,
		Previous:  value.Page - 1,
		TotalPage: totalPage,
		Data:      int(totalData),
	}
	return
}

func (t *TodoListService) FindTodoListByID(ctx context.Context, params web.Params) (response model.TodoListResponse, errService error) {
	tx := t.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err := r.(error)
			errService = exception.NewError(err, exception.ErrorInternalServer)
		}
	}()
	if params.Query == nil {
		tx.Rollback()
		errService = exception.NewError(fmt.Errorf("error query params is nil"), exception.ErrorInternalServer)
		return
	}
	userID, validUUID := params.UserID.ToUUID()
	queryParams, ok := params.Query.(web.TodoListByIDQuery)
	if !ok {
		tx.Rollback()
		errService = exception.NewError(fmt.Errorf("error parsing get by id query params"), exception.ErrorInternalServer)
		return
	}

	if validUUID != nil {
		tx.Rollback()
		errService = exception.NewError(validUUID, exception.ErrorBadRequest)
		return
	}

	validationQueryParams := t.Validator.Struct(queryParams)
	if validationQueryParams != nil {
		tx.Rollback()
		errService = exception.NewError(validationQueryParams, exception.ErrorBadRequest)
		return
	}
	value, errParsing := queryParams.ToValue()
	if errParsing != nil {
		tx.Rollback()
		errService = exception.NewError(errParsing, exception.ErrorInternalServer)
		return
	}
	todolist, err := t.Repository.GetTodoListByID(ctx, tx, value.ID, userID)
	if err != nil {
		tx.Rollback()
		errService = exception.NewError(err, exception.ErrorNotFound)
		return
	}
	tx.Commit()
	response = *todolist.ToTodoListResponse()
	return
}
