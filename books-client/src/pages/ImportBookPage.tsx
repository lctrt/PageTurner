import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { api } from '@/api/client'

export function ImportBookPage() {
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const [url, setUrl] = useState('')
  const [error, setError] = useState('')

  const importMutation = useMutation({
    mutationFn: (url: string) => api.books.importGoodreads({ url }),
    onSuccess: (book) => {
      queryClient.invalidateQueries({ queryKey: ['books'] })
      navigate(`/books/${book.id}`)
    },
    onError: (err: Error) => {
      setError(err.message)
    },
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    
    if (!url.trim()) {
      setError('Please enter a Goodreads URL')
      return
    }
    
    if (!url.includes('goodreads.com')) {
      setError('Please enter a valid Goodreads URL')
      return
    }
    
    importMutation.mutate(url)
  }

  return (
    <div className="max-w-2xl mx-auto">
      <h1 className="text-2xl font-bold text-fg-main mb-6">Import from Goodreads</h1>
      
      <div className="bg-bg-dim rounded-lg shadow-sm border border-border p-6">
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label htmlFor="url" className="block text-sm font-medium text-fg-main mb-1">
              Goodreads URL
            </label>
            <input
              type="url"
              id="url"
              value={url}
              onChange={e => setUrl(e.target.value)}
              placeholder="https://www.goodreads.com/book/show/12345.Book-Title"
              className="w-full px-4 py-3 border border-border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue bg-bg-inactive text-fg-main"
            />
            <p className="mt-2 text-sm text-fg-dim">
              Paste a link to a book page on Goodreads. We'll extract the title, author, cover image, and description.
            </p>
          </div>
          
          {error && (
            <div className="p-3 bg-bg-inactive border border-red rounded-lg">
              <p className="text-sm text-red">{error}</p>
            </div>
          )}
          
          <button
            type="submit"
            disabled={importMutation.isPending}
            className="w-full px-4 py-3 bg-blue text-bg-main font-medium rounded-lg hover:bg-blue-faint disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            {importMutation.isPending ? (
              <span className="flex items-center justify-center gap-2">
                <svg className="animate-spin h-5 w-5" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none" />
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                </svg>
                Importing...
              </span>
            ) : (
              'Import Book'
            )}
          </button>
        </form>
      </div>
      
      <div className="mt-6 p-4 bg-bg-dim rounded-lg">
        <h3 className="font-medium text-fg-main mb-2">Example URLs:</h3>
        <ul className="text-sm text-fg-dim space-y-1">
          <li>• https://www.goodreads.com/book/show/3.Harry_Potter_and_the_Sorcerer_s_Stone</li>
          <li>• https://www.goodreads.com/book/show/4099.The_Pragmatic_Programmer</li>
          <li>• https://www.goodreads.com/book/show/1.The_Hobbit</li>
        </ul>
      </div>
    </div>
  )
}
