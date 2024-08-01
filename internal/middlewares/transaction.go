package middlewares

import (
	"boilerplate-api/internal/config"
	"boilerplate-api/internal/constants"
	"boilerplate-api/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DBTransactionMiddleware struct for transaction
type DBTransactionMiddleware struct {
	logger config.Logger
	db     config.Database
}

// NewDBTransactionMiddleware new instance of transaction
func NewDBTransactionMiddleware(
	logger config.Logger,
	db config.Database,
) DBTransactionMiddleware {
	return DBTransactionMiddleware{
		logger: logger,
		db:     db,
	}
}

// DBTransactionHandle It setup the database transaction middleware
func (m DBTransactionMiddleware) DBTransactionHandle() gin.HandlerFunc {
	m.logger.Info("setting up database transaction middleware")

	return func(c *gin.Context) {
		txHandle := m.db.DB.Begin()
		m.logger.Info("beginning database transaction")

		defer func() {
			if r := recover(); r != nil {
				if err := txHandle.Error; err != nil {
					m.logger.Error("trx commit error: ", err)
				}
				txHandle.Rollback()
			}
		}()

		c.Set(constants.DBTransaction, txHandle)
		c.Next()

		if utils.StatusInList(c.Writer.Status(), []int{http.StatusOK, http.StatusCreated}) {
			m.logger.Info("committing transactions")
			if err := txHandle.Commit().Error; err != nil {
				m.logger.Error("trx commit error: ", err)
			}
		} else {
			m.logger.Info("rolling back transaction due to status code: ", c.Writer.Status())
			txHandle.Rollback()
		}
	}
}
