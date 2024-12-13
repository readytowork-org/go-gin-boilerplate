package constants

const (
	// DBTransaction is database transaction handle set at router context
	DBTransaction = "db_trx"
)

const (
	// UID authenticated user's id
	UID    = "UID"
	UserID = "user_id_db"
)

const (
	RateLimit = "RateLimit"
)

const (
	BucketName = "bucket_name"
)

// Claim custom token claims
type Claim string

var Claims = struct {
	UID      Claim
	UserIdDb Claim
	UserId   Claim
	AdminId  Claim
}{
	UID:      "UID",
	UserIdDb: "user-id-db",
	UserId:   "user-id",
	AdminId:  "admin-id",
}

func (r Claim) Name() string {
	return "claim"
}

func (r Claim) ToString() string {
	return string(r)
}

// Role roles of users
type Role string

var Roles = struct {
	SuperAdmin Role
	Admin      Role
	User       Role
	Key        string
}{
	SuperAdmin: "super-amin",
	Admin:      "admin",
	User:       "user",
	Key:        "role",
}

func (r Role) ToString() string {
	return string(r)
}

// Header Request header
type Header string

var Headers = struct {
	Authorization Header
}{
	Authorization: "Authorization",
}

func (h Header) ToString() string {
	return string(h)
}

// TokenType Authentication header token types
type TokenType string

var TokenTypes = struct {
	Bearer TokenType
}{
	Bearer: "Bearer",
}

func (h TokenType) ToString() string {
	return string(h)
}
