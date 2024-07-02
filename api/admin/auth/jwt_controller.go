package auth

import (
	"fmt"
	"net/http"
	"time"

	"boilerplate-api/api/admin/user"
	"boilerplate-api/internal/api_errors"
	"boilerplate-api/internal/auth"
	"boilerplate-api/internal/config"
	"boilerplate-api/internal/constants"
	"boilerplate-api/internal/json_response"
	"boilerplate-api/internal/request_validator"
	"boilerplate-api/internal/types"
	"boilerplate-api/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// FIXME :: refactor

// JwtAuthController struct
type JwtAuthController struct {
	logger      config.Logger
	userService user.Service
	jwtService  auth.JWTAuthService
	env         config.Env
	validator   request_validator.Validator
}

// NewJwtAuthController constructor
func NewJwtAuthController(
	logger config.Logger,
	userService user.Service,
	jwtService auth.JWTAuthService,
	env config.Env,
	validator request_validator.Validator,
) JwtAuthController {
	return JwtAuthController{
		logger:      logger,
		userService: userService,
		jwtService:  jwtService,
		env:         env,
		validator:   validator,
	}
}

func (cc JwtAuthController) LoginUserWithJWT(c *gin.Context) {
	reqData := JWTLoginRequestData{}
	// Bind the request payload to a reqData struct
	if err := c.ShouldBindJSON(&reqData); err != nil {
		cc.logger.Error("Error [ShouldBindJSON] : ", err.Error())
		err := api_errors.BadRequest.Wrap(err, "Failed to bind request data")
		status, errM := api_errors.HandleError(err)
		c.JSON(status, json_response.Error{Error: errM})
		return
	}

	// validating using custom validator
	if validationErr := cc.validator.Struct(reqData); validationErr != nil {
		cc.logger.Error("[Validate Struct] Validation error: ", validationErr.Error())
		err := api_errors.BadRequest.Wrap(validationErr, "Validation error")
		err = api_errors.SetCustomMessage(err, "Invalid input information")
		err = api_errors.AddErrorContextBlock(err, cc.validator.GenerateValidationResponse(validationErr))
		status, errM := api_errors.HandleError(err)
		c.JSON(status, json_response.Error{Error: errM})
		return
	}

	// Check if the user exists with provided email address
	user, err := cc.userService.GetOneUserWithEmail(reqData.Email)
	if err != nil {
		err := api_errors.BadRequest.New("Invalid user credentials")
		status, errM := api_errors.HandleError(err)
		c.JSON(status, json_response.Error{Error: errM})
		return
	}

	// Check if the password is correct
	// Thus password is encrypted and saved in DB, comparing plain text with it's hash
	isValidPassword := utils.CompareHashAndPlainPassword(user.Password, reqData.Password)
	if !isValidPassword {
		cc.logger.Error("[CompareHashAndPassword] hash and plain password doesnot match")
		status, errM := api_errors.HandleError(
			api_errors.BadRequest.New("Invalid user credentials"),
		)
		c.JSON(status, json_response.Error{Error: errM})
		return
	}

	// Create a new JWT access claims object
	accessClaims := auth.JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(cc.env.JwtAccessTokenExpiresAt))),
			ID:        fmt.Sprintf("%v", user.ID),
		},
		//Add other claims
	}

	// Create a new JWT Access token using the claims and the secret key
	accessToken, tokenErr := cc.jwtService.GenerateToken(accessClaims, cc.env.JwtAccessSecret)
	if tokenErr != nil {
		cc.logger.Error("[SignedString] Error getting token: ", tokenErr.Error())
		status, errM := api_errors.HandleError(
			api_errors.InternalError.New(tokenErr.Error()),
		)
		c.JSON(status, json_response.Error{Error: errM})
		return
	}

	// Create a new JWT refresh claims object
	refreshClaims := auth.JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(cc.env.JwtRefreshTokenExpiresAt))),
			ID:        fmt.Sprintf("%v", user.ID),
		},
	}

	// Create a new JWT Refresh token using the claims and the secret key
	refreshToken, refreshTokenErr := cc.jwtService.GenerateToken(refreshClaims, cc.env.JwtRefreshSecret)
	if refreshTokenErr != nil {
		cc.logger.Error("[SignedString] Error getting token: ", refreshTokenErr.Error())
		status, errM := api_errors.HandleError(
			api_errors.InternalError.New(refreshTokenErr.Error()),
		)
		c.JSON(status, json_response.Error{Error: errM})
		return
	}

	data := types.MapString{
		"user":          user.ToMap(),
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}

	c.JSON(http.StatusOK, json_response.Data[types.MapString]{Data: data})
}

func (cc JwtAuthController) RefreshJwtToken(c *gin.Context) {
	// Get the token from the request header
	header := c.GetHeader(constants.Headers.Authorization.ToString())

	tokenString, err := cc.jwtService.GetTokenFromHeader(header)
	if err != nil {
		cc.logger.Error("Error getting token from header: ", err.Error())
		err = api_errors.Unauthorized.Wrap(err, "Something went wrong")
		status, errM := api_errors.HandleError(err)
		c.JSON(status, json_response.Error{Error: errM})
		return
	}

	parsedToken, parseErr := cc.jwtService.ParseAndVerifyToken(tokenString, cc.env.JwtRefreshSecret)
	if parseErr != nil {
		cc.logger.Error("Error parsing token: ", parseErr.Error())
		err = api_errors.Unauthorized.Wrap(parseErr, "Something went wrong")
		status, errM := api_errors.HandleError(err)
		c.JSON(status, json_response.Error{Error: errM})
		return
	}

	claims, verifyErr := cc.jwtService.RetrieveClaims(parsedToken)
	if verifyErr != nil {
		cc.logger.Error("Error veriefying token: ", verifyErr.Error())
		err = api_errors.Unauthorized.Wrap(verifyErr, "Something went wrong")
		status, errM := api_errors.HandleError(err)
		c.JSON(status, json_response.Error{Error: errM})
		return
	}

	// Create a new JWT Access claims
	accessClaims := auth.JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(cc.env.JwtAccessTokenExpiresAt))),
			ID:        fmt.Sprintf("%v", claims.ID),
		},
		// Add other claims
	}

	// Create a new JWT token using the claims and the secret key
	accessToken, tokenErr := cc.jwtService.GenerateToken(accessClaims, cc.env.JwtAccessSecret)
	if tokenErr != nil {
		cc.logger.Error("[SignedString] Error getting token: ", tokenErr.Error())
		status, errM := api_errors.HandleError(
			api_errors.InternalError.New(tokenErr.Error()),
		)
		c.JSON(status, json_response.Error{Error: errM})
		return
	}

	data := types.MapString{
		"access_token": accessToken,
		"expires_at":   accessClaims.ExpiresAt,
	}

	c.JSON(http.StatusOK, json_response.Data[types.MapString]{Data: data})
}
