import { useEffect } from 'react'

/**
 * Custom hook to update the document title
 * @param title - The page title to set
 * @param suffix - Optional suffix to append (defaults to "MiniAuth")
 */
export function useDocumentTitle(title?: string, suffix?: string) {
  useEffect(() => {
    const defaultSuffix = suffix || 'MiniAuth'
    const newTitle = title ? `${title} - ${defaultSuffix}` : defaultSuffix
    document.title = newTitle
    
    // Cleanup function to reset title when component unmounts
    return () => {
      document.title = 'MiniAuth - OAuth2 Authentication Service'
    }
  }, [title, suffix])
}
