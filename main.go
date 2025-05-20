package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

var todoList *TodoList

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func main() {
	var err error
	todoList, err = NewTodoList()
	if err != nil {
		log.Fatalf("TODOリストの初期化エラー: %v", err)
	}

	// ルートハンドラ
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/todos", enableCORS(handleTodos))

	fmt.Println("サーバーを起動します: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// APIのステータス情報を返す
	status := struct {
		Status    string `json:"status"`
		Version   string `json:"version"`
		Endpoints struct {
			Todos string `json:"todos"`
		} `json:"endpoints"`
	}{
		Status:  "running",
		Version: "1.0.0",
		Endpoints: struct {
			Todos string `json:"todos"`
		}{
			Todos: "/todos",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func handleTodos(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		todos, err := todoList.GetAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(todos)

	case "POST":
		title := r.FormValue("title")
		if title == "" {
			http.Error(w, "タイトルは必須です", http.StatusBadRequest)
			return
		}
		todo, err := todoList.Add(title)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(todo)

	case "PUT":
		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			http.Error(w, "無効なIDです", http.StatusBadRequest)
			return
		}
		completed, err := todoList.Toggle(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]bool{"completed": completed})

	case "DELETE":
		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			http.Error(w, "無効なIDです", http.StatusBadRequest)
			return
		}
		if err := todoList.Delete(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})

	default:
		http.Error(w, "メソッドが許可されていません", http.StatusMethodNotAllowed)
	}
} 