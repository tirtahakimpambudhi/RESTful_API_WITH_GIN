package helper

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"go_gin/internal/config"
	"go_gin/internal/domain/model"
	"go_gin/internal/exception"
	"strings"
)

func NewAccessToken(registeredClaims *model.StandardClaimsJWT) string {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, registeredClaims)
	accessTokenStr, err := accessToken.SignedString([]byte(config.JWT.SecretKey))
	Panic(err)
	return accessTokenStr
}

func NewRefreshToken(claims *jwt.StandardClaims) string {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshTokenStr, err := refreshToken.SignedString([]byte(config.JWT.SecretKey))
	Panic(err)
	return refreshTokenStr
}

func ParseAccessToken(accessToken string) (*model.StandardClaimsJWT, error) {
	parsedAccessToken, err := jwt.ParseWithClaims(accessToken, &model.StandardClaimsJWT{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWT.SecretKey), nil
	})

	if err != nil {
		return nil, err // Return the error if there's an issue with parsing
	}

	// Check if the claims type assertion is successful
	claims, ok := parsedAccessToken.Claims.(*model.StandardClaimsJWT)
	if !ok {
		return nil, errors.New("failed to parse claims") // Return an error if the type assertion fails
	}

	return claims, nil
}

func ParseRefreshToken(refreshToken string) (*jwt.StandardClaims, error) {
	parsedRefreshToken, err := jwt.ParseWithClaims(refreshToken, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWT.SecretKey), nil
	})

	if err != nil {
		return nil, err // Return the error if there's an issue with parsing
	}

	// Check if the claims type assertion is successful
	claims, ok := parsedRefreshToken.Claims.(*jwt.StandardClaims)
	if !ok {
		return nil, errors.New("failed to parse claims") // Return an error if the type assertion fails
	}

	return claims, nil
}

func VerifyAccessToken(accessToken string) (*model.StandardClaimsJWT, error) {
	parsedToken, err := ParseAccessToken(accessToken)
	if err != nil {
		if v, ok := err.(*jwt.ValidationError); ok {
			switch v.Errors {
			case jwt.ValidationErrorSignatureInvalid:
				return nil, exception.NewError(errors.New("signature is invalid"), exception.ErrorUnauthorized)
			case jwt.ValidationErrorExpired:
				return nil, exception.NewError(errors.New("token has expired"), exception.ErrorForbidden)
			default:
				return nil, exception.NewError(err, exception.ErrorUnauthorized)
			}
		}
		return nil, exception.NewError(err, exception.ErrorUnauthorized)
	}
	return parsedToken, nil
}

func VerifyRefreshToken(refreshToken string) (*jwt.StandardClaims, error) {
	parsedToken, err := ParseRefreshToken(refreshToken)
	if err != nil {
		if v, ok := err.(*jwt.ValidationError); ok {
			switch v.Errors {
			case jwt.ValidationErrorSignatureInvalid:
				return nil, exception.NewError(errors.New("signature is invalid"), exception.ErrorUnauthorized)
			case jwt.ValidationErrorExpired:
				return nil, exception.NewError(errors.New("token has expired"), exception.ErrorForbidden)
			default:
				return nil, exception.NewError(err, exception.ErrorUnauthorized)
			}
		}
		return nil, exception.NewError(err, exception.ErrorUnauthorized)
	}
	return parsedToken, nil
}

func ExtractBearerToken(header string) (string, error) {
	if header == "" {
		return "", exception.NewError(errors.New("UNAUTHORIZED"), exception.ErrorUnauthorized)
	}
	token := strings.Split(header, " ")
	if len(token) != 2 {
		return "", exception.NewError(errors.New("incorrectly formatted header authorization"), exception.ErrorBadRequest)
	}
	return token[1], nil
}
