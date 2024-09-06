package auth

import (
	"strings"

	"boilerplate-api/internal/api_errors"
	"boilerplate-api/internal/config"
	"boilerplate-api/internal/constants"

	"github.com/golang-jwt/jwt/v4"
)

type JWTClaims struct {
	jwt.RegisteredClaims
	// ...other claims
}

type JWTAuthService struct {
	logger config.Logger
	env    config.Env
}

func NewJWTAuthService(
	logger config.Logger,
	env config.Env,
) JWTAuthService {
	return JWTAuthService{
		logger: logger,
		env:    env,
	}
}

func (m JWTAuthService) GetTokenFromHeader(header string) (string, *api_errors.ErrorResponse) {
	if header == "" {
		err := api_errors.ErrorResponse{
			Message: "Authorization token is required in header",
		}
		m.logger.Error("[GetHeader]: ", err.Message)
		return "", &err
	}

	if !strings.Contains(header, constants.TokenTypes.Bearer.ToString()) {
		err := api_errors.ErrorResponse{
			Message: "Token type is required",
		}
		m.logger.Error("Missing token type: ", err.Message)
		return "", &err
	}
	tokenString := strings.TrimSpace(strings.Replace(header, "Bearer", "", 1))
	return tokenString, nil

}

func (m JWTAuthService) ParseAndVerifyToken(tokenString, secret string) (*jwt.Token, *api_errors.ErrorResponse) {
	// Parse the token using the secret key
	token, err := jwt.ParseWithClaims(
		tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
	)
	if err != nil {
		if !strings.Contains(err.Error(), "expired") {
			err := api_errors.ErrorResponse{
				Message: "Invalid ID token",
			}
			return nil, &err
		}
		m.logger.Error("Invalid token[ParseWithClaims] :", err.Error())
		return nil, &api_errors.ErrorResponse{
			Message: err.Error(),
		}
	}
	return token, nil
}

func (m JWTAuthService) RetrieveClaims(token *jwt.Token) (*JWTClaims, *api_errors.ErrorResponse) {
	// Verify token
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		err := api_errors.ErrorResponse{
			Message: "Invalid ID token",
		}
		m.logger.Error("Invalid token [token.Valid]: ", err.Message)
		return nil, &err
	}
	return claims, nil
}

func (m JWTAuthService) GenerateToken(claims JWTClaims, secret string) (string, error) {
	// Create a new JWT token using the claims and the secret key
	tokenClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, tokenErr := tokenClaim.SignedString([]byte(secret))
	if tokenErr != nil {
		return "", tokenErr
	}
	return token, nil
}
