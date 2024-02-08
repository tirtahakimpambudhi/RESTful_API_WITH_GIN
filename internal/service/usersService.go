package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"go_gin/internal/config"
	"go_gin/internal/domain/model"
	"go_gin/internal/domain/model/web"
	"go_gin/internal/exception"
	"go_gin/pkg/bcrypts"
	"go_gin/pkg/helper"
	"gorm.io/gorm"
	"math"
	"strings"
	"time"
)

type UsersService struct {
	DB         *gorm.DB
	Repository model.UsersRepository
	Validation *validator.Validate
}

func (u *UsersService) FindUsersBySearch(ctx context.Context, params web.SearchQuery) (responses model.UsersResponses, pagination web.Pagination, errService error) {
	tx := u.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err := r.(error)
			errService = exception.NewError(err, exception.ErrorInternalServer)
		}
	}()
	validationQueryParams := u.Validation.Struct(params)
	if validationQueryParams != nil {
		tx.Rollback()
		errService = exception.NewError(validationQueryParams, exception.ErrorBadRequest)
		return
	}
	value, errParsing := params.ToValue()
	if errParsing != nil {
		tx.Rollback()
		errService = exception.NewError(errParsing, exception.ErrorInternalServer)
		return
	}
	responses = u.Repository.GetUsersBySearch(ctx, tx, *value).ToUsersResponses()
	valueSearch := []string{"%", value.Search, "%"}
	key := strings.Join(valueSearch, "")
	var totalData int64
	errCount := tx.Unscoped().Model(&model.User{}).Where("username LIKE ?", key).Where("email LIKE ?", key).Offset(int(value.Offset)).Limit(config.Other.Limit).Count(&totalData).Error
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

func (u *UsersService) FindUsers(ctx context.Context, params web.GetAllQuery) (responses model.UsersResponses, pagination web.Pagination, errService error) {
	tx := u.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err := r.(error)
			errService = exception.NewError(err, exception.ErrorInternalServer)
		}
	}()
	validationQueryParams := u.Validation.Struct(params)
	if validationQueryParams != nil {
		tx.Rollback()
		errService = exception.NewError(validationQueryParams, exception.ErrorBadRequest)
		return
	}
	value, errParsing := params.ToValue()
	if errParsing != nil {
		tx.Rollback()
		errService = exception.NewError(errParsing, exception.ErrorInternalServer)
		return
	}
	responses = u.Repository.GetUsers(ctx, tx, *value).ToUsersResponses()
	var totalData int64
	errCount := tx.Unscoped().Model(&model.User{}).Offset(int(value.Offset)).Limit(config.Other.Limit).Count(&totalData).Error
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

func (u *UsersService) RestoreUserByID(ctx context.Context, ID uuid.UUID) (errService error) {
	tx := u.DB.Begin()
	var count int64
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err := r.(error)
			errService = exception.NewError(err, exception.ErrorInternalServer)
		}
	}()
	errCount := tx.WithContext(ctx).Unscoped().Model(&model.User{}).Where("id = ?", ID).Count(&count).Error
	if errCount != nil {
		tx.Rollback()
		errService = exception.NewError(errCount, exception.ErrorInternalServer)
		return
	} else if count != 1 {
		tx.Rollback()
		errService = exception.NewError(fmt.Errorf("users with id %v not found", ID), exception.ErrorNotFound)
		return
	} else {
		u.Repository.RestoreUserByID(ctx, tx, ID)
		tx.Commit()
		errService = nil
		return
	}
}

func (u *UsersService) RestoreUsersByIDs(ctx context.Context, IDs []uuid.UUID) (errService error) {
	tx := u.DB.Begin()
	var count int64
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err := r.(error)
			errService = exception.NewError(err, exception.ErrorInternalServer)
		}
	}()
	errCount := tx.WithContext(ctx).Unscoped().Model(&model.User{}).Where("id IN ?", IDs).Count(&count).Error
	if errCount != nil {
		tx.Rollback()
		errService = exception.NewError(errCount, exception.ErrorInternalServer)
		return
	} else if count != int64(len(IDs)) {
		tx.Rollback()
		errService = exception.NewError(fmt.Errorf("users with id %v not found", IDs), exception.ErrorNotFound)
		return
	} else {
		u.Repository.RestoreUsersByIDs(ctx, tx, IDs)
		tx.Commit()
		errService = nil
		return
	}
}

