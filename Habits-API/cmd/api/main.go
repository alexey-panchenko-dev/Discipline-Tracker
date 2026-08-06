package main

import (
	"log"
	"net/http"

	"habits-api/internal/handler"
	"habits-api/internal/repository"
)

func main() {
	repo := repository.NewTaskRepo()
	taskHandler := handler.NewTaskHandler(repo)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ОК"))
	})

	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
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
	})

	log.Println("Server is running on port:8080")
	http.ListenAndServe(":8080", nil)
}