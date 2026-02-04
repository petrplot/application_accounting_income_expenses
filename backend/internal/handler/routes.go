package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	
)

// InitRoutes собирает все маршруты приложения в одном месте
func InitRoutes(r *gin.Engine, db *pgxpool.Pool) {
	// Инициализируем хендлеры
	transH := NewTransactionHandler(db)
	// budgetH := NewBudgetHandler(db)
	// reportH := NewReportHandler(db)

	api := r.Group("/api/v1")
	{
		// Группа транзакций
		t := api.Group("/transactions")
		{
			t.GET("", transH.GetAll)
			t.GET("/:id", transH.GetOne)
			t.POST("", transH.Create)
			t.PATCH("/:id", transH.Update)
			t.DELETE("/:id", transH.Delete)
		}

		// Группа отчетов
		// reports := api.Group("/reports")
		// {
		//    reports.GET("/summary", reportH.GetSummary)
		// }
	}
}
