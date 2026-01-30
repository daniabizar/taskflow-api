package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"taskflow-api/internal/models"
	"taskflow-api/internal/utils"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	db *sql.DB
}

func NewTaskHandler(db *sql.DB) *TaskHandler {
	return &TaskHandler{db: db}
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req models.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Set defaults
	if req.Priority == "" {
		req.Priority = models.PriorityMedium
	}
	if req.Category == "" {
		req.Category = models.CategoryPersonal
	}

	var task models.Task
	err := h.db.QueryRow(
		`INSERT INTO tasks (user_id, title, description, priority, category, due_date) 
		 VALUES ($1, $2, $3, $4, $5, $6) 
		 RETURNING id, user_id, title, description, priority, category, is_completed, due_date, created_at, updated_at`,
		userID, req.Title, req.Description, req.Priority, req.Category, req.DueDate,
	).Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Priority,
		&task.Category, &task.IsCompleted, &task.DueDate, &task.CreatedAt, &task.UpdatedAt)

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create task")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Task created successfully", task)
}

func (h *TaskHandler) GetTasks(c *gin.Context) {
	userID := c.GetInt("user_id")

	// Build query with filters
	query := `SELECT id, user_id, title, description, priority, category, is_completed, due_date, created_at, updated_at 
			  FROM tasks WHERE user_id = $1`
	args := []interface{}{userID}
	argCount := 1

	// Filter by priority
	if priority := c.Query("priority"); priority != "" {
		argCount++
		query += " AND priority = $" + strconv.Itoa(argCount)
		args = append(args, priority)
	}

	// Filter by category
	if category := c.Query("category"); category != "" {
		argCount++
		query += " AND category = $" + strconv.Itoa(argCount)
		args = append(args, category)
	}

	// Filter by completion status
	if isCompleted := c.Query("is_completed"); isCompleted != "" {
		argCount++
		query += " AND is_completed = $" + strconv.Itoa(argCount)
		args = append(args, isCompleted == "true")
	}

	// Search in title and description
	if search := c.Query("search"); search != "" {
		argCount++
		query += " AND (title ILIKE $" + strconv.Itoa(argCount) + " OR description ILIKE $" + strconv.Itoa(argCount) + ")"
		args = append(args, "%"+search+"%")
	}

	query += " ORDER BY created_at DESC"

	rows, err := h.db.Query(query, args...)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch tasks")
		return
	}
	defer rows.Close()

	tasks := []models.Task{}
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.UserID, &task.Title, &task.Description,
			&task.Priority, &task.Category, &task.IsCompleted, &task.DueDate,
			&task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			continue
		}
		tasks = append(tasks, task)
	}

	utils.SuccessResponse(c, http.StatusOK, "Tasks retrieved successfully", tasks)
}

func (h *TaskHandler) GetTask(c *gin.Context) {
	userID := c.GetInt("user_id")
	taskID := c.Param("id")

	var task models.Task
	err := h.db.QueryRow(
		`SELECT id, user_id, title, description, priority, category, is_completed, due_date, created_at, updated_at 
		 FROM tasks WHERE id = $1 AND user_id = $2`,
		taskID, userID,
	).Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Priority,
		&task.Category, &task.IsCompleted, &task.DueDate, &task.CreatedAt, &task.UpdatedAt)

	if err == sql.ErrNoRows {
		utils.ErrorResponse(c, http.StatusNotFound, "Task not found")
		return
	}
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch task")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Task retrieved successfully", task)
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	userID := c.GetInt("user_id")
	taskID := c.Param("id")

	var req models.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Build dynamic update query
	query := "UPDATE tasks SET updated_at = CURRENT_TIMESTAMP"
	args := []interface{}{}
	argCount := 0

	if req.Title != nil {
		argCount++
		query += ", title = $" + strconv.Itoa(argCount)
		args = append(args, *req.Title)
	}
	if req.Description != nil {
		argCount++
		query += ", description = $" + strconv.Itoa(argCount)
		args = append(args, *req.Description)
	}
	if req.Priority != nil {
		argCount++
		query += ", priority = $" + strconv.Itoa(argCount)
		args = append(args, *req.Priority)
	}
	if req.Category != nil {
		argCount++
		query += ", category = $" + strconv.Itoa(argCount)
		args = append(args, *req.Category)
	}
	if req.IsCompleted != nil {
		argCount++
		query += ", is_completed = $" + strconv.Itoa(argCount)
		args = append(args, *req.IsCompleted)
	}
	if req.DueDate != nil {
		argCount++
		query += ", due_date = $" + strconv.Itoa(argCount)
		args = append(args, *req.DueDate)
	}

	argCount++
	query += " WHERE id = $" + strconv.Itoa(argCount)
	args = append(args, taskID)

	argCount++
	query += " AND user_id = $" + strconv.Itoa(argCount)
	args = append(args, userID)

	query += " RETURNING id, user_id, title, description, priority, category, is_completed, due_date, created_at, updated_at"

	var task models.Task
	err := h.db.QueryRow(query, args...).Scan(
		&task.ID, &task.UserID, &task.Title, &task.Description, &task.Priority,
		&task.Category, &task.IsCompleted, &task.DueDate, &task.CreatedAt, &task.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		utils.ErrorResponse(c, http.StatusNotFound, "Task not found")
		return
	}
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update task")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Task updated successfully", task)
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	userID := c.GetInt("user_id")
	taskID := c.Param("id")

	result, err := h.db.Exec("DELETE FROM tasks WHERE id = $1 AND user_id = $2", taskID, userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete task")
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		utils.ErrorResponse(c, http.StatusNotFound, "Task not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Task deleted successfully", nil)
}

func (h *TaskHandler) ToggleComplete(c *gin.Context) {
	userID := c.GetInt("user_id")
	taskID := c.Param("id")

	var task models.Task
	err := h.db.QueryRow(
		`UPDATE tasks SET is_completed = NOT is_completed, updated_at = CURRENT_TIMESTAMP 
		 WHERE id = $1 AND user_id = $2 
		 RETURNING id, user_id, title, description, priority, category, is_completed, due_date, created_at, updated_at`,
		taskID, userID,
	).Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Priority,
		&task.Category, &task.IsCompleted, &task.DueDate, &task.CreatedAt, &task.UpdatedAt)

	if err == sql.ErrNoRows {
		utils.ErrorResponse(c, http.StatusNotFound, "Task not found")
		return
	}
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to toggle task")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Task toggled successfully", task)
}

func (h *TaskHandler) GetStats(c *gin.Context) {
	userID := c.GetInt("user_id")

	var stats models.TaskStats

	// Get total tasks
	h.db.QueryRow("SELECT COUNT(*) FROM tasks WHERE user_id = $1", userID).Scan(&stats.Total)

	// Get completed tasks
	h.db.QueryRow("SELECT COUNT(*) FROM tasks WHERE user_id = $1 AND is_completed = true", userID).Scan(&stats.Completed)

	// Get pending tasks
	h.db.QueryRow("SELECT COUNT(*) FROM tasks WHERE user_id = $1 AND is_completed = false", userID).Scan(&stats.Pending)

	// Get high priority tasks
	h.db.QueryRow("SELECT COUNT(*) FROM tasks WHERE user_id = $1 AND priority = 'high'", userID).Scan(&stats.HighPriority)

	// Get overdue tasks
	h.db.QueryRow(
		"SELECT COUNT(*) FROM tasks WHERE user_id = $1 AND is_completed = false AND due_date < $2",
		userID, time.Now(),
	).Scan(&stats.Overdue)

	utils.SuccessResponse(c, http.StatusOK, "Statistics retrieved successfully", stats)
}
