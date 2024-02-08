package repository

import (
	"context"
	"github.com/google/uuid"
	"go_gin/internal/config"
	"go_gin/internal/domain/model"
	"go_gin/internal/domain/model/web"
	"go_gin/pkg/helper"
	"gorm.io/gorm"
	"strings"
)

type UsersRepository struct {
}

func (u *UsersRepository) GetUsersBySearch(ctx context.Context, DB *gorm.DB, query web.SearchValue) model.Users {
	var users model.Users
	valueSearch := []string{"%", query.Search, "%"}
	key := strings.Join(valueSearch, "")
	rows, err := DB.WithContext(ctx).Unscoped().Model(&model.User{}).Where("username LIKE ?", key).Where("email LIKE ?", key).Offset(int(query.Offset)).Limit(config.Other.Limit).Rows()
	helper.Panic(err)
	defer rows.Close()
	for rows.Next() {
		user := model.User{}
		if err := DB.ScanRows(rows, &user); err != nil {
			helper.Panic(err)
			return nil
		}
		users = append(users, user)
	}
	return users
}

func (u *UsersRepository) RestoreUserByID(ctx context.Context, DB *gorm.DB, ID uuid.UUID) {
	err := DB.WithContext(ctx).Unscoped().Model(&model.User{}).Where("id = ?", ID).Not("deleted_at", nil).Update("deleted_at", nil).Error
	helper.Panic(err)
}

func (u *UsersRepository) RestoreUsersByIDs(ctx context.Context, DB *gorm.DB, IDs []uuid.UUID) {
	err := DB.WithContext(ctx).Unscoped().Model(&model.User{}).Where("id IN ?", IDs).Not("deleted_at", nil).Update("deleted_at", nil).Error
	helper.Panic(err)
}

func NewUsersRepository() model.UsersRepository {
	return &UsersRepository{}
}

func (u *UsersRepository) GetUserByID(ctx context.Context, DB *gorm.DB, ID uuid.UUID) (model.User, error) {
	user := model.User{}
	err := DB.WithContext(ctx).Where("id = ?", ID).Take(&user).Error
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (u *UsersRepository) GetUserByEmail(ctx context.Context, DB *gorm.DB, Email string) (model.User, error) {
	user := model.User{}
	err := DB.WithContext(ctx).Where("email = ?", Email).Take(&user).Error
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (u *UsersRepository) GetUsers(ctx context.Context, DB *gorm.DB, query web.GetAllValue) model.Users {
	var users model.Users
	rows, err := DB.WithContext(ctx).Unscoped().Model(&model.User{}).Offset(int(query.Offset)).Limit(config.Other.Limit).Rows()
	helper.Panic(err)
	defer rows.Close()
	for rows.Next() {
		user := model.User{}
		if err := DB.ScanRows(rows, &user); err != nil {
			helper.Panic(err)
			return nil
		}
		users = append(users, user)
	}
	return users
}

func (u *UsersRepository) CreateUser(ctx context.Context, DB *gorm.DB, user model.User) error {
	err := DB.WithContext(ctx).Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *UsersRepository) CreateUsers(ctx context.Context, DB *gorm.DB, users model.Users) error {
	var err error
	if len(users) >= config.Other.LimitInsert {
		err = DB.CreateInBatches(&users, config.Other.BatchSize).Error
	} else {
		err = DB.Create(&users).Error
	}

	if err != nil {
		return err
	}
	return nil
}

func (u *UsersRepository) UpdateUserID(ctx context.Context, DB *gorm.DB, user model.User, ID uuid.UUID) {
	err := DB.WithContext(ctx).Where("id = ?", ID).Updates(&user).Error
	helper.Panic(err)
}

func (u *UsersRepository) DeleteUserByID(ctx context.Context, DB *gorm.DB, ID uuid.UUID) {
	err := DB.WithContext(ctx).Where("id = ?", ID).Delete(&model.User{}).Error
	helper.Panic(err)
}

func (u *UsersRepository) DeleteUsersByIDs(ctx context.Context, DB *gorm.DB, IDs []uuid.UUID) {
	err := DB.WithContext(ctx).Where("id IN ?", IDs).Delete(&model.User{}).Error
	helper.Panic(err)
}

func (u *UsersRepository) UsersExistByID(ctx context.Context, DB *gorm.DB, ID uuid.UUID) bool {
	var count int64
	err := DB.Model(&model.User{}).WithContext(ctx).Where("id = ?", ID).Count(&count).Error
	helper.Panic(err)
	return count == 1
}
func (u *UsersRepository) UsersExistByIDs(ctx context.Context, DB *gorm.DB, IDs []uuid.UUID) bool {
	var count int64
	err := DB.Model(&model.User{}).WithContext(ctx).Where("id IN ?", IDs).Count(&count).Error
	helper.Panic(err)
	return count == int64(len(IDs))
}
