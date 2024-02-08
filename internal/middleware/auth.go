package middleware

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go_gin/internal/domain/model"
	"go_gin/internal/domain/model/web"
	"go_gin/internal/exception"
	"go_gin/internal/handler"
	"go_gin/pkg/helper"
	"gorm.io/gorm"
	"net/http"
)

type Middleware struct {
	Repository model.UsersRepository
	DB         *gorm.DB
}

func (m *Middleware) IsLogin(c *gin.Context) {
	_, err := c.Cookie("refreshToken")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "UNAUTHORIZED",
		})
		c.Abort()
		return
	}

	c.Next()
}
func (m *Middleware) AuthorizationAllRole(c *gin.Context) {
	idInterface, exist := c.Get("user_id")
	tx := m.DB.Begin()
	if !exist {
		err := exception.NewError(errors.New("UNAUTHORIZATION"), exception.ErrorUnauthorized)
		responseErrors := handler.NewResponseErrors(err)
		tx.Rollback()
		c.JSON(responseErrors.Status, responseErrors)
		c.Abort()
		return
	}

	id, ok := idInterface.(uuid.UUID)

	if !ok {
		err := exception.NewError(errors.New("Bad Format UUID"), exception.ErrorBadRequest)
		responseErrors := handler.NewResponseErrors(err)
		c.JSON(responseErrors.Status, responseErrors)
		tx.Rollback()
		c.Abort()
		return
	}
	user, err := m.Repository.GetUserByID(context.Background(), tx, id)
	if err != nil {
		err := exception.NewError(err, exception.ErrorNotFound)
		responseErrors := handler.NewResponseErrors(err)
		c.JSON(responseErrors.Status, responseErrors)
		tx.Rollback()
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, web.NewStandartResponse(http.StatusOK, "successfuly authentication", user.ToUserResponse()))
}
func (m *Middleware) Authentication(c *gin.Context) {
	accessToken, err := helper.ExtractBearerToken(c.GetHeader("Authorization"))

	if err != nil {
		responseErrors := handler.NewResponseErrors(err)
		c.JSON(responseErrors.Status, responseErrors)
		c.Abort()
		return
	}

	user, errJWT := helper.VerifyAccessToken(accessToken)
	if errJWT != nil {
		responseErrors := handler.NewResponseErrors(errJWT)
		c.JSON(responseErrors.Status, responseErrors)
		c.Abort()
		return
	}
	c.Set("user_id", user.ID)
	c.Set("email", user.Email)
	c.Set("username", user.Username)
	c.Next()
}

func (m *Middleware) AuthorizationRoleAdmin(c *gin.Context) {
	idInterface, exist := c.Get("user_id")
	tx := m.DB.Begin()
	if !exist {
		err := exception.NewError(errors.New("UNAUTHORIZATION"), exception.ErrorUnauthorized)
		responseErrors := handler.NewResponseErrors(err)
		tx.Rollback()
		c.JSON(responseErrors.Status, responseErrors)
		c.Abort()
		return
	}

	id, ok := idInterface.(uuid.UUID)

	if !ok {
		err := exception.NewError(errors.New("Bad Format UUID"), exception.ErrorBadRequest)
		responseErrors := handler.NewResponseErrors(err)
		c.JSON(responseErrors.Status, responseErrors)
		tx.Rollback()
		c.Abort()
		return
	}
	user, err := m.Repository.GetUserByID(context.Background(), tx, id)

	if err != nil {
		err := exception.NewError(err, exception.ErrorNotFound)
		responseErrors := handler.NewResponseErrors(err)
		c.JSON(responseErrors.Status, responseErrors)
		tx.Rollback()
		c.Abort()
		return
	}

	if user.Roles != model.Admin {
		err := exception.NewError(errors.New("UNAUTHORIZATION"), exception.ErrorUnauthorized)
		responseErrors := handler.NewResponseErrors(err)
		tx.Rollback()
		c.JSON(responseErrors.Status, responseErrors)
		c.Abort()
		return
	}

	c.Next()
}
func (m *Middleware) AuthorizationRoleModerator(c *gin.Context) {
	idInterface, exist := c.Get("user_id")
	tx := m.DB.Begin()
	if !exist {
		err := exception.NewError(errors.New("UNAUTHORIZATION"), exception.ErrorUnauthorized)
		responseErrors := handler.NewResponseErrors(err)
		tx.Rollback()
		c.JSON(responseErrors.Status, responseErrors)
		c.Abort()
		return
	}

	id, ok := idInterface.(uuid.UUID)

	if !ok {
		err := exception.NewError(errors.New("Bad Format UUID"), exception.ErrorBadRequest)
		responseErrors := handler.NewResponseErrors(err)
		c.JSON(responseErrors.Status, responseErrors)
		tx.Rollback()
		c.Abort()
		return
	}
	user, err := m.Repository.GetUserByID(context.Background(), tx, id)

	if err != nil {
		err := exception.NewError(err, exception.ErrorNotFound)
		responseErrors := handler.NewResponseErrors(err)
		c.JSON(responseErrors.Status, responseErrors)
		tx.Rollback()
		c.Abort()
		return
	}

	if user.Roles != model.Moderator {
		err := exception.NewError(errors.New("UNAUTHORIZATION"), exception.ErrorUnauthorized)
		responseErrors := handler.NewResponseErrors(err)
		tx.Rollback()
		c.JSON(responseErrors.Status, responseErrors)
		c.Abort()
		return
	}

	c.Next()
}
