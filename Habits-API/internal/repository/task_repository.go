package repository

import (
	"habits-api/internal/model" // импортируем типы для обозначения того что мы доставем и кладем в памяти
	"time"                      // импортируем тип время для извлечения момента создания задачи
)

//типы нашего хранилища
type TaskRepo struct {
	Tasks []model.Task // массив с структурами данных Task
	NextID int64 // этот айди нужен для того что бы генерировать айди для тасок
}

//функция для создания репозитория и хранения в нем тасок
func NewTaskRepo() *TaskRepo { //обьявление репозитория (ака хранилище), *TaskRepo говорим о том что функция возвращает не копию а ссылку на нее
	return &TaskRepo{ // создание новой структуры
		Tasks: make([]model.Task, 0), // таски: создать(структуру Тасок, там пока ничего нет)
		NextID: 1, // первая таска будет с айди: 1
	}
}

func (r *TaskRepo) CreateTask(title, desc string) model.Task { // создаю обращение как к методу, название функции(принимает название и описание, это будут строки) используется структура Task
	
	task := model.Task{ // создаем таску с нужными параметрами
		ID: r.NextID,
		Title: title,
		Desc: desc,
		IsDone: false,
		CreatedAt: time.Now(),
	}
	r.Tasks = append(r.Tasks, task) // добавляю новую таску к имеющимся
	
	r.NextID += 1 // что бы след таска была 2, 3, 4 и тд
	
	return task // возвращаю таск
}

func (r *TaskRepo) GetTasks() []model.Task{
	return r.Tasks
}

func (r *TaskRepo) MarkDone(id int64) (model.Task, bool) {
	for i, task := range r.Tasks {
		if task.ID == id {
			r.Tasks[i].IsDone = true
			return r.Tasks[i], true
		}
	}

	return model.Task{}, false
}

func (r *TaskRepo) DeleteTask(id int64) (model.Task, bool) {
	for i, task := range r.Tasks {
		if task.ID == id {
			r.Tasks = append(r.Tasks[:i], r.Tasks[i+1:]...)
			return task, true
		}
	}
	return model.Task{}, false
}

func (r *TaskRepo) UpdateTaskTitle(id int64, title string) (model.Task, bool) {
	for i, task := range r.Tasks {
		if task.ID == id {
			r.Tasks[i].Title = title
			return r.Tasks[i], true
		}
	}
	return model.Task{}, false
}

func (r *TaskRepo) UpdateTaskDesc(id int64, desc string) (model.Task, bool) {
	for i, task := range r.Tasks {
		if task.ID == id {
			r.Tasks[i].Title = desc
			return r.Tasks[i], true
		}
	}
	return model.Task{}, false
}