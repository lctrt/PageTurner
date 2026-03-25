import { Link } from 'react-router-dom'
import type { ReadingProgress } from '@/types'

interface ReadingProgressCardProps {
  progress: ReadingProgress
  bookTitle: string
  bookImage?: string
}

const statusColors = {
  reading: 'bg-bg-inactive text-blue',
  finished: 'bg-bg-inactive text-green',
  paused: 'bg-bg-inactive text-yellow',
}

const statusLabels = {
  reading: 'Reading',
  finished: 'Finished',
  paused: 'Paused',
}

export function ReadingProgressCard({ progress, bookTitle, bookImage }: ReadingProgressCardProps) {
  const progressPercent = progress.pages > 0 
    ? Math.round((progress.pages_read / progress.pages) * 100) 
    : 0

  return (
    <Link to={`/books/${progress.book_id}`} className="block bg-bg-dim rounded-lg shadow-sm border border-border p-4 hover:shadow-lg transition-shadow">
      <div className="flex gap-4">
        <div className="w-16 h-24 bg-bg-inactive rounded overflow-hidden flex-shrink-0">
          {bookImage ? (
            <img src={bookImage} alt={bookTitle} className="w-full h-full object-cover" />
          ) : (
            <div className="w-full h-full flex items-center justify-center text-fg-dim">
              <svg className="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
              </svg>
            </div>
          )}
        </div>
        <div className="flex-1 min-w-0">
          <h3 className="font-semibold text-fg-main truncate">{bookTitle}</h3>
          <span className={`inline-block px-2 py-0.5 text-xs font-medium rounded-full mt-1 ${statusColors[progress.status]}`}>
            {statusLabels[progress.status]}
          </span>
          <div className="mt-3">
            <div className="flex justify-between text-sm text-fg-dim mb-1">
              <span>Progress</span>
              <span>{progress.pages_read} / {progress.pages} pages ({progressPercent}%)</span>
            </div>
            <div className="h-2 bg-bg-inactive rounded-full overflow-hidden">
              <div 
                className="h-full bg-blue transition-all duration-300"
                style={{ width: `${progressPercent}%` }}
              />
            </div>
          </div>
        </div>
      </div>
    </Link>
  )
}
