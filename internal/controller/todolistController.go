package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"go_gin/internal/domain/model"
	"go_gin/internal/domain/model/web"
	"go_gin/internal/handler"
	"net/http"
)

type TodoListController struct {
	Service model.TodoListService
}

func (t *TodoListController) UpdateTodoList(c *gin.Context) {
	var request model.TodoListRequest
	ctx := context.Background()
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}
	userID := c.Param("id")
	params := web.Params{UserID: web.UserID(userID)}
	errService := t.Service.UpdateTodoList(ctx, request, params)
	if errService != nil {
		responseErrors := handler.NewResponseErrors(errService)
		c.JSON(responseErrors.Status, responseErrors)
		return
	}
	c.JSON(http.StatusOK, web.NewStandartResponse(http.StatusOK, "successfuly update todo list", nil))
}

func (t *TodoListController) DeleteTodoList(c *gin.Context) {
	var query web.TodoListByIDQuery
	ctx := context.Background()
	err := c.ShouldBindQuery(&query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "query params invalid",
		})
		return
	}
	userID := c.Param("id")
	params := web.Params{
		UserID: web.UserID(userID),
		Query:  query,
	}
	errService := t.Service.DeleteTodoList(ctx, params)
	if errService != nil {
		responseErrors := handler.NewResponseErrors(errService)
		c.JSON(responseErrors.Status, responseErrors)
		return
	}
	c.JSON(http.StatusOK, web.NewStandartResponse(http.StatusOK, "successfuly delete todo list", nil))
}

func (t *TodoListController) DeleteTodoLists(c *gin.Context) {
	ids := c.QueryArray("id")
	userID := c.Param("id")
	ctx := context.Background()
	if len(ids) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "query params invalid",
		})
		return
	}
	query := web.TodoListByIDsQuery{IDs: ids}
	params := web.Params{UserID: web.UserID(userID), Query: query}
	errService := t.Service.DeletesTodoLists(ctx, params)
	if errService != nil {
		responseErrors := handler.NewResponseErrors(errService)
		c.JSON(responseErrors.Status, responseErrors)
		return
	}
	c.JSON(http.StatusOK, web.NewStandartResponse(http.StatusOK, "successfuly delete todo list", nil))
}

func NewTodoListController(service model.TodoListService) *TodoListController {
	return &TodoListController{Service: service}
}

func (t *TodoListController) GetTodoListSearch(c *gin.Context) {
	var query web.SearchQuery
	ctx := context.Background()
	err := c.ShouldBindQuery(&query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "query params invalid",
		})
		return
	}
	userID := c.Param("id")
	params := web.Params{
		UserID: web.UserID(userID),
		Query:  query,
	}
	responses, pagination, errService := t.Service.FindTodoListsBySearch(ctx, params)
	if errService != nil {
		responseErrors := handler.NewResponseErrors(errService)
		c.JSON(responseErrors.Status, responseErrors)
		return
	}
	c.JSON(http.StatusOK, web.NewStandartResponse(http.StatusOK, "Succesfuly Get Search", map[string]interface{}{
		"todolist":   responses,
		"pagination": pagination,
	}))
}
func (t *TodoListController) GetTodoListAll(c *gin.Context) {
	var query web.GetAllQuery
	ctx := context.Background()
	err := c.ShouldBindQuery(&query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "query params invalid",
		})
		return
	}
	userID := c.Param("id")
	params := web.Params{
		UserID: web.UserID(userID),
		Query:  query,
	}
	responses, pagination, errService := t.Service.FindTodoLists(ctx, params)
	if errService != nil {
		responseErrors := handler.NewResponseErrors(errService)
		c.JSON(responseErrors.Status, responseErrors)
		return
	}
	c.JSON(http.StatusOK, web.NewStandartResponse(http.StatusOK, "successfuly get all", map[string]interface{}{
		"todolist":   responses,
		"pagination": pagination,
	}))
}

func (t *TodoListController) GetTodoListByID(c *gin.Context) {
	var query web.TodoListByIDQuery
	ctx := context.Background()
	err := c.ShouldBindQuery(&query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "query params invalid",
		})
		return
	}
	userID := c.Param("id")
	params := web.Params{
		UserID: web.UserID(userID),
		Query:  query,
	}
	response, errService := t.Service.FindTodoListByID(ctx, params)
	if errService != nil {
		responseErrors := handler.NewResponseErrors(errService)
		c.JSON(responseErrors.Status, responseErrors)
		return
	}
	c.JSON(http.StatusOK, web.NewStandartResponse(http.StatusOK, "successfuly get all", map[string]interface{}{
		"todolist": response,
	}))
}

func (t *TodoListController) CreateTodoList(c *gin.Context) {
	var request model.TodoListRequest
	ctx := context.Background()
	c.ShouldBindJSON(&request)

	userID := c.Param("id")
	params := web.Params{UserID: web.UserID(userID)}
	errService := t.Service.CreateTodoList(ctx, request, params)
	if errService != nil {
		responseErrors := handler.NewResponseErrors(errService)
		c.JSON(responseErrors.Status, responseErrors)
		return
	}
	c.JSON(http.StatusOK, web.NewStandartResponse(http.StatusOK, "successfuly create todo list", nil))
}

func (t *TodoListController) CreatesTodoLists(c *gin.Context) {
	var requests model.TodoListRequests
	ctx := context.Background()
	err := c.ShouldBindJSON(&requests)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}
	userID := c.Param("id")
	params := web.Params{UserID: web.UserID(userID)}
	errService := t.Service.CreatesTodoLists(ctx, requests, params)
	if errService != nil {
		responseErrors := handler.NewResponseErrors(errService)
		c.JSON(responseErrors.Status, responseErrors)
		return
	}
	c.JSON(http.StatusOK, web.NewStandartResponse(http.StatusOK, "successfuly creates todo list", nil))
}
