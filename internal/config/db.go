package config

import (
	"database/sql"
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
}

// NewDatabase creates a new database instance
func NewDatabase(logger Logger, dsnConfig DSNConfig) Database {
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
	logger.Infof("DSN :: %+v\n", mysqlDSNConfig)

	db, err := gorm.Open(
		mysql.New(mysql.Config{DSNConfig: &mysqlDSNConfig}),
		&gorm.Config{Logger: logger.GetGormLogger()},
	)
	if err != nil {
		logger.Panicf("Database connection failed :: %+v\n", err)
	}

	logger.Info("creating database if it doesn't exist")
	if err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dsnConfig.DBName)).Error; err != nil {
		logger.Info("couldn't create database")
		logger.Panic(err)
	}
	logger.Info("using given database")
	if err = db.Exec(fmt.Sprintf("USE %s", dsnConfig.DBName)).Error; err != nil {
		logger.Info("cannot use the given database")
		logger.Panic(err)
	}

	logger.Infof("Database connection established : %s", db.Migrator().CurrentDatabase())

	return Database{
		db,
	}
}

// MockDatabase modal
type MockDatabase struct {
	Database
	sqlDb   *sql.DB
	sqlMock sqlmock.Sqlmock
}

func NewMockDatabase() Database {
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
		log.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	sqlMock.ExpectBegin()
	sqlMock.ExpectCommit()
	sqlMock.ExpectRollback()

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDb,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		DryRun: true,
		Logger: newLogger,
	})

	//var columns = models.UserName{}.Columns()
	//// Add mock data to mock sql columns
	//now := time.Now()
	//rows := sqlMock.NewRows(columns).
	//	AddRow(now, now, nil, 1, "firebase_id", "jhon@mailinator.com", constants.UnVerifiedEmail, nil, "Jhon", "male", 20, "january", "93234242342", "234234").
	//	AddRow(now, now, nil, 1, "firebase_id", "jhon@mailinator.com", constants.UnVerifiedEmail, nil, "Jhon", "male", 20, "january", "93234242342", "234234")

	if err != nil {
		log.Fatalf("An error '%s' was not expected when opening gorm database", err)
	}

	return Database{
		DB: db,
	}
}
