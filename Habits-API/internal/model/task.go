package model

import "time"

type Task struct {
    ID        string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
    Title     string    `json:"title" example:"Выучить Swagger"`
    Desc      string    `json:"desc" example:"Разметить хэндлеры аннотациями"`
    IsDone    bool      `json:"is_done" example:"false"`
    CreatedAt time.Time `json:"created_at" example:"2026-08-09T00:00:00Z"`
}