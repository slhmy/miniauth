import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { Toaster } from 'sonner'
import { ThemeProvider } from '@/components/theme-provider'
import { AuthProvider } from '@/contexts/AuthProvider'
import { ProtectedRoute } from '@/components/ProtectedRoute'
import { MainLayout } from '@/components/layout/MainLayout'
import { LoginPage } from '@/pages/LoginPage'
import { ProfilePage } from '@/pages/ProfilePage'
import UsersPage from '@/pages/admin/Users'
import OAuthApplications from '@/pages/OAuthApplications'
import OAuthAuthorization from '@/pages/OAuthAuthorization'
import './App.css'

// Create a client for React Query
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 1000 * 60 * 5, // 5 minutes
      retry: 1,
    },
  },
})

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <ThemeProvider defaultTheme="system" storageKey="miniauth-ui-theme">
        <AuthProvider>
          <Router>
            <Routes>
              {/* Public routes */}
              <Route path="/login" element={<LoginPage />} />
              <Route path="/oauth/authorize" element={<OAuthAuthorization />} />
              
              {/* Protected routes with unified layout */}
              <Route path="/*" element={
                <ProtectedRoute>
                  <MainLayout>
                    <Routes>
                      {/* User routes */}
                      <Route path="/profile" element={<ProfilePage />} />
                      
                      {/* Admin routes - require admin privileges */}
                      <Route path="/oauth/applications" element={
                        <ProtectedRoute requireAdmin={true}>
                          <OAuthApplications />
                        </ProtectedRoute>
                      } />
                      <Route path="/admin/users" element={
                        <ProtectedRoute requireAdmin={true}>
                          <UsersPage />
                        </ProtectedRoute>
                      } />
                      
                      {/* Default redirects */}
                      <Route path="/" element={<Navigate to="/profile" replace />} />
                      <Route path="*" element={<Navigate to="/profile" replace />} />
                    </Routes>
                  </MainLayout>
                </ProtectedRoute>
              } />
            </Routes>
          </Router>
        </AuthProvider>
        <Toaster position="top-right" richColors />
      </ThemeProvider>
    </QueryClientProvider>
  )
}

export default App
