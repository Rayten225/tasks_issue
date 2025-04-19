package storage

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(constr string) (*Storage, error) {
	db, err := pgxpool.Connect(context.Background(), constr)
	if err != nil {
		return nil, err
	}
	s := Storage{
		db: db,
	}
	return &s, nil
}

type Task struct {
	ID         int
	Opened     int64
	Closed     int64
	AuthorID   int
	AssignedID int
	Title      string
	Content    string
}

func (s *Storage) AllTasks() ([]Task, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT 
			id,
			opened,
			closed,
			author_id,
			assigned_id,
			title,
			content
		FROM tasks
		ORDER BY id;
	`)
	if err != nil {
		return nil, err
	}
	var tasks []Task

	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.ID,
			&t.Opened,
			&t.Closed,
			&t.AuthorID,
			&t.AssignedID,
			&t.Title,
			&t.Content,
		)
		defer rows.Close()
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)

	}
	return tasks, rows.Err()
}

func (s *Storage) Tasks(taskID, authorID int) ([]Task, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT 
			id,
			opened,
			closed,
			author_id,
			assigned_id,
			title,
			content
		FROM tasks
		WHERE
			($1 = 0 OR id = $1) AND
			($2 = 0 OR author_id = $2)
		ORDER BY id;
	`,
		taskID,
		authorID,
	)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var tasks []Task

	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.ID,
			&t.Opened,
			&t.Closed,
			&t.AuthorID,
			&t.AssignedID,
			&t.Title,
			&t.Content,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)

	}
	return tasks, rows.Err()
}

func (s *Storage) NewTask(t Task) (int, error) {
	var id int
	err := s.db.QueryRow(context.Background(), `
		INSERT INTO tasks (title, content)
		VALUES ($1, $2) RETURNING id;
		`,
		t.Title,
		t.Content,
	).Scan(&id)
	return id, err
}

func (s *Storage) UpdateTask(t Task) (int, error) {
	var id int
	err := s.db.QueryRow(context.Background(), `
		UPDATE tasks
		SET title = $1, content = $2 
		WHERE id = $3
		RETURNING id;
		`,
		t.Title,
		t.Content,
		t.ID,
	).Scan(&id)
	return id, err
}

func (s *Storage) DeleteTask(t Task) error {
	_, err := s.db.Exec(context.Background(), `DELETE FROM tasks WHERE id = $1;`, t.ID)
	return err
}
