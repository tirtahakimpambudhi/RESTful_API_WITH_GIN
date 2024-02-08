package model

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type StandardClaimsJWT struct {
	*jwt.StandardClaims
	ID       uuid.UUID `json:"id" gorm:"column:id" validate:"required"`
	Username string    `json:"username" gorm:"column:username" validate:"required,min=5,max=100"`
	Email    string    `json:"email" gorm:"column:email;unique" validate:"required,email,max=100"`
}

func NewStandardClaimsJWT(registeredClaims *jwt.StandardClaims, ID uuid.UUID, username string, email string) *StandardClaimsJWT {
	return &StandardClaimsJWT{StandardClaims: registeredClaims, ID: ID, Username: username, Email: email}
}
