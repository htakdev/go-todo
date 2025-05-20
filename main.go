package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var todoList *TodoList
var templates = template.Must(template.ParseGlob("templates/*.html"))

func main() {
	var err error
	todoList, err = NewTodoList()
	if err != nil {
		log.Fatalf("TODOリストの初期化エラー: %v", err)
	}

	// 静的ファイルの提供
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// ルートハンドラ
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/todos", handleTodos)

	fmt.Println("サーバーを起動します: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	todos, err := todoList.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = templates.ExecuteTemplate(w, "index.html", todos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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
		w.WriteHeader(http.StatusOK)

	default:
		http.Error(w, "メソッドが許可されていません", http.StatusMethodNotAllowed)
	}
} 