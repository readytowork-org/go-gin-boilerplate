package user

import (
	"net/http"

	"boilerplate-api/internal/api_errors"
	"boilerplate-api/internal/config"
	"boilerplate-api/internal/constants"
	"boilerplate-api/internal/json_response"
	"boilerplate-api/internal/request_validator"
	"boilerplate-api/internal/utils"
	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
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

// @Tags			UserManagementApi
// @Summary		Create User
// @Description	Create one user
// @Security		Bearer
// @Produce		application/json
// @Param			data	body		CreateUserRequestData	true	"Enter JSON"
// @Success		200		{object}	json_response.Message	"OK"
// @Failure		400		{object}	json_response.Error
// @Failure		500		{object}	json_response.Error
// @Router			/api/v1/users [post]
// @Id				CreateUser
func (cc Controller) CreateUser(c *gin.Context) {
	reqData := CreateUserRequestData{}
	trx := c.MustGet(constants.DBTransaction).(*gorm.DB)

	if err := c.ShouldBindJSON(&reqData); err != nil {
		cc.logger.Error("Error [CreateUser] (ShouldBindJson) : ", err)
		err := api_errors.BadRequest.Wrap(err, "Failed to bind user data")
		status, errM := api_errors.HandleError(err)
		c.JSON(status, json_response.Error{Error: errM})
		return
	}
	if validationErr := cc.validator.Struct(reqData); validationErr != nil {
		err := api_errors.BadRequest.Wrap(validationErr, "Validation error")
		err = api_errors.SetCustomMessage(err, "Invalid input information")
		err = api_errors.AddErrorContextBlock(err, cc.validator.GenerateValidationResponse(validationErr))
		status, errM := api_errors.HandleError(err)
		c.JSON(status, json_response.Error{Error: errM})
		return
	}

	if reqData.Password != reqData.ConfirmPassword {
		cc.logger.Error("Password and confirm password not matching : ")
		err := api_errors.BadRequest.New("Password and confirm password should be same.")
		status, errM := api_errors.HandleError(err)
		c.JSON(status, json_response.Error{Error: errM})
		return
	}

	if _, err := cc.userService.GetOneUserWithEmail(reqData.Email); err != nil {
		cc.logger.Error("Error [CreateUser] [db CreateUser]: User with this email already exists")
		err := api_errors.BadRequest.New("User with this email already exists")
		status, errM := api_errors.HandleError(err)
		c.JSON(status, json_response.Error{Error: errM})
		return
	}

	if _, err := cc.userService.GetOneUserWithPhone(reqData.Phone); err != nil {
		cc.logger.Error("Error [db GetOneUserWithPhone]: User with this phone already exists")
		err := api_errors.BadRequest.New("User with this phone already exists")
		status, errM := api_errors.HandleError(err)
		c.JSON(status, json_response.Error{Error: errM})
		return
	}

	if err := cc.userService.WithTrx(trx).CreateUser(reqData.User); err != nil {
		cc.logger.Error("Error [CreateUser] [db CreateUser]: ", err.Error())
		err := api_errors.InternalError.Wrap(err, "Failed to create user")
		status, errM := api_errors.HandleError(err)
		c.JSON(status, json_response.Error{Error: errM})
		return
	}

	c.JSON(http.StatusOK, json_response.Message{
		Msg: "User Created Successfully",
	})
}

// @Tags			UserManagementApi
// @Summary		All users
// @Description	get all users
// @Security		Bearer
// @Produce		application/json
// @Param			pagination	query		Pagination	false	"query param"
// @Success		200			{object}	json_response.DataCount[GetUserResponse]
// @Failure		500			{object}	json_response.Error
// @Router			/api/v1/users [get]
// @Id				GetAllUsers
func (cc Controller) GetAllUsers(c *gin.Context) {
	pagination := utils.BuildPagination[*Pagination](c)

	users, count, err := cc.userService.GetAllUsers(*pagination)
	if err != nil {
		cc.logger.Error("Error finding user records", err.Error())
		err := api_errors.InternalError.Wrap(err, "Failed to get users data")
		status, errM := api_errors.HandleError(err)
		c.JSON(status, json_response.Error{Error: errM})
		return
	}

	c.JSON(http.StatusOK, json_response.DataCount[GetUserResponse]{
		Count: count,
		Data:  users,
	})
}

// @Tags			UserManagementApi
// @Summary		User Profile
// @Description	get user profile
// @Security		Bearer
// @Produce		application/json
// @Success		200	{object}	json_response.Data[GetUserResponse]
// @Failure		500	{object}	json_response.Error
// @Router			/api/v1/{id} [get]
// @Id				GetOneUser
func (cc Controller) GetOneUser(c *gin.Context) {
	userID, err := utils.StringToInt64(c.Param("id"))

	user, err := cc.userService.GetOneUser(userID)
	if err != nil {
		cc.logger.Error("Error finding user", err.Error())
		err := api_errors.InternalError.Wrap(err, "Failed to get user")
		status, errM := api_errors.HandleError(err)
		c.JSON(status, json_response.Error{Error: errM})
		return
	}

	c.JSON(http.StatusOK, json_response.Data[GetUserResponse]{
		Data: user,
	})
}
