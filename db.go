package main

import (
	"math/rand"
	"time"

	"github.com/cockroachdb/ddshop/robustdb"
)

func connectDB(urls []string) (*robustdb.DB, error) {
	db, err := robustdb.New(urls...)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

const schema = `
CREATE TABLE IF NOT EXISTS todos (
	id int4 PRIMARY KEY,
	title text,
	created_at timestamptz,
	completed bool
)
`

func bootstrapDB(db *robustdb.DB) error {
	_, err := db.Exec(schema)
	return err
}

type todo struct {
	ID        int32     `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	Completed bool      `json:"completed"`
}

func upsertTodo(db *robustdb.DB, t todo) error {
	if t.ID == 0 {
		t.ID = rand.Int31()
	}
	_, err := db.Exec(
		`INSERT INTO todos (id, title, created_at, completed)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET
			title = excluded.title,
			created_at = excluded.created_at,
			completed = excluded.completed`,
		t.ID, t.Title, t.CreatedAt, t.Completed)
	return err
}

func listTodos(db *robustdb.DB) ([]todo, error) {
	rows, err := db.Query(`SELECT id, title, created_at, completed FROM todos`)
	if err != nil {
		return nil, err
	}
	todos := []todo{}
	for rows.Next() {
		var t todo
		if err := rows.Scan(&t.ID, &t.Title, &t.CreatedAt, &t.Completed); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return todos, nil
}

func deleteTodo(db *robustdb.DB, todoID int64) error {
	_, err := db.Exec(`DELETE FROM todos WHERE id = $1`, todoID)
	return err
}
