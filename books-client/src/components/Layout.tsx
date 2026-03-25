import { Outlet, Link, useLocation } from 'react-router-dom'
import { useAuth } from '@/context/AuthContext'

export function Layout() {
  const location = useLocation()
  const { user, isAuthenticated, logout } = useAuth()

  const navLinks = [
    { to: '/', label: 'Books' },
    { to: '/import', label: 'Import' },
  ]

  const authLinks = isAuthenticated
    ? { to: '/my-reading', label: 'My Reading' }
    : null

  return (
    <div className="min-h-screen bg-bg-main">
      <header className="bg-bg-dim shadow-sm border-b border-border">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <Link to="/" className="flex items-center gap-2">
              <svg className="w-8 h-8 text-blue" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
              </svg>
              <span className="text-xl font-bold text-fg-main">Books</span>
            </Link>
            <nav className="flex items-center gap-6">
              {navLinks.map(link => (
                <Link
                  key={link.to}
                  to={link.to}
                  className={`text-sm font-medium transition-colors ${
                    location.pathname === link.to
                      ? 'text-blue'
                      : 'text-fg-dim hover:text-fg-main'
                  }`}
                >
                  {link.label}
                </Link>
              ))}
              {authLinks && (
                <Link
                  to={authLinks.to}
                  className={`text-sm font-medium transition-colors ${
                    location.pathname === authLinks.to
                      ? 'text-blue'
                      : 'text-fg-dim hover:text-fg-main'
                  }`}
                >
                  {authLinks.label}
                </Link>
              )}
              {isAuthenticated ? (
                <div className="flex items-center gap-4">
                  <span className="text-sm text-fg-dim">{user?.username}</span>
                  <button
                    onClick={logout}
                    className="text-sm text-fg-dim hover:text-fg-main"
                  >
                    Logout
                  </button>
                </div>
              ) : (
                <div className="flex items-center gap-4">
                  <Link
                    to="/login"
                    className="text-sm font-medium text-fg-dim hover:text-fg-main"
                  >
                    Login
                  </Link>
                  <Link
                    to="/register"
                    className="px-4 py-2 bg-blue text-bg-main text-sm font-medium rounded-lg hover:bg-blue-faint transition-colors"
                  >
                    Register
                  </Link>
                </div>
              )}
            </nav>
          </div>
        </div>
      </header>
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <Outlet />
      </main>
    </div>
  )
}
