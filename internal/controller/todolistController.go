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

// UpdateTodolist	godoc
// @Summary	Update Todolist
// @Description Update Todolist as JSON
// @Tags Todolist
// @Param id	path	string	true "Must Be UUID Format"
// @Param id	query	int	true "ID Todolist"
// @Param request	body	model.TodoListRequest	true	"Object Todolist for Update Todolist"
// @Produce	json
// @Success	200	{object}	web.StandartResponse
// @Failure 400 {object} 	handler.ResponseErrors "Bad request"
// @Failure 401 {object} 	handler.ResponseErrors "Unauthorized"
// @Failed	404	{object} 	handler.ResponseErrors "Not Found"
// @Router  /user/{id}/todolist [put]
func (t *TodoListController) UpdateTodoList(c *gin.Context) {
	var request model.TodoListRequest
	var query web.TodoListByIDQuery
	errQuery := c.ShouldBindQuery(&query)
	ctx := context.Background()
	if errQuery != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "query params invalid",
		})
		return
	}
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}
	userID := c.Param("id")
	params := web.Params{UserID: web.UserID(userID), Query: query}
	errService := t.Service.UpdateTodoList(ctx, request, params)
	if errService != nil {
		responseErrors := handler.NewResponseErrors(errService)
		c.JSON(responseErrors.Status, responseErrors)
		return
	}
	c.JSON(http.StatusOK, web.NewStandartResponse(http.StatusOK, "successfuly update todo list", nil))
}

// DeleteTodolistByID godoc
// @Summary	Delete Todolist By ID
// @Description Retrieve a object Todolist as JSON
// @Tags Todolist
// @Param id	path	string	true "Must Be UUID Format"
// @Param id 		query	int		true "ID todolist"
// @Produce	json
// @Success	200	{object}	web.StandartResponse
// @Failure 400 {object} 	handler.ResponseErrors "Bad request"
// @Failure 401 {object} 	handler.ResponseErrors "Unauthorized"
// @Failed	404	{object} 	handler.ResponseErrors "Not Found"
// @Router  /user/{id}/todolist [delete]
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

// DeleteTodolistByID godoc
// @Summary	Delete Todolist By ID
// @Description Retrieve a object Todolist as JSON
// @Tags Todolist
// @Param id	path	string	true "Must Be UUID Format"
// @Param id 		query	[]int		true "ID todolist"
// @Produce	json
// @Success	200	{object}	web.StandartResponse
// @Failure 400 {object} 	handler.ResponseErrors "Bad request"
// @Failure 401 {object} 	handler.ResponseErrors "Unauthorized"
// @Failed	404	{object} 	handler.ResponseErrors "Not Found"
// @Router  /user/{id}/todolists [delete]
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

// Todolist godoc
// @Summary Get Todolist array by search key
// @Description Retrieve a list of all Todolist as JSON
// @Tags Todolist
// @Param id	path	string 	true "Must be UUID Format"
// @Param search query string true "search keywords for task name, description"
// @Param page query int true "Page number"
// @Produce json
// @Success 200 {object} web.StandartResponse
// @Failure 400 {object} handler.ResponseErrors "Bad request"
// @Failure 401 {object} handler.ResponseErrors "Unauthorized"
// @Failed	404	{object}	handler.ResponseErrors "Not Found"
// @Router /user/{id}/todolists/s [get]
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

// GetAllTodolist godoc
// @Summary Get all Todolist array
// @Description Retrieve a list all Todolist as JSON
// @Tags Todolist
// @Param id	path	string 	true	"Must Be UUID Format"
// @Param page	query	int		true	"Page Number"
// @Produce json
// @Success	200 {object}	web.StandartResponse
// @Failure 400 {object} handler.ResponseErrors "Bad request"
// @Failure 401 {object} handler.ResponseErrors "Unauthorized"
// @Failed	404	{object}	handler.ResponseErrors "Not Found"
// @Router /user/{id}/todolists	[get]
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

// GetTodolistByID godoc
// @Summary	Get Todolist By ID
// @Description Retrieve a object Todolist as JSON
// @Tags Todolist
// @Param id	path	string	true "Must Be UUID Format"
// @Param id 		query	int		true "ID todolist"
// @Produce	json
// @Success	200	{object}	web.StandartResponse
// @Failure 400 {object} 	handler.ResponseErrors "Bad request"
// @Failure 401 {object} 	handler.ResponseErrors "Unauthorized"
// @Failed	404	{object} 	handler.ResponseErrors "Not Found"
// @Router  /user/{id}/todolist [get]
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

// CreateTodolist	godoc
// @Summary	Create New Todolist
// @Description Create New Todolist as JSON
// @Tags Todolist
// @Param id	path	string	true "Must Be UUID Format"
// @Param request	body	model.TodoListRequest	true  "Object Todolist for Create Todolist"
// @Produce	json
// @Success	200	{object}	web.StandartResponse
// @Failure 400 {object} 	handler.ResponseErrors "Bad request"
// @Failure 401 {object} 	handler.ResponseErrors "Unauthorized"
// @Failed	404	{object} 	handler.ResponseErrors "Not Found"
// @Router  /user/{id}/todolist [post]
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

// CreateTodolists	godoc
// @Summary	Create New Todolists
// @Description Create New Todolists as JSON
// @Tags Todolist
// @Param id	path	string	true "Must Be UUID Format"
// @Param request	body	model.TodoListRequests	true	"Array Todolist for Create Many"
// @Produce	json
// @Success	200	{object}	web.StandartResponse
// @Failure 400 {object} 	handler.ResponseErrors "Bad request"
// @Failure 401 {object} 	handler.ResponseErrors "Unauthorized"
// @Failed	404	{object} 	handler.ResponseErrors "Not Found"
// @Router  /user/{id}/todolists [post]
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
