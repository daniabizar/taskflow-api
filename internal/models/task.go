package models

import "time"

type Priority string
type Category string

const (
	PriorityHigh   Priority = "high"
	PriorityMedium Priority = "medium"
	PriorityLow    Priority = "low"
)

const (
	CategoryPersonal Category = "personal"
	CategoryWork     Category = "work"
	CategoryUrgent   Category = "urgent"
)

type Task struct {
	ID          int        `json:"id"`
	UserID      int        `json:"user_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Priority    Priority   `json:"priority"`
	Category    Category   `json:"category"`
	IsCompleted bool       `json:"is_completed"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Struct untuk request create task
type CreateTaskRequest struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	Priority    Priority   `json:"priority" binding:"omitempty,oneof=high medium low"`
	Category    Category   `json:"category" binding:"omitempty,oneof=personal work urgent"`
	DueDate     *time.Time `json:"due_date,omitempty"`
}

// Struct untuk request update task
type UpdateTaskRequest struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Priority    *Priority  `json:"priority" binding:"omitempty,oneof=high medium low"`
	Category    *Category  `json:"category" binding:"omitempty,oneof=personal work urgent"`
	IsCompleted *bool      `json:"is_completed"`
	DueDate     *time.Time `json:"due_date,omitempty"`
}

// Struct untuk task statistics
type TaskStats struct {
	Total        int `json:"total"`
	Completed    int `json:"completed"`
	Pending      int `json:"pending"`
	HighPriority int `json:"high_priority"`
	Overdue      int `json:"overdue"`
}
