package main

import (
	_ "habits-api/docs"

	httpSwagger "github.com/swaggo/http-swagger"

	"log"
	"net/http"

	"habits-api/internal/handler"
	"habits-api/internal/repository"
)

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// @title           Habits & Tasks API
// @version         1.0
// @host            localhost:8080
// @BasePath        /
func main() {
	db := repository.InitDB()
	defer db.Close()

	repo := repository.NewTaskRepo(db)
	taskHandler := handler.NewTaskHandler(repo)

	http.HandleFunc("/health", enableCORS(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ОК"))
	}))

	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	http.HandleFunc("/tasks", enableCORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			taskHandler.CreateTaskHandler(w, r)
		case http.MethodGet:
			taskHandler.GetTasksHandler(w, r)
		case http.MethodPut:
			taskHandler.UpdateTaskHandler(w, r)
		case http.MethodPatch:
			taskHandler.MarkDoneHandler(w, r)
		case http.MethodDelete:
			taskHandler.DeleteTaskHandler(w, r)
		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	}))

	log.Println("Server is running on port:8080")
	http.ListenAndServe(":8080", nil)
}