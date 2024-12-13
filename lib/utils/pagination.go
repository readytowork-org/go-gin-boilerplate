package utils

import (
	"github.com/gin-gonic/gin"
)

type IPagination interface {
	Build(c *gin.Context)
}

// BuildPagination -> binds the query and builds pagination
func BuildPagination[T IPagination](c *gin.Context) (m T) {
	_ = c.BindQuery(&m)

	m.Build(c)
	return m
}

// Pagination struct for Pagination
type Pagination struct {
	Sort     string `form:"sort"`
	Keyword  string `form:"keyword"`
	Offset   int    `form:"page,default=1"`
	All      bool   `form:"all"`
	PageSize int    `form:"page_size,default=10"`
}

// Build builds the pagination
func (m *Pagination) Build(c *gin.Context) {
	m.Offset = (m.Offset - 1) * m.PageSize
}
