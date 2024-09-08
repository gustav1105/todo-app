package handlers

import (
    "context"
    "errors"
    "testing"
    "todo-app/proto"

    "github.com/DATA-DOG/go-sqlmock"
    "github.com/stretchr/testify/assert"
)

// Mock sql.Result implementation
type MockResult struct{}

func (r *MockResult) LastInsertId() (int64, error) { return 0, nil }
func (r *MockResult) RowsAffected() (int64, error) { return 1, nil }

func TestAddTask_Success(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Error initializing sqlmock: %v", err)
    }
    defer db.Close()

    handler := NewTodoServiceHandler(db)
    ctx := context.Background()

    taskReq := &proto.TaskRequest{
        Title:       "Test task",
        Description: "Test description",
    }

    mock.ExpectExec("INSERT INTO tasks").
        WithArgs("Test task", "Test description", false).
        WillReturnResult(sqlmock.NewResult(1, 1))

    _, err = handler.AddTask(ctx, taskReq)
    assert.NoError(t, err)

    err = mock.ExpectationsWereMet()
    assert.NoError(t, err)
}

func TestAddTask_Failure(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Error initializing sqlmock: %v", err)
    }
    defer db.Close()

    handler := NewTodoServiceHandler(db)
    ctx := context.Background()

    taskReq := &proto.TaskRequest{
        Title:       "Test task",
        Description: "Test description",
    }

    mock.ExpectExec("INSERT INTO tasks").
        WithArgs("Test task", "Test description", false).
        WillReturnError(errors.New("DB error"))

    _, err = handler.AddTask(ctx, taskReq)
    assert.Error(t, err)
    assert.Equal(t, "DB error", err.Error())

    err = mock.ExpectationsWereMet()
    assert.NoError(t, err)
}

func TestGetTasks_Success(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Error initializing sqlmock: %v", err)
    }
    defer db.Close()

    handler := NewTodoServiceHandler(db)
    ctx := context.Background()

    rows := sqlmock.NewRows([]string{"id", "title", "description", "completed"}).
        AddRow(1, "Test task", "Test description", false)

    mock.ExpectQuery("SELECT id, title, description, completed FROM tasks").
        WillReturnRows(rows)

    _, err = handler.GetTasks(ctx, &proto.Empty{})
    assert.NoError(t, err)

    err = mock.ExpectationsWereMet()
    assert.NoError(t, err)
}
func TestCompleteTask_Success(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Error initializing sqlmock: %v", err)
    }
    defer db.Close()

    handler := NewTodoServiceHandler(db)
    ctx := context.Background()

    req := &proto.CompleteTaskRequest{
        Id: 1,
    }

    // Update the query with proper escaping for '?'
    mock.ExpectExec(`UPDATE tasks SET completed = \? WHERE id = \?`).
        WithArgs(true, int64(1)).
        WillReturnResult(sqlmock.NewResult(1, 1))

    _, err = handler.CompleteTask(ctx, req)
    assert.NoError(t, err)

    err = mock.ExpectationsWereMet()
    assert.NoError(t, err)
}

func TestCompleteTask_Failure(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Error initializing sqlmock: %v", err)
    }
    defer db.Close()

    handler := NewTodoServiceHandler(db)
    ctx := context.Background()

    req := &proto.CompleteTaskRequest{
        Id: 1,
    }

    // Update the query with proper escaping for '?'
    mock.ExpectExec(`UPDATE tasks SET completed = \? WHERE id = \?`).
        WithArgs(true, int64(1)).
        WillReturnError(errors.New("DB error"))

    _, err = handler.CompleteTask(ctx, req)
    assert.Error(t, err)
    assert.Equal(t, "DB error", err.Error())

    err = mock.ExpectationsWereMet()
    assert.NoError(t, err)
}

