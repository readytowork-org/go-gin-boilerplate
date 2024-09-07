package user

import (
	"fmt"
	"net/http"

	"boilerplate-api/internal/config"
	"boilerplate-api/internal/constants"
	"boilerplate-api/internal/json_response"
	"boilerplate-api/internal/request_validator"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	logger      config.Logger
	userService Service
	env         config.Env
	validator   request_validator.Validator
}

// NewController Creates New user controller
func NewController(
	logger config.Logger,
	userService Service,
	env config.Env,
	validator request_validator.Validator,
) Controller {
	return Controller{
		logger:      logger,
		userService: userService,
		env:         env,
		validator:   validator,
	}
}

// @Tags			UserApi
// @Summary		User Profile
// @Description	get user profile
// @Security		Bearer
// @Produce		application/json
// @Success		200	{object}	json_response.Data[CUser]
// @Failure		500	{object}	json_response.Error[string]
// @Router			/api/v1/profile [get]
// @Id				GetUserProfile
func (cc Controller) GetUserProfile(c *gin.Context) {
	userID := fmt.Sprintf("%v", c.MustGet(constants.UserID))

	user, err := cc.userService.GetOneUser(userID)
	if err != nil {
		cc.logger.Error("Error finding user profile", err.Error())
		c.JSON(
			http.StatusInternalServerError, json_response.Error[string]{
				Error:   err.Error(),
				Message: "Failed to get users profile data",
			},
		)
		return
	}

	c.JSON(
		http.StatusOK, json_response.Data[CUser]{
			Data: user,
		},
	)
}
