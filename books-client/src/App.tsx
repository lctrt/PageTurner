import { Routes, Route } from 'react-router-dom'
import { Layout } from '@/components/Layout'
import { ProtectedRoute } from '@/components/ProtectedRoute'
import { BookListPage } from '@/pages/BookListPage'
import { BookDetailPage } from '@/pages/BookDetailPage'
import { ImportBookPage } from '@/pages/ImportBookPage'
import { MyReadingPage } from '@/pages/MyReadingPage'
import { LoginPage } from '@/pages/LoginPage'
import { RegisterPage } from '@/pages/RegisterPage'

function App() {
  return (
    <Routes>
      <Route path="/" element={<Layout />}>
        <Route index element={<BookListPage />} />
        <Route path="books/:id" element={<BookDetailPage />} />
        <Route path="import" element={<ImportBookPage />} />
        <Route path="login" element={<LoginPage />} />
        <Route path="register" element={<RegisterPage />} />
        <Route
          path="my-reading"
          element={
            <ProtectedRoute>
              <MyReadingPage />
            </ProtectedRoute>
          }
        />
      </Route>
    </Routes>
  )
}

export default App
