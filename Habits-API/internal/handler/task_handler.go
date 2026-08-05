package handler

import (
	"encoding/json"
	"habits-api/internal/repository" // импортируем пакет репозитория, чтобы хэндлер умел обращаться к хранилищу задач
	"net/http"
	"strconv"
)

// структура хэндлера (контроллера)
type TaskHandler struct {
    repo *repository.TaskRepo // ссылка (указатель) на наше хранилище; хэндлер сам не хранит данные, а только управляет ими через repo
}

type CreateTaskInput struct {
	Title string `json:"title"`
	Desc  string `json:"desc"`
}

// функция-конструктор для создания хэндлера
func NewTaskHandler(repo *repository.TaskRepo) *TaskHandler { // передаем существующий репозиторий в виде указателя
    return &TaskHandler{repo: repo} // создаем и возвращаем указатель на готовый TaskHandler с привязанным к нему репозиторием
}

func (h *TaskHandler) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    if r.Method != http.MethodPost {
        http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
        return
    }

    var input CreateTaskInput
    err := json.NewDecoder(r.Body).Decode(&input)
    if err != nil {
        http.Error(w, "Некорректный JSON в теле запроса", http.StatusBadRequest)
        return 
    }

    if input.Title == "" {
        http.Error(w, "Заголовок задачи не может быть пустым", http.StatusBadRequest)
        return
    }

    createdTask := h.repo.CreateTask(input.Title, input.Desc)

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(createdTask)
}

func (h *TaskHandler) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
    if r.Method != http.MethodGet {
        http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
        return
    }

    tasks := h.repo.GetTasks()

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandler) MarkDoneHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    if r.Method != http.MethodPatch && r.Method != http.MethodPut {
        http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
        return
    }

    idStr := r.URL.Query().Get("id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        http.Error(w, "Некорректный ID", http.StatusBadRequest)
        return
    }

    updatedTask, found := h.repo.MarkDone(id)

    if !found {
        http.Error(w, "Задача не найдена", http.StatusNotFound)
        return
    } else {
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(updatedTask)
    }
}