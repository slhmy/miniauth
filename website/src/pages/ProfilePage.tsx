import React, { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useAuth } from '@/hooks/useAuth'
import { useDocumentTitle } from '@/hooks/useDocumentTitle'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { User, Mail, Building, Shield, RefreshCw, Edit } from 'lucide-react'
import { ChangePasswordDialog } from '@/components/profile/ChangePasswordDialog'
import { UpdateProfileDialog } from '@/components/profile/UpdateProfileDialog'

export const ProfilePage: React.FC = () => {
  const { user, refreshUser } = useAuth()
  const [refreshing, setRefreshing] = useState(false)
  const [changePasswordOpen, setChangePasswordOpen] = useState(false)
  const [updateProfileOpen, setUpdateProfileOpen] = useState(false)
  const { t } = useTranslation()

  // Set document title for profile page
  useDocumentTitle(t('profile.title'))

  const handleRefresh = async () => {
    setRefreshing(true)
    try {
      await refreshUser()
    } finally {
      setRefreshing(false)
    }
  }

  if (!user) {
    return (
      <div className="min-h-full flex items-center justify-center">
        <div className="text-center">
          <h2 className="text-2xl font-bold text-foreground mb-4">
            {t('common.loading')}
          </h2>
        </div>
      </div>
    )
  }

  // Check if user is admin
  const isAdmin = user.role === 'admin'

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-start">
        <div>
          <h1 className="text-3xl font-bold text-foreground">{t('profile.title')}</h1>
          <p className="text-muted-foreground mt-1">
            {t('profile.description')}
          </p>
        </div>
        <Button 
          variant="outline" 
          onClick={handleRefresh}
          disabled={refreshing}
          className="flex items-center gap-2"
        >
          <RefreshCw className={`w-4 h-4 ${refreshing ? 'animate-spin' : ''}`} />
          {t('common.refresh', 'Refresh')}
        </Button>
      </div>

      {refreshing && (
        <div className="p-3 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg">
          <div className="flex items-center space-x-2">
            <RefreshCw className="w-4 h-4 animate-spin text-blue-600 dark:text-blue-400" />
            <span className="text-sm text-blue-600 dark:text-blue-400">{t('common.loading')}</span>
          </div>
        </div>
      )}

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* User Information Card */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center justify-between">
              <div className="flex items-center">
                <User className="w-5 h-5 mr-2" />
                {t('profile.personalInfo')}
              </div>
              <Button variant="ghost" size="sm" onClick={() => setUpdateProfileOpen(true)}>
                <Edit className="w-4 h-4" />
              </Button>
            </CardTitle>
            <CardDescription>
              {t('profile.accountDetails')}
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-6">
            <div className="grid grid-cols-1 gap-4">
              <div className="space-y-2">
                <label className="text-sm font-medium text-muted-foreground">{t('common.username')}</label>
                <div className="flex items-center space-x-3 p-3 border rounded-lg bg-muted/30">
                  <User className="w-4 h-4 text-muted-foreground" />
                  <span className="text-foreground">{user.username || t('profile.noEmail')}</span>
                </div>
              </div>
              
              <div className="space-y-2">
                <label className="text-sm font-medium text-muted-foreground">{t('common.email')}</label>
                <div className="flex items-center space-x-3 p-3 border rounded-lg bg-muted/30">
                  <Mail className="w-4 h-4 text-muted-foreground" />
                  <span className="text-foreground">{user.email || t('profile.noEmail')}</span>
                </div>
              </div>

              <div className="space-y-2">
                <label className="text-sm font-medium text-muted-foreground">{t('profile.userId')}</label>
                <div className="flex items-center space-x-3 p-3 border rounded-lg bg-muted/30">
                  <Shield className="w-4 h-4 text-muted-foreground" />
                  <span className="text-foreground">#{user.id || t('profile.noEmail')}</span>
                </div>
              </div>

              {isAdmin && (
                <div className="space-y-2">
                  <label className="text-sm font-medium text-muted-foreground">{t('profile.role')}</label>
                  <div className="flex items-center space-x-3 p-3 border rounded-lg bg-primary/5">
                    <Shield className="w-4 h-4 text-primary" />
                    <Badge variant="default">{t('common.admin')}</Badge>
                  </div>
                </div>
              )}
            </div>
          </CardContent>
        </Card>

        {/* Organizations Card */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center">
              <Building className="w-5 h-5 mr-2" />
              {t('profile.organizations')} ({user.organizations?.length || 0})
            </CardTitle>
            <CardDescription>
              {t('profile.organizationsDesc')}
            </CardDescription>
          </CardHeader>
          <CardContent>
            {user.organizations && user.organizations.length > 0 ? (
              <div className="space-y-4">
                {user.organizations.map((org, index) => (
                  <div 
                    key={org.id || index} 
                    className="flex items-center justify-between p-4 border rounded-lg hover:bg-muted/30 transition-colors"
                  >
                    <div className="flex items-center space-x-3">
                      <Building className="w-5 h-5 text-muted-foreground" />
                      <div>
                        <h3 className="font-medium text-foreground">{org.name || t('profile.noEmail')}</h3>
                        <p className="text-sm text-muted-foreground">
                          {org.slug ? `/${org.slug}` : t('profile.noEmail')} • ID: {org.id || t('profile.noEmail')}
                        </p>
                      </div>
                    </div>
                    <Badge variant={
                      org.role === 'owner' || org.role === 'admin' ? 'default' : 'secondary'
                    }>
                      {org.role || 'member'}
                    </Badge>
                  </div>
                ))}
              </div>
            ) : (
              <div className="text-center py-8">
                <Building className="w-12 h-12 text-muted-foreground mx-auto mb-4 opacity-50" />
                <p className="text-muted-foreground">
                  {t('profile.noOrganizations')}
                </p>
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Account Actions */}
      <Card>
        <CardHeader>
          <CardTitle>{t('profile.accountActions')}</CardTitle>
          <CardDescription>
            {t('profile.accountActionsDesc')}
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            <Button 
              variant="outline" 
              className="justify-start h-auto p-4"
              onClick={() => setChangePasswordOpen(true)}
            >
              <div className="text-left">
                <div className="font-medium">{t('profile.changePassword')}</div>
                <div className="text-sm text-muted-foreground">{t('profile.changePasswordDesc')}</div>
              </div>
            </Button>
            
            <Button 
              variant="outline" 
              className="justify-start h-auto p-4"
              onClick={() => setUpdateProfileOpen(true)}
            >
              <div className="text-left">
                <div className="font-medium">{t('profile.updateProfile')}</div>
                <div className="text-sm text-muted-foreground">{t('profile.updateProfileDesc')}</div>
              </div>
            </Button>
            
          </div>
        </CardContent>
      </Card>

      {/* API Status Footer */}
      <div className="text-center pt-4">
        <p className="text-xs text-muted-foreground">
          {t('profile.apiFooter', { time: new Date().toLocaleTimeString() })}
        </p>
      </div>

      {/* Change Password Dialog */}
      <ChangePasswordDialog
        open={changePasswordOpen}
        onOpenChange={setChangePasswordOpen}
        onSuccess={() => {
          // Optionally refresh user data or show additional feedback
          console.log('Password changed successfully')
        }}
      />

      {/* Update Profile Dialog */}
      <UpdateProfileDialog
        open={updateProfileOpen}
        onOpenChange={setUpdateProfileOpen}
        currentUser={user ? {
          username: user.username || '',
          email: user.email || ''
        } : undefined}
        onSuccess={(updatedUser) => {
          // Refresh user data to show updated information
          refreshUser()
          console.log('Profile updated successfully', updatedUser)
        }}
      />
    </div>
  )
}

export default ProfilePage
