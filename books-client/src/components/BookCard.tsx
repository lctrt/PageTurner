import { Link } from 'react-router-dom'
import type { Book } from '@/types'
import { AuthorsList } from './AuthorsList'

interface BookCardProps {
  book: Book
}

export function BookCard({ book }: BookCardProps) {
  return (
    <Link
      to={`/books/${book.id}`}
      className="group block bg-bg-dim rounded-lg shadow-sm hover:shadow-lg transition-shadow overflow-hidden border border-border"
    >
      <div className="aspect-[2/3] bg-bg-inactive overflow-hidden">
        {book.image ? (
          <img
            src={book.image}
            alt={book.title}
            className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-200"
          />
        ) : (
          <div className="w-full h-full flex items-center justify-center text-fg-dim">
            <svg className="w-16 h-16" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
            </svg>
          </div>
        )}
      </div>
      <div className="p-4">
        <h3 className="font-semibold text-fg-main truncate">{book.title}</h3>
        <AuthorsList authors={book.authors || []} className="text-sm text-fg-dim truncate block mt-1" />
      </div>
    </Link>
  )
}
