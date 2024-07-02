package user

import (
	"fmt"
	"net/http"

	"boilerplate-api/database/models"
	"boilerplate-api/internal/api_errors"
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
// @Success		200	{object}	json_response.Data[GetUserResponse]
// @Failure		500	{object}	json_response.Error
// @Router			/api/v1/profile [get]
// @Id				GetUserProfile
func (cc Controller) GetUserProfile(c *gin.Context) {
	userID := fmt.Sprintf("%v", c.MustGet(constants.UserID))

	user, err := cc.userService.GetOneUser(userID)
	if err != nil {
		cc.logger.Error("Error finding user profile", err.Error())
		err := api_errors.InternalError.Wrap(err, "Failed to get users profile data")
		status, errM := api_errors.HandleError(err)
		c.JSON(status, json_response.Error{Error: errM})
		return
	}

	c.JSON(http.StatusOK, json_response.Data[models.User]{
		Data: user,
	})
}
