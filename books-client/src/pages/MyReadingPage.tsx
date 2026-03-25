import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { api } from '@/api/client'
import { ReadingProgressCard } from '@/components/ReadingProgressCard'
import type { ReadingStatus } from '@/types'

export function MyReadingPage() {
  const [filterStatus, setFilterStatus] = useState<ReadingStatus | 'all'>('all')

  const { data: readingList, isLoading, error } = useQuery({
    queryKey: ['my-reading'],
    queryFn: () => api.reading.getUserBooks(),
  })

  const { data: allBooks } = useQuery({
    queryKey: ['books'],
    queryFn: () => api.books.list(100, 0),
  })

  const getBookForProgress = (bookId: string) => {
    return allBooks?.find(b => b.id === bookId)
  }

  const filteredList = readingList?.filter(p => 
    filterStatus === 'all' || p.status === filterStatus
  )

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
        <p className="text-red">Failed to load reading list</p>
      </div>
    )
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold text-fg-main">My Reading</h1>
        <select
          value={filterStatus}
          onChange={e => setFilterStatus(e.target.value as ReadingStatus | 'all')}
          className="px-4 py-2 border border-border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue bg-bg-dim text-fg-main"
        >
          <option value="all">All Books</option>
          <option value="reading">Currently Reading</option>
          <option value="paused">Paused</option>
          <option value="finished">Finished</option>
        </select>
      </div>
      
      {!filteredList?.length ? (
        <div className="text-center py-12 bg-bg-dim rounded-lg border border-border">
          <svg className="mx-auto h-12 w-12 text-fg-dim" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
          </svg>
          <h3 className="mt-2 text-sm font-semibold text-fg-main">No books in your list</h3>
          <p className="mt-1 text-sm text-fg-dim">Start tracking a book from the book list page.</p>
        </div>
      ) : (
        <div className="space-y-4">
          {filteredList.map(progress => {
            const book = getBookForProgress(progress.book_id)
            return (
              <ReadingProgressCard
                key={progress.id}
                progress={progress}
                bookTitle={book?.title || 'Unknown Book'}
                bookImage={book?.image}
              />
            )
          })}
        </div>
      )}
      
      <div className="mt-8 p-4 bg-bg-dim rounded-lg border border-border">
        <h3 className="font-medium text-fg-main mb-2">Quick Stats</h3>
        <div className="grid grid-cols-3 gap-4 text-center">
          <div>
            <p className="text-2xl font-bold text-blue">
              {readingList?.filter(p => p.status === 'reading').length || 0}
            </p>
            <p className="text-sm text-fg-dim">Reading</p>
          </div>
          <div>
            <p className="text-2xl font-bold text-green">
              {readingList?.filter(p => p.status === 'finished').length || 0}
            </p>
            <p className="text-sm text-fg-dim">Finished</p>
          </div>
          <div>
            <p className="text-2xl font-bold text-yellow">
              {readingList?.filter(p => p.status === 'paused').length || 0}
            </p>
            <p className="text-sm text-fg-dim">Paused</p>
          </div>
        </div>
      </div>
    </div>
  )
}
