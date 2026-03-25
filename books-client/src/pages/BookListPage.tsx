import { useQuery } from '@tanstack/react-query'
import { api } from '@/api/client'
import { BookCard } from '@/components/BookCard'

export function BookListPage() {
  const { data: books, isLoading, error } = useQuery({
    queryKey: ['books'],
    queryFn: () => api.books.list(),
  })

  if (isLoading) {
    return (
      <div className="flex justify-center py-12">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue" />
      </div>
    )
  }

  if (error) {
    return (
      <div className="text-center py-12">
        <p className="text-red">Failed to load books</p>
      </div>
    )
  }

  if (!books?.length) {
    return (
      <div className="text-center py-12">
        <svg className="mx-auto h-12 w-12 text-fg-dim" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
        </svg>
        <h3 className="mt-2 text-sm font-semibold text-fg-main">No books yet</h3>
        <p className="mt-1 text-sm text-fg-dim">Import a book from Goodreads to get started.</p>
      </div>
    )
  }

  return (
    <div>
      <h1 className="text-2xl font-bold text-fg-main mb-6">All Books</h1>
      <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-6">
        {books.map(book => (
          <BookCard key={book.id} book={book} />
        ))}
      </div>
    </div>
  )
}
