package storage

import (
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "postgres"
	dbname = "task_issue"
)

func TestStorage(t *testing.T) {
	pwd := os.Getenv("dbpass")
	connStr := "postgres://" + user + ":" + pwd + "@" + host + ":" + strconv.Itoa(port) + "/" + dbname
	db, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Создание таблицы для теста
	_, err = db.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS tasks (
			id SERIAL PRIMARY KEY,
			opened BIGINT,
			closed BIGINT,
			author_id INT,
			assigned_id INT,
			title TEXT,
			content TEXT
		);
	`)
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}
	defer db.Exec(context.Background(), `DROP TABLE tasks;`)

	// Инициализация хранилища
	store, err := New(connStr)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}

	// Тест 1: Создание новой задачи
	task := Task{
		Title:   "Test Task",
		Content: "This is a test task",
	}
	id, err := store.NewTask(task)
	assert.NoError(t, err, "NewTask should not return error")
	assert.Greater(t, id, 0, "NewTask should return a valid ID")

	// Тест 2: Получение всех задач
	tasks, err := store.AllTasks()
	assert.NoError(t, err, "AllTasks should not return error")
	assert.Len(t, tasks, 1, "AllTasks should return exactly one task")
	assert.Equal(t, "Test Task", tasks[0].Title, "Task title should match")
	assert.Equal(t, "This is a test task", tasks[0].Content, "Task content should match")

	// Тест 3: Удаление задачи
	err = store.DeleteTask(Task{ID: id})
	assert.NoError(t, err, "DeleteTask should not return error")

	// Тест 4: Проверка, что задача удалена
	tasks, err = store.AllTasks()
	assert.NoError(t, err, "AllTasks should not return error")
	assert.Len(t, tasks, 0, "AllTasks should return no tasks after deletion")
}
