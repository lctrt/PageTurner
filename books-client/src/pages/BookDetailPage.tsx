import { useState, useEffect } from 'react'
import { useParams, Link, useNavigate } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { api } from '@/api/client'
import { useAuth } from '@/context/AuthContext'
import { AuthorsList } from '@/components/AuthorsList'
import type { Book, UpdateBookRequest, ReadingProgress, ReadingStatus } from '@/types'

const statusOptions: ReadingStatus[] = ['reading', 'paused', 'finished']

function CoverCard({ book }: { book: Book }) {
  return (
    <div className="bg-bg-dim rounded-lg shadow-sm border border-border overflow-hidden">
      <div className="bg-bg-inactive">
        {book.image ? (
          <img src={book.image} alt={book.title} className="w-full h-auto" />
        ) : (
          <div className="aspect-[2/3] flex items-center justify-center text-fg-dim">
            <svg className="w-24 h-24" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1} d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
            </svg>
          </div>
        )}
      </div>
    </div>
  )
}

function BookDetailCard({
  book,
  isEditing,
  editForm,
  setEditForm,
  onStartEdit,
  onCancelEdit,
  onSave,
  isPending,
  readingProgress,
  progressForm,
  setProgressForm,
  onUpdateProgress,
  isUpdatePending,
}: {
  book: Book
  isEditing: boolean
  editForm: { title: string; blurb: string; custom_link: string }
  setEditForm: React.Dispatch<React.SetStateAction<typeof editForm>>
  onStartEdit: () => void
  onCancelEdit: () => void
  onSave: () => void
  isPending: boolean
  readingProgress?: ReadingProgress
  progressForm: { pages: number; pages_read: number; status: ReadingStatus }
  setProgressForm: React.Dispatch<React.SetStateAction<typeof progressForm>>
  onUpdateProgress: () => void
  isUpdatePending: boolean
}) {
  const { isAuthenticated } = useAuth()

  return (
    <div className="bg-bg-dim rounded-lg shadow-sm border border-border overflow-hidden p-6">
      {isEditing ? (
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-fg-main mb-1">Title</label>
            <input
              type="text"
              value={editForm.title}
              onChange={e => setEditForm(f => ({ ...f, title: e.target.value }))}
              className="w-full px-3 py-2 border border-border rounded-md focus:outline-none focus:ring-2 focus:ring-blue bg-bg-inactive text-fg-main"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-fg-main mb-1">Blurb</label>
            <textarea
              value={editForm.blurb}
              onChange={e => setEditForm(f => ({ ...f, blurb: e.target.value }))}
              rows={4}
              className="w-full px-3 py-2 border border-border rounded-md focus:outline-none focus:ring-2 focus:ring-blue bg-bg-inactive text-fg-main"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-fg-main mb-1">Custom Link</label>
            <input
              type="url"
              value={editForm.custom_link}
              onChange={e => setEditForm(f => ({ ...f, custom_link: e.target.value }))}
              className="w-full px-3 py-2 border border-border rounded-md focus:outline-none focus:ring-2 focus:ring-blue bg-bg-inactive text-fg-main"
            />
          </div>
          <div className="flex gap-2">
            <button
              onClick={onSave}
              disabled={isPending}
              className="px-4 py-2 bg-blue text-bg-main rounded-md hover:bg-blue-faint disabled:opacity-50"
            >
              {isPending ? 'Saving...' : 'Save'}
            </button>
            <button
              onClick={onCancelEdit}
              className="px-4 py-2 border border-border rounded-md hover:bg-bg-inactive text-fg-main"
            >
              Cancel
            </button>
          </div>
        </div>
      ) : (
        <div>
          <div className="flex justify-between items-start">
            <div>
              <h1 className="text-2xl font-bold text-fg-main">{book.title}</h1>
              <AuthorsList authors={book.authors || []} className="text-fg-dim mt-1" />
            </div>
            <button
              onClick={onStartEdit}
              className="text-blue hover:text-blue-faint text-sm font-medium"
            >
              Edit
            </button>
          </div>
          
          {book.blurb && (
            <p className="mt-4 text-fg-dim">{book.blurb}</p>
          )}
          
          <div className="mt-4 flex flex-wrap gap-2">
            {book.goodreads_link && (
              <a
                href={book.goodreads_link}
                target="_blank"
                rel="noopener noreferrer"
                className="text-sm text-blue hover:underline"
              >
                View on Goodreads
              </a>
            )}
            {book.custom_link && (
              <a
                href={book.custom_link}
                target="_blank"
                rel="noopener noreferrer"
                className="text-sm text-blue hover:underline"
              >
                Custom Link
              </a>
            )}
          </div>
        </div>
      )}

      <div className="mt-8 pt-6 border-t border-border">
        {isAuthenticated ? (
          <>
            <h2 className="text-lg font-semibold text-fg-main mb-4">Your Reading Progress</h2>
            
            {readingProgress && (
              <div className="mb-4">
                <div className="flex justify-between text-sm text-fg-dim mb-1">
                  <span>Current: {readingProgress.pages_read} / {readingProgress.pages} pages</span>
                  <span>
                    {readingProgress.pages > 0 
                      ? Math.round((readingProgress.pages_read / readingProgress.pages) * 100)
                      : 0}%
                  </span>
                </div>
                <div className="h-2 bg-bg-inactive rounded-full overflow-hidden">
                  <div 
                    className="h-full bg-blue transition-all duration-300"
                    style={{ width: `${readingProgress.pages > 0 ? (readingProgress.pages_read / readingProgress.pages) * 100 : 0}%` }}
                  />
                </div>
              </div>
            )}
            
            <div className="grid grid-cols-3 gap-4">
              <div>
                <label className="block text-sm font-medium text-fg-main mb-1">Total Pages</label>
                <input
                  type="number"
                  value={progressForm.pages}
                  onChange={e => setProgressForm(f => ({ ...f, pages: parseInt(e.target.value) || 0 }))}
                  className="w-full px-3 py-2 border border-border rounded-md focus:outline-none focus:ring-2 focus:ring-blue bg-bg-inactive text-fg-main"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-fg-main mb-1">Pages Read</label>
                <input
                  type="number"
                  value={progressForm.pages_read}
                  onChange={e => setProgressForm(f => ({ ...f, pages_read: parseInt(e.target.value) || 0 }))}
                  className="w-full px-3 py-2 border border-border rounded-md focus:outline-none focus:ring-2 focus:ring-blue bg-bg-inactive text-fg-main"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-fg-main mb-1">Status</label>
                <select
                  value={progressForm.status}
                  onChange={e => setProgressForm(f => ({ ...f, status: e.target.value as ReadingStatus }))}
                  className="w-full px-3 py-2 border border-border rounded-md focus:outline-none focus:ring-2 focus:ring-blue bg-bg-inactive text-fg-main"
                >
                  {statusOptions.map(status => (
                    <option key={status} value={status}>
                      {status.charAt(0).toUpperCase() + status.slice(1)}
                    </option>
                  ))}
                </select>
              </div>
            </div>
            
            <button
              onClick={onUpdateProgress}
              disabled={isUpdatePending}
              className="mt-4 px-4 py-2 bg-blue text-bg-main rounded-md hover:bg-blue-faint disabled:opacity-50"
            >
              {isUpdatePending ? 'Updating...' : 'Update Progress'}
            </button>
          </>
        ) : (
          <>
            <h2 className="text-lg font-semibold text-fg-main mb-4">Track Your Reading</h2>
            <p className="text-fg-dim mb-4">
              Sign in to track your reading progress for this book.
            </p>
            <Link
              to="/login"
              className="inline-block px-4 py-2 bg-blue text-bg-main rounded-md hover:bg-blue-faint"
            >
              Sign in to track
            </Link>
          </>
        )}
      </div>
    </div>
  )
}

