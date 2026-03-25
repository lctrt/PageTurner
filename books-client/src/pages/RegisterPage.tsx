import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { authApi } from '@/api/auth'
import { useAuth } from '@/context/AuthContext'

export function RegisterPage() {
  const navigate = useNavigate()
  const { login } = useAuth()
  const [form, setForm] = useState({ username: '', email: '', password: '', confirmPassword: '' })
  const [error, setError] = useState('')
  const [isLoading, setIsLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')

    if (form.password !== form.confirmPassword) {
      setError('Passwords do not match')
      return
    }

    if (form.password.length < 6) {
      setError('Password must be at least 6 characters')
      return
    }

    setIsLoading(true)

    try {
      const response = await authApi.register({
        username: form.username,
        email: form.email,
        password: form.password,
      })
      login(response.token, response.user)
      navigate('/')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Registration failed')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="min-h-[80vh] flex items-center justify-center">
      <div className="max-w-md w-full bg-bg-dim rounded-lg shadow-sm border border-border p-8">
        <div className="text-center mb-8">
          <h1 className="text-2xl font-bold text-fg-main">Create Account</h1>
          <p className="text-fg-dim mt-2">Start tracking your reading</p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label htmlFor="username" className="block text-sm font-medium text-fg-main mb-1">
              Username
            </label>
            <input
              type="text"
              id="username"
              value={form.username}
              onChange={e => setForm(f => ({ ...f, username: e.target.value }))}
              required
              className="w-full px-4 py-2 border border-border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue bg-bg-inactive text-fg-main"
            />
          </div>

          <div>
            <label htmlFor="email" className="block text-sm font-medium text-fg-main mb-1">
              Email
            </label>
            <input
              type="email"
              id="email"
              value={form.email}
              onChange={e => setForm(f => ({ ...f, email: e.target.value }))}
              required
              className="w-full px-4 py-2 border border-border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue bg-bg-inactive text-fg-main"
            />
          </div>

          <div>
            <label htmlFor="password" className="block text-sm font-medium text-fg-main mb-1">
              Password
            </label>
            <input
              type="password"
              id="password"
              value={form.password}
              onChange={e => setForm(f => ({ ...f, password: e.target.value }))}
              required
              className="w-full px-4 py-2 border border-border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue bg-bg-inactive text-fg-main"
            />
          </div>

          <div>
            <label htmlFor="confirmPassword" className="block text-sm font-medium text-fg-main mb-1">
              Confirm Password
            </label>
            <input
              type="password"
              id="confirmPassword"
              value={form.confirmPassword}
              onChange={e => setForm(f => ({ ...f, confirmPassword: e.target.value }))}
              required
              className="w-full px-4 py-2 border border-border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue bg-bg-inactive text-fg-main"
            />
          </div>

          {error && (
            <div className="p-3 bg-bg-inactive border border-red rounded-lg">
              <p className="text-sm text-red">{error}</p>
            </div>
          )}

          <button
            type="submit"
            disabled={isLoading}
            className="w-full py-3 bg-blue text-bg-main font-medium rounded-lg hover:bg-blue-faint disabled:opacity-50 transition-colors"
          >
            {isLoading ? 'Creating account...' : 'Create Account'}
          </button>
        </form>

        <p className="text-center text-sm text-fg-dim mt-6">
          Already have an account?{' '}
          <Link to="/login" className="text-blue hover:underline">
            Sign in
          </Link>
        </p>
      </div>
    </div>
  )
}
