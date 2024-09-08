package handlers

import (
    "context"
    "database/sql"
    "log"
    "todo-app/internal/models"
    "todo-app/proto"
)

// TodoServiceHandler struct to hold DB reference and embed UnimplementedTodoServiceServer
type TodoServiceHandler struct {
    proto.UnimplementedTodoServiceServer
    DB *sql.DB
}

// NewTodoServiceHandler constructor with DB dependency
func NewTodoServiceHandler(db *sql.DB) proto.TodoServiceServer {
    return &TodoServiceHandler{
        DB: db,
    }
}

// AddTask implementation to add a task to the database
func (h *TodoServiceHandler) AddTask(ctx context.Context, req *proto.TaskRequest) (*proto.Empty, error) {
    query := "INSERT INTO tasks (title, description, completed) VALUES (?, ?, ?)"
    _, err := h.DB.ExecContext(ctx, query, req.Title, req.Description, false)
    if err != nil {
        log.Printf("Error adding task: %v", err)
        return nil, err
    }
    return &proto.Empty{}, nil
}

// GetTasks implementation to fetch all tasks from the database
func (h *TodoServiceHandler) GetTasks(ctx context.Context, req *proto.Empty) (*proto.TaskResponse, error) {
    query := "SELECT id, title, description, completed FROM tasks"
    rows, err := h.DB.QueryContext(ctx, query)
    if err != nil {
        log.Printf("Error fetching tasks: %v", err)
        return nil, err
    }
    defer rows.Close()

    var tasks []*proto.Task
    for rows.Next() {
        var task models.Task
        err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Completed)
        if err != nil {
            log.Printf("Error scanning task: %v", err)
            return nil, err
        }
        tasks = append(tasks, &proto.Task{
            Id:          task.ID,
            Title:       task.Title,
            Description: task.Description,
            Completed:   task.Completed,
        })
    }

    return &proto.TaskResponse{Tasks: tasks}, nil
}

// CompleteTask implementation to mark a task as completed
func (h *TodoServiceHandler) CompleteTask(ctx context.Context, req *proto.CompleteTaskRequest) (*proto.Empty, error) {
    query := "UPDATE tasks SET completed = ? WHERE id = ?"
    _, err := h.DB.ExecContext(ctx, query, true, req.Id)
    if err != nil {
        log.Printf("Error updating task: %v", err)
        return nil, err
    }
    return &proto.Empty{}, nil
}

