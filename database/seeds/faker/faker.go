package faker

import (
	_ "boilerplate-api/__mocks/mock_data"
	"boilerplate-api/internal/config"
	"gorm.io/gorm"
)

type Config struct {
	SkipClean []string
}

type Faker struct {
	logger   config.Logger
	database *gorm.DB
	config   Config
}

func NewFaker(
	database *gorm.DB,
	logger config.Logger,
	config Config,
) *Faker {
	return &Faker{
		logger:   logger,
		database: database,
		config:   config,
	}
}

func (f *Faker) Seed(data ...interface{}) error {
	txHandle := f.database.Begin()

	for _, _data := range data {
		createQuery := txHandle.Create(_data)

		if err := createQuery.Error; err != nil {
			txHandle.Rollback()
			return err
		}
	}

	txHandle.Commit()
	return nil
}

func (f *Faker) UnSeed() error {
	f.logger.Info("UnSeeding data from database")

	err := f.database.Transaction(
		func(tx *gorm.DB) error {
			err := tx.Exec("SET FOREIGN_KEY_CHECKS = 0;").Error
			if err != nil {
				return err
			}
			var statements []string

			tableQuery := tx.Table("information_schema.tables").
				Select("CONCAT('TRUNCATE TABLE ', table_name, ';')").
				Where("table_schema = DATABASE()")

			if len(f.config.SkipClean) > 0 {
				tableQuery = tableQuery.Where("table_name NOT IN (?)", f.config.SkipClean)
			}

			err = tableQuery.
				Scan(&statements).
				Error
			if err != nil {
				return err
			}

			for _, statement := range statements {
				statementErr := tx.Exec(statement).Error
				if statementErr != nil {
					return err
				}
			}

			err = tx.Exec("SET FOREIGN_KEY_CHECKS = 1;").Error
			if err != nil {
				return err
			}
			return nil
		},
	)
	if err != nil {
		return err
	}
	return nil
}
