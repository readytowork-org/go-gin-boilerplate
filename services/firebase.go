package services

import "context"

// AuthErrorResponse structure
type AuthErrorResponse struct {
	Message   string `json:"message"`
	ErrorType int    `json:"error_type"`
}

// FirebaseToken Replace this with firebase.google.com/go/auth auth.Token
type FirebaseToken struct {
	AuthTime int64                  `json:"auth_time"`
	Issuer   string                 `json:"iss"`
	Audience string                 `json:"aud"`
	Expires  int64                  `json:"exp"`
	IssuedAt int64                  `json:"iat"`
	Subject  string                 `json:"sub,omitempty"`
	UID      string                 `json:"uid,omitempty"`
	Claims   map[string]interface{} `json:"-"`
}

type IFirebaseMiddlewareService interface {
	VerifyToken(token string) (*FirebaseToken, *AuthErrorResponse)
}

type IFirebaseAuthService interface {
	GetUserByEmail(context context.Context, email string) (interface{}, *error)
	CreateUser(displayName, email, password, role string) (interface{}, *AuthErrorResponse)
}
