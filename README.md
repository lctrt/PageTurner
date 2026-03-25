# PageTurner

PageTurner is a tool to track what you're reading.

## Features:
* import a book from a goodreads link
* update reading status (page read, finished)
* edit book details

## Tech Stack
- Golang API
- PostgreSQL
- Vite + React + TypeScript
- React Router for navigation
- React Query for data fetching
- TailwindCSS for styling
- Docker Compose for Container setup

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/health` | Health check |
| `POST` | `/auth/register` | Register new user |
| `POST` | `/auth/login` | Login and get JWT token |
| `POST` | `/books` | Create a book |
| `POST` | `/books/import` | Import book from Goodreads link |
| `GET` | `/books` | List books (query: `?limit=20&offset=0`) |
| `GET` | `/books/{id}` | Get book by ID |
| `PUT` | `/books/{id}` | Update book |
| `GET` | `/me/books` | Get authenticated user's reading list (auth required) |
| `GET` | `/me/books/{bookId}/status` | Get reading status (auth required) |
| `PUT` | `/me/books/{bookId}/status` | Update reading status (auth required) |

### Create Book Request
```json
{
  "title": "The Pragmatic Programmer",
  "authors": ["David Thomas", "Andrew Hunt"],
  "blurb": "Your journey to mastery",
  "goodreads_link": "https://goodreads.com/book/123"
}
```

### Import from Goodreads Request
```json
{
  "url": "https://www.goodreads.com/book/show/12345.Title"
}
```

### Update Reading Status Request
```json
{
  "pages": 400,
  "pages_read": 150,
  "status": "reading"
}
```

Valid status values: `reading`, `finished`, `paused`

### Register Request
```json
{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "secret123"
}
```

### Login Request
```json
{
  "username": "johndoe",
  "password": "secret123"
}
```

### Auth Response
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "uuid",
    "username": "johndoe",
    "email": "john@example.com"
  }
}
```

Protected endpoints require `Authorization: Bearer <token>` header.

## Getting Started

```bash
docker-compose up --build
```
