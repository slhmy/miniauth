import { useState } from 'react'
import type { ReactNode } from 'react'
import { useTranslation } from 'react-i18next'
import { 
  Users, 
  Menu, 
  X, 
  LogOut,
  Shield,
  User,
  ChevronDown,
  Settings
} from 'lucide-react'
import type { LucideIcon } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'
import { ThemeToggle } from '@/components/theme-toggle'
import { LanguageToggle } from '@/components/language-toggle'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { useIsMobile } from '@/hooks/use-mobile'
import { useAuth } from '@/hooks/useAuth'
import { useNavigate, useLocation } from 'react-router-dom'

interface MainLayoutProps {
  children: ReactNode
}

interface NavigationItem {
  id: string
  labelKey: string
  icon: LucideIcon
  path: string
  adminOnly?: boolean
}

const getNavigationItems = (): NavigationItem[] => [
  {
    id: 'profile',
    labelKey: 'navigation.profile',
    icon: User,
    path: '/profile'
  },
  {
    id: 'oauth-applications',
    labelKey: 'oauth.applications',
    icon: Settings,
    path: '/oauth/applications',
    adminOnly: true
  },
  {
    id: 'user-management',
    labelKey: 'navigation.userManagement',
    icon: Users,
    path: '/admin/users',
    adminOnly: true
  },
]

