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

// GetAll godoc
// @Security Bearer
// @Summary Get Users array
// @Description Retrieve a list of all users as JSON
// @Tags Admin
// @Param page query int true "Page number"
// @Produce json
// @Success 200 {object} web.StandartResponse
// @Failure 400 {object} handler.ResponseErrors "Bad request"
// @Failure 401 {object} handler.ResponseErrors "Unauthorized"
// @Failure 403 {object} handler.ResponseErrors "Forbidden"
// @Router /admin/users [get]
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

// GetByID godoc
// @Security Bearer
// @Summary Get User By ID for ADMIN
// @Description Retrieve user details by ID
// @Tags Admin
// @Param id path string true "Must be in UUID format"
// @Produce json
// @Success 200 {object} web.StandartResponse
// @Failure 400 {object} handler.ResponseErrors "Bad request format"
// @Failure 404 {object} handler.ResponseErrors "User not found"
// @Router /admin/user/{id} [get]
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

// GetUserSearch godoc
// @Security Bearer
// @Summary Get Users array
// @Description Retrieve a list of users as JSON based on search criteria
// @Tags Admin
// @Param search query string true "Search users by key"
// @Param page query int true "Page number"
// @Produce json
// @Success 200 {object} web.StandartResponse
// @Failure 400 {object} handler.ResponseErrors "Invalid query parameters"
// @Failure 401 {object} handler.ResponseErrors "Unauthorized"
// @Failure 403 {object} handler.ResponseErrors "Forbidden"
// @Router /admin/users/search [get]
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

// CreateUsers godoc
// @Security Bearer
// @Summary Register for all roles
// @Description Create new many users
// @Tags Admin
// @Param request body model.UsersRequests true "Registers Request Body"
// @Produce json
// @Success 200 {object} web.StandartResponse
// @Failure 400 {object} handler.ResponseErrors "Bad request format"
// @Failure 401 {object} handler.ResponseErrors "Unauthorized"
// @Failure 409 {object} handler.ResponseErrors "Conflict"
// @Router /admin/registers [post]
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

// RestoreUserByID godoc
// @Security Bearer
// @Summary Restore User for ADMIN role
// @Description Restore a user by ID
// @Tags Admin
// @Param id path string true "Must be in UUID format"
// @Produce json
// @Success 200 {object} web.StandartResponse
// @Failure 400 {object} handler.ResponseErrors "Bad request format"
// @Failure 401 {object} handler.ResponseErrors "Unauthorized"
// @Failure 404 {object} handler.ResponseErrors "User not found"
// @Router /admin/user/{id} [patch]
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

// RestoreUsersByIDs godoc
// @Security Bearer
// @Summary Restore Users for ADMIN role
// @Description Restore users by IDs
// @Tags Admin
// @Param id query []string true "Must be in UUID format"
// @Produce json
// @Success 200 {object} web.StandartResponse
// @Failure 400 {object} handler.ResponseErrors "Bad request format"
// @Failure 401 {object} handler.ResponseErrors "Unauthorized"
// @Failure 404 {object} handler.ResponseErrors "User not found"
// @Router /admin/users [patch]
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

// RegisterUsers	godoc
// CreateUser godoc
// @Summary Register for all roles
// @Description Create new users
// @Tags All
// @Param request body model.UserRequest true "Register Request Body"
// @Produce json
// @Success 200 {object} web.StandartResponse
// @Failure 400 {object} handler.ResponseErrors "Bad request format"
// @Failure 409 {object} handler.ResponseErrors "Conflict"
// @Router /register [post]
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

// RefreshTokenUser godoc
// @Summary Get Refresh Token
// @Description Responds with a Refresh Token
// @Tags All
// @Produce json
// @Success 200 {object} web.StandartResponse
// @Failure 400 {object} handler.ResponseErrors "Bad request format"
// @Failure 401 {object} handler.ResponseErrors "Unauthorized"
// @Failure 403 {object} handler.ResponseErrors "Forbidden"
// @Router /refresh [get]
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

// LoginUser godoc
// @Summary Login for all roles
// @Description Responds with the access token
// @Tags All
// @Param request body model.UserLoginUpdateRequest true "Login Request Body"
// @Produce json
// @Success 200 {object} web.StandartResponse
// @Failure 400 {object} handler.ResponseErrors "Bad request format"
// @Failure 404 {object} handler.ResponseErrors "User not found"
// @Router /login [post]
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

// UpdateUserID godoc
// @Summary Update User for all roles
// @Description Update user by ID
// @Tags All
// @Param request body model.UserLoginUpdateRequest true "Update Request Body"
// @Param id path string true "Must be in UUID format"
// @Produce json
// @Success 200 {object} web.StandartResponse
// @Failure 400 {object} handler.ResponseErrors "Bad request format"
// @Failure 404 {object} handler.ResponseErrors "User not found"
// @Router /user/{id} [put]
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

// LogoutUser godoc
// @Summary Logout User for all roles
// @Description Logout users
// @Tags All
// @Produce json
// @Success 200 {object} web.StandartResponse
// @Failure 400 {object} handler.ResponseErrors "Bad request format"
// @Failure 404 {object} handler.ResponseErrors "User not found"
// @Router /logout [delete]
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

// DeleteUserByID godoc
// @Security Bearer
// @Summary Delete or Banned User for ADMIN role
// @Description Delete user by ID
// @Tags Admin
// @Param id path string true "Must be in UUID format"
// @Produce json
// @Success 200 {object} web.StandartResponse
// @Failure 400 {object} handler.ResponseErrors "Bad request format"
// @Failure 401 {object} handler.ResponseErrors "Unauthorized"
// @Failure 404 {object} handler.ResponseErrors "User not found"
// @Router /admin/user/{id} [delete]
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

// DeleteUsersByIDs godoc
// @Security Bearer
// @Summary Delete or Banned User for ADMIN role
// @Description Delete users by IDs
// @Tags Admin
// @Param id query []string true "Must be in UUID format"
// @Produce json
// @Success 200 {object} web.StandartResponse
// @Failure 400 {object} handler.ResponseErrors "Bad request format"
// @Failure 401 {object} handler.ResponseErrors "Unauthorized"
// @Failure 404 {object} handler.ResponseErrors "User not found"
// @Router /admin/users [delete]
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