export function BookDetailPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  
  const [isEditing, setIsEditing] = useState(false)
  const [editForm, setEditForm] = useState({
    title: '',
    blurb: '',
    custom_link: '',
  })
  
  const [progressForm, setProgressForm] = useState({
    pages: 0,
    pages_read: 0,
    status: 'reading' as ReadingStatus,
  })

  const { data: book, isLoading, error } = useQuery({
    queryKey: ['books', id],
    queryFn: () => api.books.get(id!),
    enabled: !!id,
  })

  const { data: readingProgress } = useQuery({
    queryKey: ['reading', id],
    queryFn: () => api.reading.getStatus(id!),
    enabled: !!id,
  })

  useEffect(() => {
    if (book && !isEditing) {
      setEditForm({
        title: book.title || '',
        blurb: book.blurb || '',
        custom_link: book.custom_link || '',
      })
    }
  }, [book, isEditing])

  useEffect(() => {
    if (readingProgress) {
      setProgressForm({
        pages: readingProgress.pages,
        pages_read: readingProgress.pages_read,
        status: readingProgress.status,
      })
    }
  }, [readingProgress])

  const updateBook = useMutation({
    mutationFn: (data: UpdateBookRequest) => api.books.update(id!, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['books', id] })
      setIsEditing(false)
    },
  })

  const updateReading = useMutation({
    mutationFn: (data: Parameters<typeof api.reading.updateStatus>[1]) => 
      api.reading.updateStatus(id!, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['reading', id] })
      queryClient.invalidateQueries({ queryKey: ['my-reading'] })
    },
  })

  if (isLoading) {
    return (
      <div className="flex justify-center py-12">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue" />
      </div>
    )
  }

  if (error || !book) {
    return (
      <div className="text-center py-12">
        <p className="text-red">Book not found</p>
        <Link to="/" className="text-blue hover:underline mt-2 block">Back to books</Link>
      </div>
    )
  }

  const handleSaveEdit = () => {
    updateBook.mutate({
      title: editForm.title,
      blurb: editForm.blurb,
      custom_link: editForm.custom_link,
    })
  }

  const handleUpdateProgress = () => {
    updateReading.mutate({
      pages: progressForm.pages,
      pages_read: progressForm.pages_read,
      status: progressForm.status,
    })
  }

  return (
    <div>
      <button 
        onClick={() => navigate('/')}
        className="text-sm text-fg-dim hover:text-fg-main mb-4 flex items-center gap-1"
      >
        ← Back to books
      </button>
      
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <CoverCard book={book} />
        <div className="md:col-span-2">
          <BookDetailCard
            book={book}
            isEditing={isEditing}
            editForm={editForm}
            setEditForm={setEditForm}
            onStartEdit={() => setIsEditing(true)}
            onCancelEdit={() => setIsEditing(false)}
            onSave={handleSaveEdit}
            isPending={updateBook.isPending}
            readingProgress={readingProgress}
            progressForm={progressForm}
            setProgressForm={setProgressForm}
            onUpdateProgress={handleUpdateProgress}
            isUpdatePending={updateReading.isPending}
          />
        </div>
      </div>
    </div>
  )
}
