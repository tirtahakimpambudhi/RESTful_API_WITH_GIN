package controller

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go_gin/internal/domain/model"
	"go_gin/internal/domain/model/web"
	"go_gin/internal/exception"
	"go_gin/internal/handler"
	"go_gin/pkg/helper"
	"net/http"
)

type UsersController struct {
	Service model.UsersService
}

func (u *UsersController) GetBySearch(c *gin.Context) {
	var queryParams web.SearchQuery
	ctx := context.Background()
	c.ShouldBindQuery(&queryParams)

	users, pagination, err := u.Service.FindUsersBySearch(ctx, queryParams)
	if err != nil {
		responseErrors := handler.NewResponseErrors(err)
		c.JSON(responseErrors.Status, responseErrors)
		return
	}
	c.JSON(http.StatusOK, web.NewStandartResponse(http.StatusOK, "Succesefully Get All", map[string]interface{}{
		"users":      users,
		"pagination": pagination,
	}))
}

func (u *UsersController) RestoreUserByID(c *gin.Context) {
	ctx := context.Background()
	IDString := c.Param("id")
	ID, err := uuid.Parse(IDString)
	badFormatErrorUUID := helper.NewCustomError(err, exception.ErrorBadRequest)

	if badFormatErrorUUID != nil {
		responseErrors := handler.NewResponseErrors(badFormatErrorUUID)
		c.JSON(responseErrors.Status, responseErrors)
		return
	}
	err = u.Service.RestoreUserByID(ctx, ID)
	if err != nil {
		responseErrors := handler.NewResponseErrors(err)
		c.JSON(responseErrors.Status, responseErrors)
		return
	}
	c.JSON(http.StatusOK, web.NewStandartResponse(http.StatusOK, "Succesefully Delete", nil))
}

func (u *UsersController) RestoreUsersByIDs(c *gin.Context) {
	var IDs []uuid.UUID
	ctx := context.Background()
	IDsString := c.QueryArray("id")
	if len(IDsString) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "query params invalid",
		})
		return
	}
	for _, s := range IDsString {
		ID, err := uuid.Parse(s)
		badFormatErrorUUID := helper.NewCustomError(err, exception.ErrorBadRequest)
		if badFormatErrorUUID != nil {
			responseErrors := handler.NewResponseErrors(badFormatErrorUUID)
			c.JSON(responseErrors.Status, responseErrors)
			return
		}
		IDs = append(IDs, ID)
	}
	err := u.Service.RestoreUsersByIDs(ctx, IDs)
	if err != nil {
		responseErrors := handler.NewResponseErrors(err)
		c.JSON(responseErrors.Status, responseErrors)
		return
	}
	c.JSON(http.StatusOK, web.NewStandartResponse(http.StatusOK, "Succesefully Deletes", nil))
}

func NewUsersController(service model.UsersService) model.UsersController {
	return &UsersController{Service: service}
}

func (u *UsersController) GetAll(c *gin.Context) {
	var queryParams web.GetAllQuery
	ctx := context.Background()
	c.ShouldBindQuery(&queryParams)

	users, pagination, err := u.Service.FindUsers(ctx, queryParams)
	if err != nil {
		responseErrors := handler.NewResponseErrors(err)
		c.JSON(responseErrors.Status, responseErrors)
		return
	}
	c.JSON(http.StatusOK, web.NewStandartResponse(http.StatusOK, "Succesefully Get All", map[string]interface{}{
		"users":      users,
		"pagination": pagination,
	}))
}

func (u *UsersController) GetByID(c *gin.Context) {
	IDString := c.Param("id")
	ID, err := uuid.Parse(IDString)
	badFormatErrorUUID := helper.NewCustomError(err, exception.ErrorBadRequest)

	if badFormatErrorUUID != nil {
		responseErrors := handler.NewResponseErrors(badFormatErrorUUID)
		c.JSON(responseErrors.Status, responseErrors)
		return
	}
	user, err := u.Service.FindUserByID(c, ID)
	if err != nil {
		responseErrors := handler.NewResponseErrors(err)
		c.JSON(responseErrors.Status, responseErrors)
		return
	}
	c.JSON(http.StatusOK, web.NewStandartResponse(http.StatusOK, "Succesefully Get By ID", map[string]interface{}{
		"user": user,
	}))
}

func (u *UsersController) CreateUser(c *gin.Context) {
	var userRequest model.UserRequest
	var ID uuid.UUID

	c.ShouldBindJSON(&userRequest)
	ctx := context.Background()
	uid := uuid.New()
	userRequest.ID = uid
	ID = uid
	err := u.Service.CreateUser(ctx, userRequest)

	if err != nil {
		responseErrors := handler.NewResponseErrors(err)
		c.JSON(responseErrors.Status, responseErrors)
		return
	}
	c.JSON(http.StatusOK, web.NewStandartResponse(http.StatusOK, "Succesefully Create", map[string]interface{}{
		"id": ID,
	}))

}

func (u *UsersController) LoginUser(c *gin.Context) {
	var userRequest model.UserLoginUpdateRequest

	c.ShouldBindJSON(&userRequest)
	ctx := context.Background()

	accessToken, refreshToken, err := u.Service.LoginUsers(ctx, userRequest)
	cookie := &http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		HttpOnly: true,
		Path:     "/",
	}
	if err != nil {
		responseErrors := handler.NewResponseErrors(err)
		c.JSON(responseErrors.Status, responseErrors)
		return
	}
	http.SetCookie(c.Writer, cookie)
	c.JSON(http.StatusOK, web.NewStandartResponse(http.StatusOK, "Successfully Login", map[string]interface{}{
		"accessToken": accessToken,
	}))
}

