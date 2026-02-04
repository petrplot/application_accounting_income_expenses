package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"fmt"
	"github.com/petrplot/internal/handler"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Загрузка конфигурации
	if err := godotenv.Load(); err != nil {
		log.Println("Инфо: .env файл не найден")
	}

	// 2. Подключение к БД
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"), // Рекомендую добавить это в .env (localhost)
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)

	dbPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Ошибка создания пула БД: %v", err)
	}
	defer dbPool.Close()

	if err := dbPool.Ping(context.Background()); err != nil {
		log.Fatalf("База недоступна: %v", err)
	}

	// 3. Инициализация сервера
	r := gin.Default()

	// Healthcheck
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Регистрация всех маршрутов из пакета handler
	handler.InitRoutes(r, dbPool)

	// 4. Запуск
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Сервер запущен на :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Критическая ошибка сервера: %v", err)
	}
}