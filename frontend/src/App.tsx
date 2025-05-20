import { useState, useEffect, FormEvent, ChangeEvent } from 'react'
import './App.css'

interface Todo {
  id: number
  title: string
  completed: boolean
  created_at: string
}

function App() {
  const [todos, setTodos] = useState<Todo[]>([])
  const [newTodo, setNewTodo] = useState('')

  useEffect(() => {
    fetchTodos()
  }, [])

  const fetchTodos = async () => {
    try {
      const response = await fetch('/api/todos')
      const data = await response.json()
      setTodos(data)
    } catch (error) {
      console.error('TODOの取得に失敗しました:', error)
    }
  }

  const addTodo = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    if (!newTodo.trim()) return

    try {
      const response = await fetch('/api/todos', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: `title=${encodeURIComponent(newTodo)}`
      })

      if (response.ok) {
        setNewTodo('')
        fetchTodos()
      }
    } catch (error) {
      console.error('TODOの追加に失敗しました:', error)
    }
  }

  const toggleTodo = async (id: number) => {
    try {
      const response = await fetch('/api/todos', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: `id=${id}`
      })

      if (response.ok) {
        fetchTodos()
      }
    } catch (error) {
      console.error('TODOの更新に失敗しました:', error)
    }
  }

  const deleteTodo = async (id: number) => {
    try {
      const response = await fetch('/api/todos', {
        method: 'DELETE',
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: `id=${id}`
      })

      if (response.ok) {
        fetchTodos()
      }
    } catch (error) {
      console.error('TODOの削除に失敗しました:', error)
    }
  }

  return (
    <div className="container">
      <h1>TODOアプリ</h1>
      
      <form onSubmit={addTodo} className="add-todo">
        <input
          type="text"
          value={newTodo}
          onChange={(e: ChangeEvent<HTMLInputElement>) => setNewTodo(e.target.value)}
          placeholder="新しいTODOを入力"
        />
        <button type="submit">追加</button>
      </form>

      <div className="todo-list">
        {todos.map((todo: Todo) => (
          <div key={todo.id} className="todo-item">
            <input
              type="checkbox"
              checked={todo.completed}
              onChange={() => toggleTodo(todo.id)}
            />
            <span className={`todo-title ${todo.completed ? 'completed' : ''}`}>
              {todo.title}
            </span>
            <span className="todo-date">
              {new Date(todo.created_at).toLocaleString()}
            </span>
            <button onClick={() => deleteTodo(todo.id)}>削除</button>
          </div>
        ))}
      </div>
    </div>
  )
}

export default App 