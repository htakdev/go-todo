import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import App from './App'

// モックの実装
const mockFetch = vi.fn()
global.fetch = mockFetch

describe('App', () => {
  beforeEach(() => {
    // 各テスト前にモックをリセット
    mockFetch.mockReset()
  })

  it('初期状態でTODOリストを表示する', async () => {
    // モックのレスポンスを設定
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: async () => [
        { id: 1, title: 'Test Todo 1', completed: false },
        { id: 2, title: 'Test Todo 2', completed: true },
      ],
    })

    render(<App />)

    // ローディング中は「Loading...」が表示される
    expect(screen.getByText('Loading...')).toBeInTheDocument()

    // TODOリストが表示されるのを待つ
    await waitFor(() => {
      expect(screen.getByText('Test Todo 1')).toBeInTheDocument()
      expect(screen.getByText('Test Todo 2')).toBeInTheDocument()
    })
  })

  it('新しいTODOを追加できる', async () => {
    // 初期のTODOリストをモック
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: async () => [],
    })

    // 新しいTODOの追加レスポンスをモック
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: async () => ({ id: 1, title: 'New Todo', completed: false }),
    })

    render(<App />)

    // 入力フィールドに新しいTODOを入力
    const input = screen.getByPlaceholderText('新しいTODOを入力')
    await userEvent.type(input, 'New Todo')

    // 追加ボタンをクリック
    const addButton = screen.getByText('追加')
    await userEvent.click(addButton)

    // 新しいTODOが表示されるのを待つ
    await waitFor(() => {
      expect(screen.getByText('New Todo')).toBeInTheDocument()
    })
  })

  it('TODOの完了状態を切り替えられる', async () => {
    // 初期のTODOリストをモック
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: async () => [{ id: 1, title: 'Test Todo', completed: false }],
    })

    // 完了状態の切り替えレスポンスをモック
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: async () => ({ completed: true }),
    })

    render(<App />)

    // TODOが表示されるのを待つ
    await waitFor(() => {
      expect(screen.getByText('Test Todo')).toBeInTheDocument()
    })

    // 完了ボタンをクリック
    const toggleButton = screen.getByText('完了')
    await userEvent.click(toggleButton)

    // 完了状態が切り替わるのを待つ
    await waitFor(() => {
      expect(screen.getByText('未完了')).toBeInTheDocument()
    })
  })

  it('TODOを削除できる', async () => {
    // 初期のTODOリストをモック
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: async () => [{ id: 1, title: 'Test Todo', completed: false }],
    })

    // 削除レスポンスをモック
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: async () => ({}),
    })

    render(<App />)

    // TODOが表示されるのを待つ
    await waitFor(() => {
      expect(screen.getByText('Test Todo')).toBeInTheDocument()
    })

    // 削除ボタンをクリック
    const deleteButton = screen.getByText('削除')
    await userEvent.click(deleteButton)

    // TODOが削除されるのを待つ
    await waitFor(() => {
      expect(screen.queryByText('Test Todo')).not.toBeInTheDocument()
    })
  })

  it('エラー時にエラーメッセージを表示する', async () => {
    // エラーレスポンスをモック
    mockFetch.mockRejectedValueOnce(new Error('Network error'))

    render(<App />)

    // エラーメッセージが表示されるのを待つ
    await waitFor(() => {
      expect(screen.getByText('Error: Network error')).toBeInTheDocument()
    })
  })
}) 