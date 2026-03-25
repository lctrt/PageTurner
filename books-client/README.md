# Books Client

React frontend for the Books API.

## Tech Stack

- Vite + React + TypeScript
- React Router for navigation
- React Query for data fetching
- TailwindCSS for styling

## Setup

```bash
# Install dependencies
npm install

# Run development server
npm run dev

# Build for production
npm run build
```

## Configuration

The API base URL is configured via Vite proxy in `vite.config.ts`. 
Requests to `/api/*` are proxied to `http://localhost:8080`.

Update the proxy target in `vite.config.ts` to point to your API server.

## Pages

- `/` - Book library grid
- `/books/:id` - Book details with reading progress
- `/import` - Import book from Goodreads URL
- `/login` - User login
- `/register` - User registration
- `/my-reading` - User's reading list (requires auth)
