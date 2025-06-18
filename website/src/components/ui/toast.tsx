import { toast as sonnerToast } from 'sonner'

export const toast = {
  success: (message: string) => sonnerToast.success(message),
  error: (message: string) => sonnerToast.error(message),
  info: (message: string) => sonnerToast.info(message),
  warning: (message: string) => sonnerToast.warning(message),
  promise: sonnerToast.promise,
  loading: (message: string) => sonnerToast.loading(message),
  dismiss: sonnerToast.dismiss,
}
