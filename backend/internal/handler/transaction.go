package handler

import (
	"context"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionHandler struct {
	db *pgxpool.Pool
}

func NewTransactionHandler(db *pgxpool.Pool) *TransactionHandler {
    return &TransactionHandler{
        db: db, // Сохраняем подключение к БД внутри структуры
    }
}

// TransactionRequest — то, что мы ждем от пользователя (JSON)
type TransactionRequest struct {
	Title      string  `json:"title" binding:"required"`      // Наименование
	CategoryID int     `json:"category_id" binding:"required"` // Категория
	Quantity   int     `json:"quantity" binding:"required"`   // Количество
	UnitPrice  float64 `json:"unit_price" binding:"required"` // Цена за 1шт
	Type       string  `json:"type" binding:"required"`       // income или expense
}

// Create — метод для создания записи
func (h *TransactionHandler) Create(c *gin.Context) {
	var req TransactionRequest

	// 1. Читаем JSON из запроса
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные входные данные"})
		return
	}

	// 2. Расчет общей суммы (твое требование: сумма за все покупки)
	totalAmount := float64(req.Quantity) * req.UnitPrice

	// 3. SQL запрос для вставки в PostgreSQL
	query := `
		INSERT INTO transactions (title, category_id, quantity, unit_price, total_amount, type, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
		RETURNING id, created_at`

	var id int
	var createdAt interface{}

	err := h.db.QueryRow(context.Background(), query,
		req.Title, req.CategoryID, req.Quantity, req.UnitPrice, totalAmount, req.Type,
	).Scan(&id, &createdAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка базы данных: " + err.Error()})
		return
	}

	// 4. Возвращаем результат
	c.JSON(http.StatusCreated, gin.H{
		"id":           id,
		"status":       "success",
		"total_amount": totalAmount,
		"date":         createdAt,
	})
}

// GetAll — получить все транзакции
func (h *TransactionHandler) GetAll(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Список всех транзакций"})
}

// GetOne — получить одну транзакцию по ID
func (h *TransactionHandler) GetOne(c *gin.Context) {
	id := c.Param("id") // извлекаем ID из URL
	c.JSON(http.StatusOK, gin.H{"message": "Транзакция ID: " + id})
}

// Update — частично изменить запись
func (h *TransactionHandler) Update(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Обновление транзакции ID: " + id})
}

// Delete — удалить запись
func (h *TransactionHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Удаление транзакции ID: " + id})
}
