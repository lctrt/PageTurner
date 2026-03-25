import type {
  Book,
  CreateBookRequest,
  ImportGoodreadsRequest,
  UpdateBookRequest,
  ReadingProgress,
  UpdateReadingStatusRequest,
} from '@/types'

const API_BASE = '/api'
const TOKEN_KEY = 'books_auth_token'

function getToken(): string | null {
  return localStorage.getItem(TOKEN_KEY)
}

function authHeaders(): HeadersInit {
  const token = getToken()
  const headers: HeadersInit = { 'Content-Type': 'application/json' }
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }
  return headers
}

async function handleResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    if (response.status === 401) {
      localStorage.removeItem(TOKEN_KEY)
      window.location.href = '/login'
    }
    const error = await response.text()
    throw new Error(error || `HTTP ${response.status}`)
  }
  return response.json()
}

export const api = {
  books: {
    list: async (limit = 20, offset = 0): Promise<Book[]> => {
      const res = await fetch(`${API_BASE}/books?limit=${limit}&offset=${offset}`)
      return handleResponse(res)
    },

    get: async (id: string): Promise<Book> => {
      const res = await fetch(`${API_BASE}/books/${id}`)
      return handleResponse(res)
    },

    create: async (data: CreateBookRequest): Promise<Book> => {
      const res = await fetch(`${API_BASE}/books`, {
        method: 'POST',
        headers: authHeaders(),
        body: JSON.stringify(data),
      })
      return handleResponse(res)
    },

    update: async (id: string, data: UpdateBookRequest): Promise<Book> => {
      const res = await fetch(`${API_BASE}/books/${id}`, {
        method: 'PUT',
        headers: authHeaders(),
        body: JSON.stringify(data),
      })
      return handleResponse(res)
    },

    importGoodreads: async (data: ImportGoodreadsRequest): Promise<Book> => {
      const res = await fetch(`${API_BASE}/books/import`, {
        method: 'POST',
        headers: authHeaders(),
        body: JSON.stringify(data),
      })
      return handleResponse(res)
    },
  },

  reading: {
    getUserBooks: async (): Promise<ReadingProgress[]> => {
      const res = await fetch(`${API_BASE}/me/books`, {
        headers: authHeaders(),
      })
      return handleResponse(res)
    },

    getStatus: async (bookId: string): Promise<ReadingProgress> => {
      const res = await fetch(`${API_BASE}/me/books/${bookId}/status`, {
        headers: authHeaders(),
      })
      return handleResponse(res)
    },

    updateStatus: async (
      bookId: string,
      data: UpdateReadingStatusRequest
    ): Promise<ReadingProgress> => {
      const res = await fetch(`${API_BASE}/me/books/${bookId}/status`, {
        method: 'PUT',
        headers: authHeaders(),
        body: JSON.stringify(data),
      })
      return handleResponse(res)
    },
  },
}
