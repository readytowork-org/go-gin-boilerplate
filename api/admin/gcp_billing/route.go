package gcp_billing

import (
	"boilerplate-api/internal/router"
)

// SetupRoutes -> sets up route for util entities
func SetupRoutes(
	controller Controller,
	router router.Router,
) {
	gcpBilling := router.V1.Group("/gcp-billing")
	{
		gcpBilling.GET("", controller.GetCost)
		gcpBilling.GET("budget", controller.GetBudgetInfo)
		gcpBilling.POST("budget", controller.CreateUpdateBudget)
	}
}
