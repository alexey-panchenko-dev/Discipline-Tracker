package repository

import (
	"database/sql"
	"habits-api/internal/model" // импортируем типы для обозначения того что мы доставем и кладем в памяти
)

//типы нашего хранилища
type TaskRepo struct {
	db *sql.DB
}

//функция для создания репозитория и хранения в нем тасок
func NewTaskRepo(db *sql.DB) *TaskRepo { //обьявление репозитория (ака хранилище), *TaskRepo говорим о том что функция возвращает не копию а ссылку на нее
	return &TaskRepo{db: db}
}

func (r *TaskRepo) CreateTask(title, desc string) (model.Task, error) {
	query := `
	INSERT INTO tasks (title, desc_text)
	VALUES ($1, $2)
	RETURNING id, title, desc_text, is_done, created_at
	`

	var task model.Task
	err := r.db.QueryRow(query, title, desc).Scan(
		&task.ID,
		&task.Title,
		&task.Desc,
		&task.IsDone,
		&task.CreatedAt,
	)

	return task, err
}

func (r *TaskRepo) GetTasks() ([]model.Task, error) {
	query := `SELECT id, title, desc_text, is_done, created_at FROM tasks ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]model.Task, 0)

	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Desc, &task.IsDone, &task.CreatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *TaskRepo) MarkDone(id string) (model.Task, bool) {
	query := `
		UPDATE tasks
		SET is_done = TRUE
		WHERE id = $1
		RETURNING id, title, desc_text, is_done, created_at
	`

	var task model.Task
	err := r.db.QueryRow(query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Desc,
		&task.IsDone,
		&task.CreatedAt,
	)

	if err != nil {
		return model.Task{}, false
	}

	return task, true
}

func (r *TaskRepo) DeleteTask(id string) (model.Task, bool) {
	query := `
		DELETE FROM tasks 
		WHERE id = $1 
		RETURNING id, title, desc_text, is_done, created_at
	`

	var task model.Task

	err := r.db.QueryRow(query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Desc,
		&task.IsDone,
		&task.CreatedAt,
	)

	if err != nil {
		return model.Task{}, false
	}
	return task, true
}


func (r *TaskRepo) UpdateTaskTitle(id string, title string) (model.Task, bool) {
	query := `
		UPDATE tasks 
		SET title = $1 
		WHERE id = $2 
		RETURNING id, title, desc_text, is_done, created_at
	`

	var task model.Task

	err := r.db.QueryRow(query, title, id).Scan(
		&task.ID,
		&task.Title,
		&task.Desc,
		&task.IsDone,
		&task.CreatedAt,
	)

	if err != nil {
		return model.Task{}, false
	}
	return task, true
}

func (r *TaskRepo) UpdateTaskDesc(id string, desc string) (model.Task, bool) {
	query := `
		UPDATE tasks 
		SET desc_text = $1 
		WHERE id = $2 
		RETURNING id, title, desc_text, is_done, created_at
	`

	var task model.Task
	err := r.db.QueryRow(query, desc, id).Scan(
		&task.ID,
		&task.Title,
		&task.Desc,
		&task.IsDone,
		&task.CreatedAt,
	)

	if err != nil {
		return model.Task{}, false
	}
	return task, true
}