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
// @Success		200		{object}	json_response.Message	"User Created Successfully"
// @Failure		400		{object}	json_response.Error[string]
// @Failure		422		{object}	json_response.Error[[]api_errors.ValidationError]
// @Failure		500		{object}	json_response.Error[string]
// @Router			/api/v1/users [post]
// @Id				CreateUser
func (cc Controller) CreateUser(c *gin.Context) {
	reqData := CreateUserRequestData{}
	trx := c.MustGet(constants.DBTransaction).(*gorm.DB)

	if err := c.ShouldBindJSON(&reqData); err != nil {
		cc.logger.Error("Error [CreateUser] (ShouldBindJson) : ", err)
		c.JSON(http.StatusBadRequest, json_response.Error[string]{
			Error:   err.Error(),
			Message: "Failed to bind user data",
		})
		return
	}
	if validationErr := cc.validator.Struct(reqData); validationErr != nil {
		c.JSON(http.StatusUnprocessableEntity, json_response.Error[[]api_errors.ValidationError]{
			Error:   cc.validator.GenerateValidationResponse(validationErr),
			Message: "Invalid input information",
		})
		return
	}

	if reqData.Password != reqData.ConfirmPassword {
		cc.logger.Error("Password and confirm password not matching : ")
		c.JSON(http.StatusBadRequest, json_response.Error[string]{
			Error:   "Failed to create User",
			Message: "Password and confirm password should be same.",
		})
		return
	}

	if _, err := cc.userService.GetOneUserWithEmail(reqData.Email); err != nil {
		cc.logger.Error("Error [CreateUser] [db CreateUser]: User with this email already exists")
		c.JSON(http.StatusBadRequest, json_response.Error[string]{
			Error:   "Failed to create User",
			Message: "User with this email already exists",
		})
		return
	}

	if _, err := cc.userService.GetOneUserWithPhone(reqData.Phone); err != nil {
		cc.logger.Error("Error [db GetOneUserWithPhone]: User with this phone already exists")
		c.JSON(http.StatusBadRequest, json_response.Error[string]{
			Error:   "Failed to create User",
			Message: "User with this phone already exists",
		})
		return
	}

	if err := cc.userService.WithTrx(trx).CreateUser(reqData.User); err != nil {
		cc.logger.Error("Error [CreateUser] [db CreateUser]: ", err.Error())
		c.JSON(http.StatusInternalServerError, json_response.Error[string]{
			Error:   err.Error(),
			Message: "Failed to create User",
		})
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
// @Failure		500			{object}	json_response.Error[string]
// @Router			/api/v1/users [get]
// @Id				GetAllUsers
func (cc Controller) GetAllUsers(c *gin.Context) {
	pagination := utils.BuildPagination[*Pagination](c)

	users, count, err := cc.userService.GetAllUsers(*pagination)
	if err != nil {
		cc.logger.Error("Error finding user records", err.Error())
		c.JSON(http.StatusInternalServerError, json_response.Error[string]{
			Error:   err.Error(),
			Message: "Failed to get users data",
		})
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
// @Failure		500	{object}	json_response.Error[string]
// @Router			/api/v1/{id} [get]
// @Id				GetOneUser
func (cc Controller) GetOneUser(c *gin.Context) {
	userID, errResponse := utils.StringToInt64(c.Param("id"))
	if errResponse != nil {
		cc.logger.Error("Error finding user", errResponse.Message)
		c.JSON(http.StatusInternalServerError, json_response.Error[string]{
			Error:   errResponse.Message,
			Message: "Failed to get user",
		})
		return
	}

	user, err := cc.userService.GetOneUser(userID)
	if err != nil {
		cc.logger.Error("Error finding user", err.Error())
		c.JSON(http.StatusInternalServerError, json_response.Error[string]{
			Error:   err.Error(),
			Message: "Failed to get user",
		})
		return
	}

	c.JSON(http.StatusOK, json_response.Data[GetUserResponse]{
		Data: user,
	})
}
