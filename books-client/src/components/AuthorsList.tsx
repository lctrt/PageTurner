import { useState } from 'react'
import type { Author } from '@/types'

interface AuthorsListProps {
  authors: Author[]
  className?: string
}

export function AuthorsList({ authors, className = '' }: AuthorsListProps) {
  const [expanded, setExpanded] = useState(false)

  if (!authors || authors.length === 0) {
    return <span className={className}>Unknown author</span>
  }

  const displayAuthors = expanded ? authors : authors.slice(0, 2)
  const hasMore = authors.length > 2

  return (
    <span className={className}>
      {displayAuthors.map(a => a.name).join(', ')}
      {hasMore && (
        <button
          onClick={() => setExpanded(!expanded)}
          className="ml-1 text-blue hover:text-blue-faint"
        >
          {expanded ? ' (less)' : ` (+${authors.length - 2} more)`}
        </button>
      )}
    </span>
  )
}