func (u *UsersService) LogoutUsers(ctx context.Context, refreshToken string) (errService error) {
	tx := u.DB.Begin()
	var user model.Users
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err := r.(error)
			errService = exception.NewError(err, exception.ErrorInternalServer)
		}
	}()
	err := tx.Where("refresh_token = ?", refreshToken).Take(&user).Error
	if err != nil {
		errService = exception.NewError(err, exception.ErrorNotFound)
		return
	}
	errUpdate := tx.Model(&user).Update("refresh_token", nil).Error
	if errUpdate != nil {
		errService = exception.NewError(errUpdate, exception.ErrorNotFound)
		return
	}
	errService = nil
	return
}

func NewUsersService(DB *gorm.DB, repository model.UsersRepository, validate *validator.Validate) model.UsersService {
	return &UsersService{DB: DB, Repository: repository, Validation: validate}
}
func (u *UsersService) RefreshTokenUser(ctx context.Context, refreshToken string) (accessToken string, errService error) {
	tx := u.DB.Begin()

	var user model.User

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err := r.(error)
			errService = exception.NewError(err, exception.ErrorInternalServer)
		}
	}()

	// Check if there's an error retrieving the user by refresh token
	if err := tx.Where("refresh_token = ?", refreshToken).Take(&user).Error; err != nil {
		tx.Rollback()
		errService = exception.NewError(err, exception.ErrorNotFound)
		return
	}

	// Verify the refresh token
	_, errJWT := helper.VerifyRefreshToken(refreshToken)
	if errJWT != nil {
		tx.Rollback()
		errService = errJWT
		return
	}

	// Create a new access token
	claims := model.NewStandardClaimsJWT(&jwt.StandardClaims{
		Issuer:    config.JWT.AppName,
		ExpiresAt: time.Now().Add(time.Minute * time.Duration(config.JWT.Exp)).Unix(),
	}, user.ID, user.Username, user.Email)

	accessToken = helper.NewAccessToken(claims)

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		errService = exception.NewError(err, exception.ErrorInternalServer)
		return
	}

	errService = nil
	return
}

func (u *UsersService) LoginUsers(ctx context.Context, userLogin model.UserLoginUpdateRequest) (accessToken, refreshToken string, errService error) {
	tx := u.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err := r.(error)
			errService = exception.NewError(err, exception.ErrorInternalServer)
		}
	}()
	validationError := u.Validation.Struct(userLogin)
	badRequest := helper.NewCustomError(validationError, exception.ErrorBadRequest)
	if badRequest != nil {
		errService = badRequest
		return
	}
	user, err := u.Repository.GetUserByEmail(ctx, tx, userLogin.Email)
	if userLogin.Email != user.Email || userLogin.Username != user.Username || !bcrypts.CheckPasswordHash(userLogin.Password, user.Password) {
		errService = exception.NewError(errors.New("Email or Username or Password Wrong"), exception.ErrorBadRequest)
		return
	}
	claims := model.NewStandardClaimsJWT(&jwt.StandardClaims{
		Issuer:    config.JWT.AppName,
		ExpiresAt: time.Now().Add(time.Minute * time.Duration(config.JWT.Exp)).Unix(),
	}, user.ID, user.Username, user.Email)
	accessToken = helper.NewAccessToken(claims)
	refreshToken = helper.NewRefreshToken(&jwt.StandardClaims{
		Issuer:    config.JWT.AppName,
		ExpiresAt: time.Now().Add(time.Hour * time.Duration(config.JWT.Exp)).Unix(),
	})

	errServer := tx.Model(&user).Update("refresh_token", refreshToken).Error
	if err != nil {
		tx.Rollback()
		errService = exception.NewError(err, exception.ErrorNotFound)
		return
	} else if errServer != nil {
		tx.Rollback()
		errServer = exception.NewError(errServer, exception.ErrorInternalServer)
		return
	} else {
		tx.Commit()
		errService = nil
	}

	return
}

func (u *UsersService) FindUserByID(ctx context.Context, ID uuid.UUID) (userResponse model.UserResponse, errService error) {
	tx := u.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err := r.(error)
			errService = exception.NewError(err, exception.ErrorInternalServer)
			userResponse = model.UserResponse{}
		}
	}()
	user, err := u.Repository.GetUserByID(ctx, tx, ID)
	err = helper.NewCustomError(err, exception.ErrorNotFound)
	if err != nil {
		tx.Rollback()
		userResponse = model.UserResponse{}
		errService = err
		return
	} else {
		tx.Commit()
		userResponse = *user.ToUserResponse()
		errService = nil
	}
	return
}

