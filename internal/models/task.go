package models

type Task struct {
    ID          int32  `db:"id"`          // Task ID
    Title       string `db:"title"`       // Task title
    Description string `db:"description"` // Task description
    Completed   bool   `db:"completed"`   // Task completion status
}