export function MainLayout({ children }: MainLayoutProps) {
  const [sidebarOpen, setSidebarOpen] = useState(false)
  const isMobile = useIsMobile()
  const { logout, user } = useAuth()
  const navigate = useNavigate()
  const location = useLocation()
  const { t } = useTranslation()

  // Check if user is admin
  const isAdmin = user?.role === 'admin'

  // Get navigation items with current translations
  const navigationItems = getNavigationItems()

  // Filter navigation items based on admin status
  const filteredNavigationItems = navigationItems.filter(item => 
    !item.adminOnly || isAdmin
  )

  const handleLogout = async () => {
    await logout()
    navigate('/login')
  }

  const handleNavigation = (path: string) => {
    navigate(path)
    if (isMobile) {
      setSidebarOpen(false)
    }
  }

  // UserMenu component for the top bar
  const UserMenu = () => (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost" className="flex items-center space-x-2 text-sm hover:bg-accent">
          <User className="h-4 w-4" />
          <span>{user?.username || t('common.user')}</span>
          {isAdmin && (
            <span className="ml-2 inline-flex items-center px-2 py-1 rounded-full text-xs bg-primary/10 text-primary">
              {t('common.admin')}
            </span>
          )}
          <ChevronDown className="h-4 w-4" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-56">
        <DropdownMenuLabel>
          <div className="flex flex-col space-y-1">
            <p className="text-sm font-medium leading-none">{user?.username || t('common.user')}</p>
            <p className="text-xs leading-none text-muted-foreground">
              {user?.email || t('profile.noEmail')}
            </p>
          </div>
        </DropdownMenuLabel>
        <DropdownMenuSeparator />
        <DropdownMenuItem onClick={() => handleNavigation('/profile')}>
          <User className="mr-2 h-4 w-4" />
          <span>{t('common.profile')}</span>
        </DropdownMenuItem>
        {isAdmin && (
          <>
            <DropdownMenuSeparator />
            <DropdownMenuItem onClick={() => handleNavigation('/admin/users')}>
              <Users className="mr-2 h-4 w-4" />
              <span>{t('navigation.userManagement')}</span>
            </DropdownMenuItem>
          </>
        )}
        <DropdownMenuSeparator />
        <DropdownMenuItem onClick={handleLogout} className="text-red-600 dark:text-red-400">
          <LogOut className="mr-2 h-4 w-4" />
          <span>{t('common.logout')}</span>
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  )

  const Sidebar = () => (
    <div className="flex h-full flex-col border-r bg-card text-card-foreground">
      {/* Header */}
      <div className="flex h-16 items-center justify-between px-6 border-b border-border">
        <div className="flex items-center space-x-2">
          <Shield className="h-6 w-6 text-primary" />
          <span className="text-lg font-semibold text-foreground">
            MiniAuth {isAdmin && <span className="text-sm text-muted-foreground">• {t('common.admin')}</span>}
          </span>
        </div>
        {isMobile && (
          <Button
            variant="ghost"
            size="sm"
            onClick={() => setSidebarOpen(false)}
            className="h-8 w-8 p-0 text-foreground hover:text-accent-foreground"
          >
            <X className="h-5 w-5" />
          </Button>
        )}
      </div>

      {/* Navigation */}
      <nav className="flex-1 space-y-1 p-4">
        {/* Regular Navigation */}
        {filteredNavigationItems.filter(item => !item.adminOnly).map((item) => {
          const Icon = item.icon
          const isActive = location.pathname === item.path
          return (
            <button
              key={item.id}
              onClick={() => handleNavigation(item.path)}
              className={`flex w-full items-center space-x-3 rounded-md px-3 py-2 text-sm font-medium transition-colors ${
                isActive 
                  ? 'bg-accent text-accent-foreground' 
                  : 'text-muted-foreground hover:bg-accent hover:text-accent-foreground'
              }`}
            >
              <Icon className="h-5 w-5" />
              <span>{t(item.labelKey)}</span>
            </button>
          )
        })}

        {/* Admin Section */}
        {isAdmin && (
          <>
            <Separator className="my-4 bg-border" />
            <div className="px-3 py-2">
              <p className="text-xs font-semibold text-muted-foreground uppercase tracking-wider">
                {t('navigation.administration')}
              </p>
            </div>
            {filteredNavigationItems.filter(item => item.adminOnly).map((item) => {
              const Icon = item.icon
              const isActive = location.pathname === item.path
              return (
                <button
                  key={item.id}
                  onClick={() => handleNavigation(item.path)}
                  className={`flex w-full items-center space-x-3 rounded-md px-3 py-2 text-sm font-medium transition-colors ${
                    isActive 
                      ? 'bg-accent text-accent-foreground' 
                      : 'text-muted-foreground hover:bg-accent hover:text-accent-foreground'
                  }`}
                >
                  <Icon className="h-5 w-5" />
                  <span>{t(item.labelKey)}</span>
                </button>
              )
            })}
          </>
        )}
      </nav>

      {/* Removed the Logout button from sidebar footer */}
    </div>
  )

  return (
    <div className="flex h-screen w-full bg-background">
      {/* Desktop Sidebar */}
      {!isMobile && (
        <div className="w-64 flex-shrink-0 shadow-sm">
          <Sidebar />
        </div>
      )}

      {/* Mobile Sidebar Overlay */}
      {isMobile && sidebarOpen && (
        <div className="fixed inset-0 z-50 flex">
          <div 
            className="flex-1 bg-black/80"
            onClick={() => setSidebarOpen(false)}
          />
          <div className="w-64 flex-shrink-0">
            <Sidebar />
          </div>
        </div>
      )}

      {/* Main Content */}
      <div className="flex flex-1 flex-col overflow-hidden min-w-0 w-full">
        {/* Top Bar */}
        <header className="flex h-16 items-center justify-between border-b bg-background px-6 shadow-sm flex-shrink-0 w-full">
          {isMobile && (
            <Button
              variant="ghost"
              size="sm"
              onClick={() => setSidebarOpen(true)}
              className="h-8 w-8 p-0 text-foreground hover:text-accent-foreground"
            >
              <Menu className="h-5 w-5" />
            </Button>
          )}
          
          <div className="flex items-center space-x-4">
            {/* Page title could be added here based on current route */}
          </div>

          <div className="flex items-center space-x-4">
            {/* Theme Toggle */}
            <ThemeToggle />
            
            {/* Language Toggle */}
            <LanguageToggle />
            
            {/* User Menu - Replaces old user info display */}
            <UserMenu />
          </div>
        </header>

        {/* Main Content Area */}
        <main className="flex-1 overflow-auto bg-background w-full">
          <div className="p-6 w-full h-full">
            {children}
          </div>
        </main>
      </div>
    </div>
  )
}
