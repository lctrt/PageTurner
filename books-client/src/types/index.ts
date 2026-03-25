export interface Author {
  id: string
  name: string
}

export interface Publisher {
  id: string
  name: string
}

export interface User {
  id: string
  username: string
  email: string
}

export interface Book {
  id: string
  title: string
  authors: Author[]
  blurb: string
  image: string
  goodreads_link: string
  custom_link: string
  create_at: string
  update_at: string
}

export type ReadingStatus = 'reading' | 'finished' | 'paused'

export interface ReadingProgress {
  id: string
  user_id: string
  book_id: string
  pages: number
  pages_read: number
  status: ReadingStatus
  create_at: string
  update_at: string
}

export interface CreateBookRequest {
  title: string
  authors: string[]
  blurb?: string
  image?: string
  goodreads_link?: string
  custom_link?: string
}

export interface ImportGoodreadsRequest {
  url: string
}

export interface UpdateBookRequest {
  title?: string
  authors?: string[]
  blurb?: string
  image?: string
  goodreads_link?: string
  custom_link?: string
}

export interface UpdateReadingStatusRequest {
  pages?: number
  pages_read?: number
  status?: ReadingStatus
}

export interface RegisterRequest {
  username: string
  email: string
  password: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface AuthResponse {
  token: string
  user: User
}
