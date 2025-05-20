package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// Todo はTODOアイテムを表す構造体です
type Todo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

// TodoList はTODOアイテムのリストを管理する構造体です
type TodoList struct {
	db *sql.DB
}

// NewTodoList は新しいTodoListを作成します
func NewTodoList() (*TodoList, error) {
	// PostgreSQLの接続情報
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=todo sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("データベース接続エラー: %v", err)
	}

	// 接続テスト
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("データベース接続テストエラー: %v", err)
	}

	// テーブル作成
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS todos (
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			completed BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("テーブル作成エラー: %v", err)
	}

	return &TodoList{db: db}, nil
}

// Add は新しいTODOを追加します
func (tl *TodoList) Add(title string) (Todo, error) {
	var todo Todo
	err := tl.db.QueryRow(
		"INSERT INTO todos (title) VALUES ($1) RETURNING id, title, completed, created_at",
		title,
	).Scan(&todo.ID, &todo.Title, &todo.Completed, &todo.CreatedAt)

	if err != nil {
		return Todo{}, fmt.Errorf("TODO追加エラー: %v", err)
	}
	return todo, nil
}

// GetAll は全てのTODOを返します
func (tl *TodoList) GetAll() ([]Todo, error) {
	rows, err := tl.db.Query("SELECT id, title, completed, created_at FROM todos ORDER BY created_at DESC")
	if err != nil {
		return nil, fmt.Errorf("TODO取得エラー: %v", err)
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Completed, &todo.CreatedAt); err != nil {
			return nil, fmt.Errorf("TODO読み込みエラー: %v", err)
		}
		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("TODO読み込みエラー: %v", err)
	}

	return todos, nil
}

// Toggle はTODOの完了状態を切り替えます
func (tl *TodoList) Toggle(id int) (bool, error) {
	var completed bool
	err := tl.db.QueryRow(
		"UPDATE todos SET completed = NOT completed WHERE id = $1 RETURNING completed",
		id,
	).Scan(&completed)

	if err != nil {
		return false, fmt.Errorf("TODO更新エラー: %v", err)
	}
	return completed, nil
}

// Delete はTODOを削除します
func (tl *TodoList) Delete(id int) error {
	result, err := tl.db.Exec("DELETE FROM todos WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("TODO削除エラー: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("削除結果確認エラー: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("指定されたIDのTODOが見つかりません")
	}

	return nil
} 