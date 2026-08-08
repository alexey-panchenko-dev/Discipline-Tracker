package handler

import (
	"encoding/json"
	"habits-api/internal/model"
	"habits-api/internal/repository"
	"net/http"
)

// TaskHandler — структура контроллера, хранит ссылку на репозиторий
type TaskHandler struct {
	repo *repository.TaskRepo
}

// CreateTaskInput — структура для создания задачи
type CreateTaskInput struct {
	Title string `json:"title" example:"Выучить Swagger"`
	Desc  string `json:"desc" example:"Разметить хэндлеры аннотациями для генерации документации"`
}

// UpdateTaskInput — структура для частичного обновления задачи
type UpdateTaskInput struct {
	Title *string `json:"title,omitempty" example:"Обновленный заголовок"`
	Desc  *string `json:"desc,omitempty" example:"Обновленное описание задачи"`
}

// NewTaskHandler — конструктор для создания хэндлера
func NewTaskHandler(repo *repository.TaskRepo) *TaskHandler {
	return &TaskHandler{repo: repo}
}

// CreateTaskHandler godoc
// @Summary      Создать новую задачу
// @Description  Создает задачу с переданным заголовком и описанием. Генерирует уникальный UUID.
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        input body CreateTaskInput true "Данные для создания задачи"
// @Success      201  {object}  model.Task
// @Failure      400  {string}  string "Некорректный JSON или пустой заголовок"
// @Failure      405  {string}  string "Метод не поддерживается"
// @Failure      500  {string}  string "Ошибка сохранения задачи в БД"
// @Router       /tasks [post]
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

	createdTask, err := h.repo.CreateTask(input.Title, input.Desc)
	if err != nil {
		http.Error(w, "Ошибка сохранения задачи в БД", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdTask)
}

// GetTasksHandler godoc
// @Summary      Получить список всех задач
// @Description  Возвращает массив всех задач, отсортированных по дате создания (новые первыми).
// @Tags         tasks
// @Produce      json
// @Success      200  {array}   model.Task
// @Failure      405  {string}  string "Метод не поддерживается"
// @Failure      500  {string}  string "Ошибка получения задач из БД"
// @Router       /tasks [get]
func (h *TaskHandler) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	tasks, err := h.repo.GetTasks()
	if err != nil {
		http.Error(w, "Ошибка получения задач из БД", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

// MarkDoneHandler godoc
// @Summary      Отметить задачу как выполненную
// @Description  Переводит флаг `is_done` задачи в статус `true` по её UUID.
// @Tags         tasks
// @Produce      json
// @Param        id   query     string  true  "UUID задачи" format(uuid)
// @Success      200  {object}  model.Task
// @Failure      400  {string}  string "Не указан ID задачи"
// @Failure      404  {string}  string "Задача не найдена"
// @Failure      405  {string}  string "Метод не поддерживается"
// @Router       /tasks [patch]
func (h *TaskHandler) MarkDoneHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPatch && r.Method != http.MethodPut {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Не указан ID задачи", http.StatusBadRequest)
		return
	}

	updatedTask, found := h.repo.MarkDone(idStr)
	if !found {
		http.Error(w, "Задача не найдена", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedTask)
}

// DeleteTaskHandler godoc
// @Summary      Удалить задачу
// @Description  Удаляет задачу из базы данных по её UUID и возвращает удалённый объект.
// @Tags         tasks
// @Produce      json
// @Param        id   query     string  true  "UUID задачи" format(uuid)
// @Success      200  {object}  model.Task
// @Failure      400  {string}  string "Не указан ID задачи"
// @Failure      404  {string}  string "Задача не найдена"
// @Failure      405  {string}  string "Метод не поддерживается"
// @Router       /tasks [delete]
func (h *TaskHandler) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodDelete {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Не указан ID задачи", http.StatusBadRequest)
		return
	}

	deletedTask, isDelete := h.repo.DeleteTask(idStr)
	if !isDelete {
		http.Error(w, "Задача не найдена", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(deletedTask)
}

// UpdateTaskHandler godoc
// @Summary      Обновить заголовок и/или описание задачи
// @Description  Позволяет передать новый title, новый desc или оба поля сразу для задачи с указанным UUID.
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id    query    string           true  "UUID задачи" format(uuid)
// @Param        input body     UpdateTaskInput  true  "Поля для обновления"
// @Success      200   {object} model.Task
// @Failure      400   {string} string "Не указан ID, некорректный JSON или не передано полей"
// @Failure      404   {string} string "Задача не найдена"
// @Router       /tasks [put]
func (h *TaskHandler) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Не указан ID задачи", http.StatusBadRequest)
		return
	}

	var input UpdateTaskInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	if input.Title == nil && input.Desc == nil {
		http.Error(w, "Не передано ни одного поля для обновления", http.StatusBadRequest)
		return
	}

	var updatedTask model.Task
	var found bool

	if input.Title != nil {
		updatedTask, found = h.repo.UpdateTaskTitle(idStr, *input.Title)
		if !found {
			http.Error(w, "Задача не найдена", http.StatusNotFound)
			return
		}
	}

	if input.Desc != nil {
		updatedTask, found = h.repo.UpdateTaskDesc(idStr, *input.Desc)
		if !found {
			http.Error(w, "Задача не найдена", http.StatusNotFound)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedTask)
}