func (u *UsersService) CreateUser(ctx context.Context, user model.UserRequest) (errService error) {
	tx := u.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err := r.(error)
			errService = exception.NewError(err, exception.ErrorInternalServer)
		}
	}()

	validationError := u.Validation.Struct(user)
	badRequest := helper.NewCustomError(validationError, exception.ErrorBadRequest)
	var adminCount int64
	errCount := tx.WithContext(ctx).Model(model.User{}).Where("roles = ?", model.Admin).Count(&adminCount).Error
	err := u.Repository.CreateUser(ctx, tx, *user.ToUser())
	conflict := helper.NewCustomError(err, exception.ErrorConflict)
	if errCount != nil {
		tx.Rollback()
		errService = exception.NewError(errCount, exception.ErrorInternalServer)
		return
	}
	if user.Roles == model.Admin {
		if adminCount >= 1 {
			tx.Rollback()
			errService = exception.NewError(fmt.Errorf("admin is maximum one person"), exception.ErrorUnauthorized)
			return
		}
	}
	if badRequest != nil {
		tx.Rollback()
		errService = badRequest
		return
	} else if conflict != nil {
		tx.Rollback()
		errService = conflict
		return
	} else {
		tx.Commit()
		errService = nil
	}
	return
}

func (u *UsersService) CreateUsers(ctx context.Context, users model.UsersRequests) (errService error) {
	tx := u.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err := r.(error)
			errService = exception.NewError(err, exception.ErrorInternalServer)
		}
	}()

	usersStruct := model.UsersRequestsValidation{
		Users: users,
	}
	validationError := u.Validation.Struct(usersStruct)
	badRequest := helper.NewCustomError(validationError, exception.ErrorBadRequest)
	var adminCount int64
	errCount := tx.WithContext(ctx).Model(model.User{}).Where("roles = ?", model.Admin).Count(&adminCount).Error
	err := u.Repository.CreateUsers(ctx, tx, users.ToUsers())
	conflict := helper.NewCustomError(err, exception.ErrorConflict)
	if errCount != nil {
		tx.Rollback()
		errService = exception.NewError(errCount, exception.ErrorInternalServer)
		return
	}
	for _, user := range users {
		if user.Roles == model.Admin {
			if adminCount >= 1 {
				tx.Rollback()
				errService = exception.NewError(fmt.Errorf("admin is maximum one person"), exception.ErrorUnauthorized)
				return
			}
		}
	}
	if badRequest != nil {
		tx.Rollback()
		errService = badRequest
		return
	} else if conflict != nil {
		tx.Rollback()
		errService = conflict
		return
	} else {
		tx.Commit()
		errService = nil
	}
	return
}

func (u *UsersService) UpdateUserID(ctx context.Context, user model.UserLoginUpdateRequest, ID uuid.UUID) (errService error) {
	tx := u.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err := r.(error)
			errService = exception.NewError(err, exception.ErrorInternalServer)
		}
	}()

	validationError := u.Validation.Struct(user)
	badRequest := helper.NewCustomError(validationError, exception.ErrorBadRequest)

	exist := u.Repository.UsersExistByID(ctx, tx, ID)
	if user.Password != "" {
		user.Password, _ = bcrypts.HashPassword(user.Password, config.Other.SaltLevel)
	}
	u.Repository.UpdateUserID(ctx, tx, *user.ToUser(), ID)
	if badRequest != nil {
		tx.Rollback()
		errService = badRequest
		return
	} else if !exist {
		tx.Rollback()
		errService = exception.NewError(fmt.Errorf("Users With ID %v Not Found", ID), exception.ErrorNotFound)
		return
	} else {
		tx.Commit()
		errService = nil
	}
	return
}

func (u *UsersService) DeleteUserByID(ctx context.Context, ID uuid.UUID) (errService error) {
	tx := u.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err := r.(error)
			errService = exception.NewError(err, exception.ErrorInternalServer)
		}
	}()
	exist := u.Repository.UsersExistByID(ctx, tx, ID)

	u.Repository.DeleteUserByID(ctx, tx, ID)
	if !exist {
		tx.Rollback()
		errService = exception.NewError(fmt.Errorf("Users With ID %v Not Found", ID), exception.ErrorNotFound)
		return
	} else {
		tx.Commit()
		errService = nil
	}
	return
}

func (u *UsersService) DeleteUsersByIDs(ctx context.Context, IDs []uuid.UUID) (errService error) {
	tx := u.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err := r.(error)
			errService = exception.NewError(err, exception.ErrorInternalServer)
		}
	}()

	exist := u.Repository.UsersExistByIDs(ctx, tx, IDs)

	u.Repository.DeleteUsersByIDs(ctx, tx, IDs)
	if !exist {
		tx.Rollback()
		errService = exception.NewError(fmt.Errorf("Users With ID %v Not Found", IDs), exception.ErrorNotFound)
		return
	} else {
		tx.Commit()
		errService = nil
	}
	return
}

//type ctrl struct {
//	db
//	ctx
//}
