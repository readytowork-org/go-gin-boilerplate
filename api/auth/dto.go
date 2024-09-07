package auth

// JWTLoginRequestData Request body data to authenticate user with jwt-auth
type JWTLoginRequestData struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}
