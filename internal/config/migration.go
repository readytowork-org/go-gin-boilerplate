package config

import (
	"errors"
	"fmt"
	"regexp"

	goMySql "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Migrations Migration Struct
type Migrations struct {
	logger   Logger
	migrator *migrate.Migrate
}

// NewMigrations return new Migrations struct
func NewMigrations(
	dsnConfig DSNConfig,
	logger Logger,
	envPath EnvPath,
) Migrations {
	path := getMigrationFolder(envPath.ToString())
	path = fmt.Sprintf("file://%s/", path)

	mysqlDSNConfig := goMySql.Config{
		User:                 dsnConfig.UserName,
		Passwd:               dsnConfig.Password,
		DBName:               dsnConfig.DBName,
		Net:                  dsnConfig.Network,
		Addr:                 dsnConfig.Address,
		ParseTime:            dsnConfig.ParseTime,
		Loc:                  dsnConfig.TimeLocation,
		AllowNativePasswords: true,
		CheckConnLiveness:    true,
	}

	migrator, err := migrate.New(path, fmt.Sprintf("mysql://%+v", mysqlDSNConfig.FormatDSN()))
	if err != nil {
		logger.Panic("Error in migration: ", err.Error())
	}

	return Migrations{
		logger:   logger,
		migrator: migrator,
	}
}

// MigrateUp migrates all table
func (m Migrations) MigrateUp() {
	m.logger.Info("--- Running Migration Up ---")
	err := m.migrator.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		m.logger.Info("Error in migration steps: ", err.Error())
	}
}

/*
getMigrationFolder path from env path.

e.g:

	../../<.test.env/.env> => ../../migration
	<.test.env/.env> => migration
*/
func getMigrationFolder(envPath string) string {
	m1 := regexp.MustCompile(`(\.(\w+))+`)
	return m1.ReplaceAllString(envPath, "database/migration")
}
