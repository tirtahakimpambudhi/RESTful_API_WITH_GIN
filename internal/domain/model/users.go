package model

import (
	"database/sql"
	"database/sql/driver"
	"github.com/google/uuid"
	"go_gin/pkg/bcrypts"
	"gorm.io/gorm"
	"time"
)

type UserRole string

func (ct *UserRole) Scan(value interface{}) error {
	if value == nil {
		*ct = ""
		return nil
	}
	switch v := value.(type) {
	case []byte:
		*ct = UserRole(v)
	case string:
		*ct = UserRole(v)
	}
	return nil
}

func (ct UserRole) Value() (driver.Value, error) {
	return string(ct), nil
}

const (
	Admin     UserRole = "ADMIN"
	Moderator UserRole = "MODERATOR"
	Basic     UserRole = "BASIC"
)

type User struct {
	ID           uuid.UUID      `json:"id" gorm:"primaryKey;column:id"`
	Username     string         `json:"username" gorm:"column:username"`
	Email        string         `json:"email" gorm:"column:email;unique"`
	Password     string         `json:"password" gorm:"column:password"`
	RefreshToken sql.NullString `json:"refreshToken" gorm:"column:refresh_token"`
	Roles        UserRole       `gorm:"type:user_role;default:BASIC" json:"roles"`
	CreatedAt    time.Time      `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time      `json:"updatedAt" gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"deletedAt" gorm:"column:deleted_at"`
	TodoLists    TodoLists      `json:"todo_lists" gorm:"foreignKey:user_id;references:id"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	hashPass, err := bcrypts.HashPassword(u.Password, 10)
	if err != nil {
		return err
	}
	u.Password = hashPass
	return nil
}

type UserResponse struct {
	ID        uuid.UUID      `json:"id" gorm:"primaryKey;column:id"`
	Username  string         `json:"username" gorm:"column:username"`
	Email     string         `json:"email" gorm:"column:email;unique"`
	CreatedAt time.Time      `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"column:deleted_at"`
}

type UserRequest struct {
	ID       uuid.UUID `json:"id" gorm:"column:id" validate:"required"`
	Username string    `json:"username" gorm:"column:username" validate:"required,min=5,max=100"`
	Email    string    `json:"email" gorm:"column:email;unique" validate:"required,email,max=100"`
	Password string    `json:"password" gorm:"column:password" validate:"required,min=8"`
	Roles    UserRole  `json:"role" gorm:"type:enum('ADMIN', 'BASIC', 'MODERATOR');column:roles" validate:"required,eq=ADMIN|eq=MODERATOR|eq=BASIC"`
}

type UserLoginUpdateRequest struct {
	Username string `json:"username" gorm:"column:username" validate:"required,min=5,max=100"`
	Email    string `json:"email" gorm:"column:email;unique" validate:"required,email,max=100"`
	Password string `json:"password" gorm:"column:password" validate:"required,min=8"`
}

func (u *User) ToUserResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: u.DeletedAt,
	}
}

func (u *UserRequest) ToUser() *User {
	return &User{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		Password: u.Password,
		Roles:    u.Roles,
	}
}

func (u *UserLoginUpdateRequest) ToUser() *User {
	return &User{
		Username: u.Username,
		Email:    u.Email,
		Password: u.Password,
	}
}
