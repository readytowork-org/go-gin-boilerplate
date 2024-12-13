package seeds

import (
	"context"

	"boilerplate-api/lib/config"
	"boilerplate-api/lib/constants"
	"boilerplate-api/services"
)

// AdminSeed  Admin seeding
type AdminSeed struct {
	Seed
	logger          config.Logger
	firebaseService services.IFirebaseAdminSeed
	adminEmail      string
	adminPass       string
	adminName       string
}

// NewAdminSeed creates admin seed
func NewAdminSeed(
	logger config.Logger,
	authService services.IFirebaseAdminSeed,
	adminEmail string,
	adminPass string,
	adminName string,
) AdminSeed {
	return AdminSeed{
		logger:          logger,
		firebaseService: authService,
		adminEmail:      adminEmail,
		adminPass:       adminPass,
		adminName:       adminName,
	}
}

// Run the seed data
func (c AdminSeed) Run() {
	c.logger.Info("🌱 seeding admin data...")

	_, err := c.firebaseService.GetUserByEmail(context.Background(), c.adminEmail)
	if err != nil {
		_, errResponse := c.firebaseService.CreateUser(
			c.adminName, c.adminEmail, c.adminPass,
			string(constants.Roles.SuperAdmin),
		)
		if errResponse != nil {
			c.logger.Error("Firebase Admin user can't be created: ", errResponse.Message)
			return
		}

		c.logger.Info("Firebase Admin UserName Created, email: ", c.adminEmail, " password: ", c.adminPass)
	}

	c.logger.Info("Admin already exist")
}
