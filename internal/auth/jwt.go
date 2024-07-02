package auth

import (
	"strings"

	"boilerplate-api/internal/api_errors"
	"boilerplate-api/internal/config"
	"boilerplate-api/internal/constants"

	"github.com/golang-jwt/jwt/v4"
)

// FIXME :: refactor

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

func (m JWTAuthService) GetTokenFromHeader(header string) (string, error) {
	if header == "" {
		err := api_errors.BadRequest.New("Authorization token is required in header")
		err = api_errors.SetCustomMessage(err, "Authorization token is required in header")
		m.logger.Error("[GetHeader]: ", err.Error())
		return "", err
	}

	if !strings.Contains(header, constants.TokenTypes.Bearer.ToString()) {
		err := api_errors.BadRequest.New("Token type is required")
		m.logger.Error("Missing token type: ", err.Error())
		return "", err
	}
	tokenString := strings.TrimSpace(strings.Replace(header, "Bearer", "", 1))
	return tokenString, nil

}

func (m JWTAuthService) ParseAndVerifyToken(tokenString, secret string) (*jwt.Token, error) {
	// Parse the token using the secret key
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		if !strings.Contains(err.Error(), "expired") {
			m.logger.Error("Invalid token[ParseWithClaims] :", err.Error())
			err := api_errors.BadRequest.New("Invalid ID token")
			return nil, err
		}
		m.logger.Error("Invalid token[ParseWithClaims] :", err.Error())
		return nil, err
	}
	return token, nil

}

func (m JWTAuthService) RetrieveClaims(token *jwt.Token) (*JWTClaims, error) {
	// Verify token
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		err := api_errors.BadRequest.New("Invalid token")
		err = api_errors.SetCustomMessage(err, "Invalid token")
		m.logger.Error("Invalid token [token.Valid]: ", err.Error())
		return nil, err
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
