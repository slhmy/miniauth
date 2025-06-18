import React, { useState } from 'react'
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
import { AuthApi, HandlersChangePasswordRequest } from '@/api'

interface ChangePasswordDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  onSuccess?: () => void
}

export const ChangePasswordDialog: React.FC<ChangePasswordDialogProps> = ({
  open,
  onOpenChange,
  onSuccess
}) => {
  const { t } = useTranslation()
  const [formData, setFormData] = useState({
    currentPassword: '',
    newPassword: '',
    confirmPassword: ''
  })
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)

    // Validate passwords match
    if (formData.newPassword !== formData.confirmPassword) {
      setError(t('forms.passwordsDoNotMatch'))
      return
    }

    // Validate password length
    if (formData.newPassword.length < 6) {
      setError(t('forms.passwordTooShort'))
      return
    }

    setLoading(true)
    try {
      const authApi = new AuthApi(undefined, '/api')
      
      const request: HandlersChangePasswordRequest = {
        currentPassword: formData.currentPassword,
        newPassword: formData.newPassword
      }

      await authApi.meChangePasswordPut(request)
      
      // Reset form
      setFormData({
        currentPassword: '',
        newPassword: '',
        confirmPassword: ''
      })
      
      onOpenChange(false)
      onSuccess?.()
      
      // Show success toast
      toast.success(t('profile.passwordChanged'))
      
    } catch (error: unknown) {
      console.error('Change password error:', error)
      
      if (error && typeof error === 'object' && 'response' in error) {
        const axiosError = error as { response?: { status?: number } }
        if (axiosError.response?.status === 400) {
          setError(t('profile.incorrectCurrentPassword'))
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
    setFormData({
      currentPassword: '',
      newPassword: '',
      confirmPassword: ''
    })
    setError(null)
    onOpenChange(false)
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>{t('profile.changePasswordTitle')}</DialogTitle>
          <DialogDescription>
            {t('profile.changePasswordDesc')}
          </DialogDescription>
        </DialogHeader>
        
        <form onSubmit={handleSubmit} className="space-y-4">
          {error && (
            <Alert variant="destructive">
              {error}
            </Alert>
          )}
          
          <div className="space-y-2">
            <label htmlFor="currentPassword" className="text-sm font-medium">
              {t('profile.currentPassword')}
            </label>
            <Input
              id="currentPassword"
              type="password"
              value={formData.currentPassword}
              onChange={(e) => setFormData(prev => ({ 
                ...prev, 
                currentPassword: e.target.value 
              }))}
              required
              disabled={loading}
            />
          </div>
          
          <div className="space-y-2">
            <label htmlFor="newPassword" className="text-sm font-medium">
              {t('profile.newPassword')}
            </label>
            <Input
              id="newPassword"
              type="password"
              value={formData.newPassword}
              onChange={(e) => setFormData(prev => ({ 
                ...prev, 
                newPassword: e.target.value 
              }))}
              required
              disabled={loading}
              minLength={6}
            />
          </div>
          
          <div className="space-y-2">
            <label htmlFor="confirmPassword" className="text-sm font-medium">
              {t('profile.confirmNewPassword')}
            </label>
            <Input
              id="confirmPassword"
              type="password"
              value={formData.confirmPassword}
              onChange={(e) => setFormData(prev => ({ 
                ...prev, 
                confirmPassword: e.target.value 
              }))}
              required
              disabled={loading}
              minLength={6}
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