func (u *UsersController) LogoutUser(c *gin.Context) {
	ctx := context.Background()
	cookieValue, err := c.Cookie("refreshToken")
	if err != nil {
		c.JSON(http.StatusPermanentRedirect, gin.H{
			"status":  http.StatusPermanentRedirect,
			"message": "No Cookie",
		})
	}
	errService := u.Service.LogoutUsers(ctx, cookieValue)
	if errService != nil {
		responseErrors := handler.NewResponseErrors(err)
		c.JSON(responseErrors.Status, "/api")
		return
	}
	cookie := &http.Cookie{
		Name:     "refreshToken",
		MaxAge:   -1,
		HttpOnly: true,
		Value:    "",
		Path:     "/",
	}
	http.SetCookie(c.Writer, cookie)
	c.JSON(http.StatusOK, web.NewStandartResponse(http.StatusOK, "Successfully Refresh Token", map[string]interface{}{
		"data": nil,
	}))
}
func (u *UsersController) RefreshTokenUser(c *gin.Context) {
	cookie, err := c.Cookie("refreshToken")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "No Cookie",
		})
		return
	}
	ctx := context.Background()
	accessToken, err := u.Service.RefreshTokenUser(ctx, cookie)
	if err != nil {
		responseErrors := handler.NewResponseErrors(err)
		c.JSON(responseErrors.Status, responseErrors)
		return
	}
	c.JSON(http.StatusOK, web.NewStandartResponse(http.StatusOK, "Successfully Refresh Token", map[string]interface{}{
		"accessToken": accessToken,
	}))
}

func (u *UsersController) CreateUsers(c *gin.Context) {
	var usersRequests model.UsersRequests
	var IDs []uuid.UUID

	ctx := context.Background()
	c.ShouldBindJSON(&usersRequests)
	for i, _ := range usersRequests {
		id := uuid.New()
		usersRequests[i].ID = id
		IDs = append(IDs, id)
	}
	err := u.Service.CreateUsers(ctx, usersRequests)
	if err != nil {
		responseErrors := handler.NewResponseErrors(err)
		c.JSON(responseErrors.Status, responseErrors)
		return
	}
	c.JSON(http.StatusOK, web.NewStandartResponse(http.StatusOK, "Succesefully Creates", map[string]interface{}{
		"ids": IDs,
	}))
}

func (u *UsersController) UpdateUserID(c *gin.Context) {
	var userRequest model.UserLoginUpdateRequest
	ctx := context.Background()
	IDString := c.Param("id")
	ID, err := uuid.Parse(IDString)
	c.ShouldBindJSON(&userRequest)
	badFormatErrorUUID := helper.NewCustomError(err, exception.ErrorBadRequest)

	if badFormatErrorUUID != nil {
		responseErrors := handler.NewResponseErrors(badFormatErrorUUID)
		c.JSON(responseErrors.Status, responseErrors)
		return
	}
	err = u.Service.UpdateUserID(ctx, userRequest, ID)
	fmt.Println(err, ID, userRequest)
	if err != nil {
		responseErrors := handler.NewResponseErrors(err)
		c.JSON(responseErrors.Status, responseErrors)
		return
	}
	c.JSON(http.StatusOK, web.NewStandartResponse(http.StatusOK, "Succesefully Update", map[string]interface{}{
		"id": ID,
	}))

}

func (u *UsersController) DeleteUserByID(c *gin.Context) {
	ctx := context.Background()
	IDString := c.Param("id")
	ID, err := uuid.Parse(IDString)
	badFormatErrorUUID := helper.NewCustomError(err, exception.ErrorBadRequest)

	if badFormatErrorUUID != nil {
		responseErrors := handler.NewResponseErrors(badFormatErrorUUID)
		c.JSON(responseErrors.Status, responseErrors)
		return
	}
	err = u.Service.DeleteUserByID(ctx, ID)
	if err != nil {
		responseErrors := handler.NewResponseErrors(err)
		c.JSON(responseErrors.Status, responseErrors)
		return
	}
	c.JSON(http.StatusOK, web.NewStandartResponse(http.StatusOK, "Succesefully Delete", nil))
}

func (u *UsersController) DeleteUsersByIDs(c *gin.Context) {
	var IDs []uuid.UUID
	ctx := context.Background()

	IDsString := c.QueryArray("id")
	if len(IDsString) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "query params invalid",
		})
		return
	}
	for _, s := range IDsString {
		ID, err := uuid.Parse(s)
		badFormatErrorUUID := helper.NewCustomError(err, exception.ErrorBadRequest)
		if badFormatErrorUUID != nil {
			responseErrors := handler.NewResponseErrors(badFormatErrorUUID)
			c.JSON(responseErrors.Status, responseErrors)
			return
		}
		IDs = append(IDs, ID)
	}
	err := u.Service.DeleteUsersByIDs(ctx, IDs)
	if err != nil {
		responseErrors := handler.NewResponseErrors(err)
		c.JSON(responseErrors.Status, responseErrors)
		return
	}
	c.JSON(http.StatusOK, web.NewStandartResponse(http.StatusOK, "Succesefully Deletes", nil))
}
