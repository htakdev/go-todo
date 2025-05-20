package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleIndex(t *testing.T) {
	// テスト用のリクエストを作成
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// レスポンスレコーダーを作成
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleIndex)

	// リクエストを実行
	handler.ServeHTTP(rr, req)

	// ステータスコードの確認
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// レスポンスのContent-Typeの確認
	expectedContentType := "application/json"
	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("handler returned wrong content type: got %v want %v",
			contentType, expectedContentType)
	}

	// レスポンスのJSONをパース
	var response map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	// 期待されるフィールドの確認
	expectedFields := []string{"status", "version", "endpoints"}
	for _, field := range expectedFields {
		if _, exists := response[field]; !exists {
			t.Errorf("response missing expected field: %s", field)
		}
	}

	// ステータスの値の確認
	if status, ok := response["status"].(string); !ok || status != "running" {
		t.Errorf("handler returned wrong status: got %v want %v",
			status, "running")
	}
}

func TestHandleTodos(t *testing.T) {
	// テスト用のTODOリストを初期化
	var err error
	todoList, err = NewTodoList()
	if err != nil {
		t.Fatal(err)
	}

	// TODOの追加テスト
	t.Run("Add Todo", func(t *testing.T) {
		formData := bytes.NewBufferString("title=Test Todo")
		req, err := http.NewRequest("POST", "/todos", formData)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(enableCORS(handleTodos))
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		var response map[string]interface{}
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}

		if title, ok := response["title"].(string); !ok || title != "Test Todo" {
			t.Errorf("handler returned wrong title: got %v want %v",
				title, "Test Todo")
		}
	})

	// TODOの取得テスト
	t.Run("Get Todos", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/todos", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(enableCORS(handleTodos))
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		var response []map[string]interface{}
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}

		if len(response) == 0 {
			t.Error("handler returned empty todo list")
		}
	})

	// TODOの完了状態の切り替えテスト
	t.Run("Toggle Todo", func(t *testing.T) {
		formData := bytes.NewBufferString("id=1")
		req, err := http.NewRequest("PUT", "/todos", formData)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(enableCORS(handleTodos))
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		var response map[string]bool
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}

		if _, exists := response["completed"]; !exists {
			t.Error("handler response missing completed field")
		}
	})

	// TODOの削除テスト
	t.Run("Delete Todo", func(t *testing.T) {
		// まずTODOを追加
		formData := bytes.NewBufferString("title=Test Todo for Delete")
		req, err := http.NewRequest("POST", "/todos", formData)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(enableCORS(handleTodos))
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code for POST: got %v want %v",
				status, http.StatusOK)
		}

		// 追加したTODOのIDを取得
		var addedTodo map[string]interface{}
		if err := json.NewDecoder(rr.Body).Decode(&addedTodo); err != nil {
			t.Fatal(err)
		}

		// 追加したTODOを削除
		formData = bytes.NewBufferString(fmt.Sprintf("id=%v", addedTodo["id"]))
		req, err = http.NewRequest("DELETE", "/todos", formData)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		rr = httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code for DELETE: got %v want %v",
				status, http.StatusOK)
		}

		// 削除が成功したことを確認
		var response map[string]string
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}

		if status, ok := response["status"]; !ok || status != "success" {
			t.Errorf("handler returned wrong response: got %v want %v",
				response, map[string]string{"status": "success"})
		}
	})
} 