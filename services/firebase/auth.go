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
func (fb *AuthService) Create(userRequest AuthUser, setClaims ...func(claims types.MapString) types.MapString) (string, error) {
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
			errMessage := fmt.Errorf("%s with this email already exits in Firebase", userRequest.Role)
			return "", api_errors.BadRequest.Wrap(errMessage, errMsg)
		}
		return "", api_errors.InternalError.Wrap(err, errMsg)
	}

	claims := types.MapString{constants.Roles.Key: userRequest.Role}

	for _, setClaim := range setClaims {
		claims = setClaim(claims)
	}

	err = fb.SetClaim(u.UID, claims)
	if err != nil {
		return "Internal Server Error", err
	}
	return u.UID, err
}

// CreateUser creates a new user with email and password
func (fb *AuthService) CreateUser(userRequest AuthUser) (string, error) {
	return fb.Create(userRequest, func(claims types.MapString) types.MapString {
		claims[constants.Claims.UserId.Name()] = userRequest.UserID
		return claims
	})
}

// CreateAdmin creates a new admin with email and password
func (fb *AuthService) CreateAdmin(userRequest AuthUser) (string, error) {
	return fb.Create(userRequest, func(claims types.MapString) types.MapString {
		if userRequest.Role != constants.Roles.Admin.ToString() {
			claims[constants.Claims.AdminId.ToString()] = userRequest.AdminID
		}
		return claims
	})
}

// VerifyToken verify passed firebase id token
func (fb *AuthService) VerifyToken(idToken string) (*auth.Token, error) {
	token, err := fb.VerifyIDToken(context.Background(), idToken)
	return token, err
}

// SetClaim set's claim to firebase user
func (fb *AuthService) SetClaim(uid string, claims map[string]interface{}) error {
	err := fb.SetCustomUserClaims(context.Background(), uid, claims)
	return err
}

// UpdateEmailVerification update firebase user email verify
func (fb *AuthService) UpdateEmailVerification(uid string) error {
	params := (&auth.UserToUpdate{}).
		EmailVerified(true)
	_, err := fb.UpdateUser(context.Background(), uid, params)
	return err
}

// DisableUser true for disabled; false for enabled.
func (fb *AuthService) DisableUser(uid string, disable bool) error {
	params := (&auth.UserToUpdate{}).
		Disabled(disable)
	_, err := fb.UpdateUser(context.Background(), uid, params)
	return err
}

// UpdateFirebaseAdmin handles the common operation to update admin in Firebase for OneStore Admin and Admin
func (fb *AuthService) UpdateFirebaseAdmin(UID string, newUserData, oldUserData AuthUser) error {
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
			return err
		}
	}
	return nil
}

func (fb *AuthService) CreateCustomToken(ctx context.Context, uid string) (string, error) {
	token, err := fb.CustomToken(ctx, uid)
	if err != nil {
		return "", err
	}
	return token, nil
}
