import '@testing-library/jest-dom'
import { afterEach } from 'vitest'
import { cleanup } from '@testing-library/react'

// 各テスト後にクリーンアップを実行
afterEach(() => {
  cleanup()
}) 