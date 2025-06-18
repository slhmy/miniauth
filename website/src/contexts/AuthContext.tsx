import { createContext } from 'react'
import type { HandlersGetUserResponse } from '@/api'

export interface AuthContextType {
  user: HandlersGetUserResponse | null
  login: (email: string, password: string) => Promise<boolean>
  register: (username: string, email: string, password: string) => Promise<{ success: boolean; message?: string }>
  logout: () => Promise<void>
  refreshUser: () => Promise<void>
  loading: boolean
  isAuthenticated: boolean
}

// Create the context
export const AuthContext = createContext<AuthContextType | undefined>(undefined)
