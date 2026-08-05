package model

import "time"

type Task struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Desc      string `json:"desc"`
	IsDone    bool   `json:"is_done"`
	CreatedAt time.Time `json:"created_at"`
}

// описываю структуру для задач что бы го мог выделить нужное колво места в памяти для каждой таски