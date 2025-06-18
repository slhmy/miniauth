import React, { useState } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { useAuth } from '@/hooks/useAuth'
import { useDocumentTitle } from '@/hooks/useDocumentTitle'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { LanguageToggle } from '@/components/language-toggle'
import { ThemeToggle } from '@/components/theme-toggle'
import { Loader2, AlertCircle, Shield } from 'lucide-react'

export const LoginPage: React.FC = () => {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [username, setUsername] = useState('')
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')
  const [isSignUp, setIsSignUp] = useState(false)
  const { login, register, loading } = useAuth()
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const { t } = useTranslation()

  // Set document title based on login/signup mode
  useDocumentTitle(isSignUp ? t('auth.signUp') : t('auth.signIn'))

  // Check if this is an OAuth redirect
  const isOAuthRedirect = searchParams.get('oauth_redirect') === 'true'
  
  // Store OAuth parameters for later use
  const oauthParams = {
    client_id: searchParams.get('client_id'),
    redirect_uri: searchParams.get('redirect_uri'),
    scope: searchParams.get('scope'),
    state: searchParams.get('state'),
    response_type: searchParams.get('response_type'),
    code_challenge: searchParams.get('code_challenge'),
    code_challenge_method: searchParams.get('code_challenge_method'),
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setSuccess('')

    if (!email || !password || (isSignUp && !username)) {
      setError(t('auth.fillAllFields'))
      return
    }

    if (isSignUp && password !== confirmPassword) {
      setError(t('auth.passwordMismatch'))
      return
    }

    if (isSignUp) {
      // Handle registration
      const result = await register(username, email, password)
      if (result.success) {
        setSuccess(t('auth.registrationSuccess'))
        handlePostLoginRedirect()
      } else {
        setError(result.message || t('auth.registrationFailed'))
      }
    } else {
      // Handle login
      const success = await login(email, password)
      if (success) {
        handlePostLoginRedirect()
      } else {
        setError(t('auth.invalidCredentials'))
      }
    }
  }

  const handlePostLoginRedirect = () => {
    if (isOAuthRedirect && oauthParams.client_id) {
      // Build OAuth authorization URL with preserved parameters
      const params = new URLSearchParams()
      Object.entries(oauthParams).forEach(([key, value]) => {
        if (value) {
          params.set(key, value)
        }
      })
      navigate(`/oauth/authorize?${params.toString()}`)
    } else {
      navigate('/profile')
    }
  }

  const toggleMode = () => {
    setIsSignUp(!isSignUp)
    setError('')
    setSuccess('')
    setUsername('')
    setEmail('')
    setPassword('')
    setConfirmPassword('')
  }

  return (
    <div className="min-h-screen flex flex-col bg-background">
      {/* Header */}
      <header className="flex h-16 items-center justify-between border-b bg-background px-6 shadow-sm flex-shrink-0 w-full">
        <div className="flex items-center space-x-2">
          <Shield className="h-6 w-6 text-primary" />
          <span className="text-lg font-semibold text-foreground">MiniAuth</span>
        </div>
        
        <div className="flex items-center space-x-4">
          <ThemeToggle />
          <LanguageToggle />
        </div>
      </header>

      {/* Main Content */}
      <div className="flex-1 flex items-center justify-center px-4 py-8">
        <Card className="w-full max-w-md">
          <CardHeader className="space-y-1">
            <CardTitle className="text-2xl font-bold text-center">
              {isSignUp ? t('auth.createAccount') : t('auth.welcomeBack')}
            </CardTitle>
            <CardDescription className="text-center">
              {isOAuthRedirect 
                ? t('auth.oauthLoginDescription', { defaultValue: 'Please sign in to authorize the application' })
                : (isSignUp ? t('auth.signUpToAccount') : t('auth.signInToAccount'))
              }
            </CardDescription>
            {isOAuthRedirect && oauthParams.client_id && (
              <div className="mt-2 p-2 bg-blue-50 dark:bg-blue-950 border border-blue-200 dark:border-blue-800 rounded-md">
                <p className="text-sm text-blue-700 dark:text-blue-300 text-center">
                  {t('auth.oauthClientInfo', { 
                    defaultValue: 'Authorizing access for application',
                    client_id: oauthParams.client_id 
                  })}
                </p>
              </div>
            )}
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit} className="space-y-4">
              {error && (
                <Alert variant="destructive">
                  <AlertCircle className="h-4 w-4" />
                  <AlertDescription>
                    {error}
                  </AlertDescription>
                </Alert>
              )}

              {success && (
                <Alert className="border-green-500 text-green-700 bg-green-50 dark:bg-green-950 dark:text-green-300">
                  <AlertCircle className="h-4 w-4" />
                  <AlertDescription>
                    {success}
                  </AlertDescription>
                </Alert>
              )}
              
              {isSignUp && (
                <div className="space-y-2">
                  <label htmlFor="username" className="text-sm font-medium">
                    {t('common.username')}
                  </label>
                  <Input
                    id="username"
                    type="text"
                    placeholder={t('common.username')}
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                    disabled={loading}
                    required
                  />
                </div>
              )}
              
              <div className="space-y-2">
                <label htmlFor="email" className="text-sm font-medium">
                  {t('common.email')}
                </label>
                <Input
                  id="email"
                  type="email"
                  placeholder={t('common.email')}
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  disabled={loading}
                  required
                />
              </div>
              
              <div className="space-y-2">
                <label htmlFor="password" className="text-sm font-medium">
                  {t('common.password')}
                </label>
                <Input
                  id="password"
                  type="password"
                  placeholder={t('common.password')}
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  disabled={loading}
                  required
                />
              </div>
              
              {isSignUp && (
                <div className="space-y-2">
                  <label htmlFor="confirmPassword" className="text-sm font-medium">
                    {t('common.confirmPassword')}
                  </label>
                  <Input
                    id="confirmPassword"
                    type="password"
                    placeholder={t('common.confirmPassword')}
                    value={confirmPassword}
                    onChange={(e) => setConfirmPassword(e.target.value)}
                    disabled={loading}
                    required
                  />
                </div>
              )}
              
              <Button
                type="submit"
                className="w-full"
                disabled={loading}
              >
                {loading ? (
                  <>
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    {t('common.loading')}
                  </>
                ) : (
                  isSignUp ? t('auth.signUp') : t('common.login')
                )}
              </Button>
              
              <div className="text-center">
                <p className="text-sm text-muted-foreground inline">
                  {isSignUp ? t('auth.alreadyHaveAccount') : t('auth.dontHaveAccount')}
                </p>
                <Button
                  type="button"
                  variant="link"
                  onClick={toggleMode}
                  disabled={loading}
                  className="text-sm p-0 h-auto ml-1"
                >
                  {isSignUp ? t('auth.signIn') : t('auth.signUp')}
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
