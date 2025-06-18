import { useState, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Eye, EyeOff, AlertCircle } from 'lucide-react'
import type { HandlersAdminUserInfo, HandlersAdminCreateUserRequest, HandlersAdminUpdateUserRequest, DatabaseUserRole } from '@/api/api'

interface UserFormProps {
  user?: HandlersAdminUserInfo
  isCreateMode: boolean
  onSubmit: (data: HandlersAdminCreateUserRequest | HandlersAdminUpdateUserRequest) => Promise<void>
  onCancel: () => void
  isLoading?: boolean
  error?: string
}

export default function UserForm({ user, isCreateMode, onSubmit, onCancel, isLoading, error }: UserFormProps) {
  const { t } = useTranslation()
  const [formData, setFormData] = useState({
    username: user?.username || '',
    email: user?.email || '',
    password: '',
    confirmPassword: '',
    role: (user?.role as DatabaseUserRole) || 'user' as DatabaseUserRole
  })
  const [showPassword, setShowPassword] = useState(false)
  const [showConfirmPassword, setShowConfirmPassword] = useState(false)
  const [validationErrors, setValidationErrors] = useState<Record<string, string>>({})

  useEffect(() => {
    if (user && !isCreateMode) {
      setFormData(prev => ({
        ...prev,
        username: user.username || '',
        email: user.email || '',
        role: (user.role as DatabaseUserRole) || 'user' as DatabaseUserRole
      }))
    }
  }, [user, isCreateMode])

  const validateForm = () => {
    const errors: Record<string, string> = {}

    // Username validation
    if (!formData.username.trim()) {
      errors.username = t('forms.required')
    }

    // Email validation
    if (!formData.email.trim()) {
      errors.email = t('forms.required')
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email)) {
      errors.email = t('forms.invalidEmail')
    }

    // Password validation (only for create mode or when password is provided)
    if (isCreateMode) {
      if (!formData.password) {
        errors.password = t('forms.required')
      } else if (formData.password.length < 6) {
        errors.password = t('forms.passwordTooShort')
      }

      if (!formData.confirmPassword) {
        errors.confirmPassword = t('forms.required')
      } else if (formData.password !== formData.confirmPassword) {
        errors.confirmPassword = t('forms.passwordsDoNotMatch')
      }
    } else if (formData.password) {
      // In edit mode, only validate password if it's provided
      if (formData.password.length < 6) {
        errors.password = t('forms.passwordTooShort')
      }
      if (formData.password !== formData.confirmPassword) {
        errors.confirmPassword = t('forms.passwordsDoNotMatch')
      }
    }

    setValidationErrors(errors)
    return Object.keys(errors).length === 0
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!validateForm()) {
      return
    }

    try {
      if (isCreateMode) {
        const createData: HandlersAdminCreateUserRequest = {
          username: formData.username.trim(),
          email: formData.email.trim(),
          password: formData.password
        }
        await onSubmit(createData)
      } else {
        const updateData: HandlersAdminUpdateUserRequest = {
          username: formData.username.trim(),
          email: formData.email.trim(),
          role: formData.role
        }
        await onSubmit(updateData)
      }
    } catch {
      // Error handling is done by parent component
    }
  }

  const handleInputChange = (field: string, value: string) => {
    setFormData(prev => ({ ...prev, [field]: value }))
    // Clear validation error when user starts typing
    if (validationErrors[field]) {
      setValidationErrors(prev => {
        // eslint-disable-next-line @typescript-eslint/no-unused-vars
        const { [field]: _, ...rest } = prev
        return rest
      })
    }
  }

  return (
    <div className="space-y-4">
      <form onSubmit={handleSubmit} className="space-y-4">
        {error && (
          <Alert variant="destructive">
            <AlertCircle className="h-4 w-4" />
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        )}

        {/* Username */}
        <div className="space-y-2">
          <label htmlFor="username" className="text-sm font-medium">
            {t('common.username')} *
          </label>
          <Input
            id="username"
            type="text"
            value={formData.username}
            onChange={(e) => handleInputChange('username', e.target.value)}
            placeholder={t('forms.enterUsername')}
            className={validationErrors.username ? 'border-red-500' : ''}
          />
          {validationErrors.username && (
            <p className="text-sm text-red-500">{validationErrors.username}</p>
          )}
        </div>

        {/* Email */}
        <div className="space-y-2">
          <label htmlFor="email" className="text-sm font-medium">
            {t('common.email')} *
          </label>
          <Input
            id="email"
            type="email"
            value={formData.email}
            onChange={(e) => handleInputChange('email', e.target.value)}
            placeholder={t('forms.enterEmail')}
            className={validationErrors.email ? 'border-red-500' : ''}
          />
          {validationErrors.email && (
            <p className="text-sm text-red-500">{validationErrors.email}</p>
          )}
        </div>

        {/* Role (only for edit mode) */}
        {!isCreateMode && (
          <div className="space-y-2">
            <label htmlFor="role" className="text-sm font-medium">
              {t('admin.role')}
            </label>
            <select 
              id="role" 
              value={formData.role} 
              onChange={(e) => handleInputChange('role', e.target.value as DatabaseUserRole)}
              className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
            >
              <option value="user">{t('roles.user')}</option>
              <option value="admin">{t('roles.admin')}</option>
            </select>
          </div>
        )}

        {/* Password */}
        <div className="space-y-2">
          <label htmlFor="password" className="text-sm font-medium">
            {t('common.password')} {isCreateMode ? '*' : '(optional)'}
          </label>
          <div className="relative">
            <Input
              id="password"
              type={showPassword ? 'text' : 'password'}
              value={formData.password}
              onChange={(e) => handleInputChange('password', e.target.value)}
              placeholder={isCreateMode ? t('forms.enterPassword') : t('forms.leaveEmptyKeepCurrent')}
              className={validationErrors.password ? 'border-red-500 pr-10' : 'pr-10'}
            />
            <Button
              type="button"
              variant="ghost"
              size="sm"
              className="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent"
              onClick={() => setShowPassword(!showPassword)}
            >
              {showPassword ? (
                <EyeOff className="h-4 w-4" />
              ) : (
                <Eye className="h-4 w-4" />
              )}
            </Button>
          </div>
          {validationErrors.password && (
            <p className="text-sm text-red-500">{validationErrors.password}</p>
          )}
        </div>

        {/* Confirm Password */}
        {(isCreateMode || formData.password) && (
          <div className="space-y-2">
            <label htmlFor="confirmPassword" className="text-sm font-medium">
              {t('forms.confirmPassword')} *
            </label>
            <div className="relative">
              <Input
                id="confirmPassword"
                type={showConfirmPassword ? 'text' : 'password'}
                value={formData.confirmPassword}
                onChange={(e) => handleInputChange('confirmPassword', e.target.value)}
                placeholder={t('forms.confirmPasswordPlaceholder')}
                className={validationErrors.confirmPassword ? 'border-red-500 pr-10' : 'pr-10'}
              />
              <Button
                type="button"
                variant="ghost"
                size="sm"
                className="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent"
                onClick={() => setShowConfirmPassword(!showConfirmPassword)}
              >
                {showConfirmPassword ? (
                  <EyeOff className="h-4 w-4" />
                ) : (
                  <Eye className="h-4 w-4" />
                )}
              </Button>
            </div>
            {validationErrors.confirmPassword && (
              <p className="text-sm text-red-500">{validationErrors.confirmPassword}</p>
            )}
          </div>
        )}

        {/* Actions */}
        <div className="flex gap-3 pt-4">
          <Button
            type="button"
            variant="outline"
            onClick={onCancel}
            className="flex-1"
            disabled={isLoading}
          >
            {t('common.cancel')}
          </Button>
          <Button
            type="submit"
            className="flex-1"
            disabled={isLoading}
          >
            {isLoading 
              ? t('common.loading') 
              : isCreateMode 
                ? t('common.create') 
                : t('common.update')
            }
          </Button>
        </div>
      </form>
    </div>
  )
}
