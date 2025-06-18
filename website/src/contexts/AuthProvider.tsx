import React, { useState, useEffect } from 'react'
import type { ReactNode } from 'react'
import type { HandlersGetUserResponse } from '@/api'
import { AuthApi, UsersApi } from '@/api'
import type { AxiosResponse } from 'axios'
import { AuthContext, type AuthContextType } from './AuthContext'

interface AuthProviderProps {
  children: ReactNode
}

const authApi = new AuthApi(undefined, '/api')
const usersApi = new UsersApi(undefined, '/api')

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [user, setUser] = useState<HandlersGetUserResponse | null>(null)
  const [loading, setLoading] = useState(true)

  // Check if user is authenticated on app start
  useEffect(() => {
    refreshUser()
  }, [])

  const refreshUser = async () => {
    try {
      setLoading(true)
      const response: AxiosResponse<HandlersGetUserResponse> = await authApi.meGet()
      setUser(response.data)
    } catch {
      console.log('User not authenticated')
      setUser(null)
    } finally {
      setLoading(false)
    }
  }

  const login = async (email: string, password: string): Promise<boolean> => {
    try {
      setLoading(true)
      await authApi.loginPost({ email, password })
      
      // After successful login, get user info
      const response: AxiosResponse<HandlersGetUserResponse> = await authApi.meGet()
      setUser(response.data)
      
      return true
    } catch (error) {
      console.error('Login failed:', error)
      setUser(null)
      return false
    } finally {
      setLoading(false)
    }
  }

  const register = async (username: string, email: string, password: string): Promise<{ success: boolean; message?: string }> => {
    try {
      setLoading(true)
      const response = await usersApi.usersPost({ username, email, password })
      
      // Registration successful, now try to login automatically
      const loginSuccess = await login(email, password)
      
      return { 
        success: loginSuccess, 
        message: response.data.message || 'Registration successful' 
      }
    } catch (error) {
      console.error('Registration failed:', error)
      let message = 'Registration failed'
      
      // Extract error message from response if available
      if (error && typeof error === 'object' && 'response' in error) {
        const axiosError = error as { response?: { data?: { message?: string; error?: string } } }
        if (axiosError.response?.data?.message) {
          message = axiosError.response.data.message
        } else if (axiosError.response?.data?.error) {
          message = axiosError.response.data.error
        }
      } else if (error && typeof error === 'object' && 'message' in error) {
        message = (error as { message: string }).message
      }
      
      return { success: false, message }
    } finally {
      setLoading(false)
    }
  }

  const logout = async () => {
    try {
      setLoading(true)
      await authApi.logoutPost()
    } catch (error) {
      console.error('Logout error:', error)
    } finally {
      setUser(null)
      setLoading(false)
    }
  }

  const value: AuthContextType = {
    user,
    login,
    register,
    logout,
    refreshUser,
    loading,
    isAuthenticated: !!user,
  }

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}
