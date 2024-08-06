package firebase

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
