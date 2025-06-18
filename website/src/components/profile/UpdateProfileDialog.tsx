import React, { useState, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Alert } from '@/components/ui/alert'
import { toast } from '@/components/ui/toast'
import { AuthApi, HandlersUpdateProfileRequest, HandlersGetUserResponse } from '@/api'

interface UpdateProfileDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  onSuccess?: (updatedUser: HandlersGetUserResponse) => void
  currentUser?: {
    username: string
    email: string
  }
}

export const UpdateProfileDialog: React.FC<UpdateProfileDialogProps> = ({
  open,
  onOpenChange,
  onSuccess,
  currentUser
}) => {
  const { t } = useTranslation()
  const [formData, setFormData] = useState({
    username: '',
    email: ''
  })
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  // Initialize form data when dialog opens or currentUser changes
  useEffect(() => {
    if (currentUser) {
      setFormData({
        username: currentUser.username || '',
        email: currentUser.email || ''
      })
    }
  }, [currentUser, open])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)

    // Validate email format
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    if (!emailRegex.test(formData.email)) {
      setError(t('forms.invalidEmail'))
      return
    }

    // Validate username length
    if (formData.username.length < 3) {
      setError(t('forms.required'))
      return
    }

    setLoading(true)
    try {
      const authApi = new AuthApi(undefined, '/api')
      
      const request: HandlersUpdateProfileRequest = {
        username: formData.username,
      }

      const response = await authApi.meProfilePut(request)
      
      onOpenChange(false)
      onSuccess?.(response.data)
      
      // Show success toast
      toast.success(t('profile.profileUpdated'))
      
    } catch (error: unknown) {
      console.error('Update profile error:', error)
      
      if (error && typeof error === 'object' && 'response' in error) {
        const axiosError = error as { response?: { status?: number; data?: { error?: string } } }
        if (axiosError.response?.status === 400) {
          setError(axiosError.response.data?.error || t('errors.somethingWentWrong'))
        } else {
          setError(t('errors.somethingWentWrong'))
        }
      } else {
        setError(t('errors.somethingWentWrong'))
      }
    } finally {
      setLoading(false)
    }
  }

  const handleCancel = () => {
    // Reset to original values
    if (currentUser) {
      setFormData({
        username: currentUser.username || '',
        email: currentUser.email || ''
      })
    }
    setError(null)
    onOpenChange(false)
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>{t('profile.updateProfileTitle')}</DialogTitle>
          <DialogDescription>
            {t('profile.updateProfileDesc')}
          </DialogDescription>
        </DialogHeader>
        
        <form onSubmit={handleSubmit} className="space-y-4">
          {error && (
            <Alert variant="destructive">
              {error}
            </Alert>
          )}
          
          <div className="space-y-2">
            <label htmlFor="username" className="text-sm font-medium">
              {t('common.username')}
            </label>
            <Input
              id="username"
              type="text"
              value={formData.username}
              onChange={(e) => setFormData(prev => ({ 
                ...prev, 
                username: e.target.value 
              }))}
              required
              disabled={loading}
              minLength={3}
              maxLength={50}
              placeholder={t('forms.enterUsername')}
            />
          </div>
          
          <DialogFooter>
            <Button 
              type="button" 
              variant="outline" 
              onClick={handleCancel}
              disabled={loading}
            >
              {t('common.cancel')}
            </Button>
            <Button type="submit" disabled={loading}>
              {loading ? t('common.loading') : t('common.save')}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
