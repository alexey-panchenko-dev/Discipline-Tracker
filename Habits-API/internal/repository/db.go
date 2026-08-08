package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	// Драйвер PostgreSQL: регистрирует себя в "database/sql" под именем "postgres"
	_ "github.com/lib/pq"
)

// InitDB настраивает подключение к БД в Docker и создает нужную таблицу
func InitDB() *sql.DB {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	// Строка подключения: где ищем базу, под каким юзером и с каким паролем
	connStr := fmt.Sprintf("host=%s port=5432 user=postgres password=mysecretpassword dbname=habits_db sslmode=disable", dbHost)

	// Инициализируем объект базы (пока просто конфиг, реального соединения ещё нет)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Ошибка подключения к DB: %v", err)
	}

	// Делаем реальный запрос в базу, чтобы проверить, поднята ли она
	if err := db.Ping(); err != nil {
		log.Fatalf("Не удалось подключиться к БД (Ping failed): %v", err)
	}

	log.Println("Успешное подключение к DB")

	// SQL-запрос для автоматического создания таблицы при старте, если её ещё нет
	createTableSql := `
	CREATE TABLE IF NOT EXISTS tasks (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		title VARCHAR(255) NOT NULL,
		desc_text TEXT DEFAULT '',
		is_done BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`

	// Выполняем SQL-запрос без ожидания возврата строк
	_, err = db.Exec(createTableSql)
	if err != nil {
		log.Fatalf("Ошибка создания таблицы %v", err)
	}

	// Возвращаем пуллер/клиент базы для работы из других частей приложения
	return db
}