package config

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	goMySql "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database modal
type Database struct {
	*gorm.DB
	dsn             string
	dbType          DBType
	ConnectionError *error
}

// NewDatabase creates a new database instance
func NewDatabase(logger Logger, dsnConfig DSNConfig) *Database {
	database := Database{
		dbType: dsnConfig.DBType,
	}

	var dialector *gorm.Dialector

	switch database.dbType {
	case DBTypeSql:
		mysqlDSNConfig := goMySql.Config{
			User:                 dsnConfig.UserName,
			Passwd:               dsnConfig.Password,
			Net:                  dsnConfig.Network,
			Addr:                 dsnConfig.Address,
			ParseTime:            dsnConfig.ParseTime,
			Loc:                  dsnConfig.TimeLocation,
			AllowNativePasswords: true,
			CheckConnLiveness:    true,
		}
		_dialector := mysql.New(
			mysql.Config{
				DSNConfig: &mysqlDSNConfig,
			},
		)
		dialector = &_dialector
		mysqlDSNConfig.DBName = dsnConfig.DBName
		database.dsn = mysqlDSNConfig.FormatDSN()
		break
	}

	if dialector == nil || database.dsn == "" || dsnConfig.Address == "" {
		err := errors.New("database not configured --- Using Mock Database")
		logger.Error(err)

		database = *NewMockDatabase()
		database.ConnectionError = &err
		return &database
	}

	db, err := gorm.Open(*dialector, &gorm.Config{Logger: logger.GetGormLogger()})
	if err != nil {
		_err := errors.New(
			fmt.Sprintf(
				"Database connection failed\n Please check dsn:: %+v\n+%v", database.dsn, err.Error(),
			),
		)
		logger.Error(_err)
		database.ConnectionError = &_err
		return &database
	}
	database.DB = db

	logger.Info("creating database if it doesn't exist")
	if err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dsnConfig.DBName)).Error; err != nil {
		logger.Error("Couldn't create database")
		database.ConnectionError = &err
		return &database
	}

	logger.Info("using given database")
	if err = db.Exec(fmt.Sprintf("USE %s", dsnConfig.DBName)).Error; err != nil {
		logger.Error("Cannot use the given database")
		database.ConnectionError = &err
		return &database
	}

	logger.Infof("Database connection established : %s", db.Migrator().CurrentDatabase())

	return &database
}

func (d Database) DSN() string {
	return d.dsn
}

func (d Database) Type() string {
	return d.dbType.ToString()
}

// MockDatabase modal
type MockDatabase struct {
	*Database
	sqlDb   *sql.DB
	sqlMock sqlmock.Sqlmock
}

func NewMockDatabase() *Database {
	var connectionError *error
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Disable color
		},
	)

	sqlDb, sqlMock, err := sqlmock.New()
	if err != nil {
		globalLog.Error("error opening a stub database connection")
		connectionError = &err
	}
	sqlMock.ExpectBegin()
	sqlMock.ExpectCommit()
	sqlMock.ExpectRollback()

	db, err := gorm.Open(
		mysql.New(
			mysql.Config{
				Conn:                      sqlDb,
				SkipInitializeWithVersion: true,
			},
		), &gorm.Config{
			DryRun: true,
			Logger: newLogger,
		},
	)

	if err != nil {
		globalLog.Errorf("An error was not expected when opening gorm database")
		connectionError = &err
	}

	return &Database{
		DB:              db,
		dsn:             "",
		dbType:          DBTypeSql,
		ConnectionError: connectionError,
	}
}
