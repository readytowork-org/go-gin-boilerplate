package firebase

import (
	"context"
	"fmt"
	"strings"

	"boilerplate-api/internal/api_errors"
	"boilerplate-api/internal/config"
	"boilerplate-api/internal/constants"
	"boilerplate-api/internal/types"
	"firebase.google.com/go"

	"firebase.google.com/go/auth"
)

type AuthUser struct {
	Password    string
	Role        string
	DisplayName *string
	Email       string
	AdminID     int64
	UserID      int64
}

// AuthService structure
type AuthService struct {
	*auth.Client
}

// NewFirebaseAuthService creates new firebase service
func NewFirebaseAuthService(logger config.Logger, app *firebase.App) AuthService {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	firebaseAuth, err := app.Auth(ctx)
	if err != nil {
		logger.Fatalf("Firebase Authentication: %v", err)
	}

	return AuthService{
		Client: firebaseAuth,
	}
}

// Create creates a new user with email and password
func (fb *AuthService) Create(
	userRequest AuthUser, setClaims ...func(claims types.MapString) types.MapString) (
	string, *api_errors.ErrorResponse) {

	params := (&auth.UserToCreate{}).
		Email(userRequest.Email).
		Password(userRequest.Password)

	if userRequest.DisplayName != nil && *userRequest.DisplayName != "" {
		params = params.DisplayName(*userRequest.DisplayName)
	}

	u, err := fb.Client.CreateUser(context.Background(), params)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to create %s", userRequest.Role)
		if strings.Contains(err.Error(), "EMAIL_EXISTS") {
			errMessage := fmt.Sprintf("%s with this email already exits in Firebase", userRequest.Role)
			return "", &api_errors.ErrorResponse{
				ErrorType: api_errors.BadRequest,
				Message:   errMessage,
			}
		}
		return "", &api_errors.ErrorResponse{
			ErrorType: api_errors.InternalError,
			Message:   errMsg,
		}
	}

	claims := types.MapString{constants.Roles.Key: userRequest.Role}

	for _, setClaim := range setClaims {
		claims = setClaim(claims)
	}

	errResponse := fb.SetClaim(u.UID, claims)
	if errResponse != nil {
		return "", &api_errors.ErrorResponse{
			Message:   errResponse.Message,
			ErrorType: api_errors.InternalError,
		}
	}
	return u.UID, nil
}

// CreateUser creates a new user with email and password
func (fb *AuthService) CreateUser(userRequest AuthUser) (string, *api_errors.ErrorResponse) {
	return fb.Create(userRequest, func(claims types.MapString) types.MapString {
		claims[constants.Claims.UserId.Name()] = userRequest.UserID
		return claims
	})
}

// CreateAdmin creates a new admin with email and password
func (fb *AuthService) CreateAdmin(userRequest AuthUser) (string, *api_errors.ErrorResponse) {
	return fb.Create(userRequest, func(claims types.MapString) types.MapString {
		if userRequest.Role != constants.Roles.Admin.ToString() {
			claims[constants.Claims.AdminId.ToString()] = userRequest.AdminID
		}
		return claims
	})
}

// VerifyToken verify passed firebase id token
func (fb *AuthService) VerifyToken(idToken string) (*auth.Token, *api_errors.ErrorResponse) {
	token, err := fb.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		return nil, &api_errors.ErrorResponse{
			Message: err.Error(),
		}
	}
	return token, nil
}

// SetClaim set's claim to firebase user
func (fb *AuthService) SetClaim(uid string, claims map[string]interface{}) *api_errors.ErrorResponse {
	err := fb.SetCustomUserClaims(context.Background(), uid, claims)
	return &api_errors.ErrorResponse{
		Message: err.Error(),
	}
}

// UpdateEmailVerification update firebase user email verify
func (fb *AuthService) UpdateEmailVerification(uid string) *api_errors.ErrorResponse {
	params := (&auth.UserToUpdate{}).
		EmailVerified(true)
	_, err := fb.UpdateUser(context.Background(), uid, params)
	return &api_errors.ErrorResponse{
		Message: err.Error(),
	}
}

// DisableUser true for disabled; false for enabled.
func (fb *AuthService) DisableUser(uid string, disable bool) *api_errors.ErrorResponse {
	params := (&auth.UserToUpdate{}).
		Disabled(disable)
	_, err := fb.UpdateUser(context.Background(), uid, params)
	return &api_errors.ErrorResponse{
		Message: err.Error(),
	}
}

// UpdateFirebaseAdmin handles the common operation to update admin in Firebase for OneStore Admin and Admin
func (fb *AuthService) UpdateFirebaseAdmin(UID string, newUserData, oldUserData AuthUser) *api_errors.ErrorResponse {
	fbAdmin := &auth.UserToUpdate{}

	if newUserData.Email != "" && newUserData.Email != oldUserData.Email {
		fbAdmin = fbAdmin.Email(newUserData.Email)
	}

	if newUserData.Password != "" {
		fbAdmin = fbAdmin.Password(newUserData.Password)
	}

	if newUserData.DisplayName != nil && newUserData.DisplayName != oldUserData.DisplayName {
		fbAdmin = fbAdmin.DisplayName(*newUserData.DisplayName)
	}

	if fbAdmin != nil {
		if _, err := fb.UpdateUser(context.Background(), UID, fbAdmin); err != nil {
			return &api_errors.ErrorResponse{
				Message: err.Error(),
			}
		}
	}
	return nil
}

func (fb *AuthService) CreateCustomToken(ctx context.Context, uid string) (string, *api_errors.ErrorResponse) {
	token, err := fb.CustomToken(ctx, uid)
	if err != nil {
		return "", &api_errors.ErrorResponse{
			Message: err.Error(),
		}
	}
	return token, nil
}
