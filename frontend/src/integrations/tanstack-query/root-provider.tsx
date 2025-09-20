import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { isAxiosError } from 'axios'

const MAX_RETRIES = 3
const HTTP_STATUS_TO_NOT_RETRY = [400, 401, 403, 404]

export function getContext() {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        retry: (failureCount, error) => {
          if (failureCount > MAX_RETRIES) {
            return false
          }
          return !(
            isAxiosError(error) &&
            HTTP_STATUS_TO_NOT_RETRY.includes(error.response?.status || 0)
          )
        },
      },
    },
  })
  return {
    queryClient,
  }
}

export function Provider({
  children,
  queryClient,
}: {
  children: React.ReactNode
  queryClient: QueryClient
}) {
  return (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  )
}
