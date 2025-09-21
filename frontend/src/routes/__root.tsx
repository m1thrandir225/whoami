import { TanstackDevtools } from '@tanstack/react-devtools'
import {
  Outlet,
  createRootRouteWithContext,
  redirect,
} from '@tanstack/react-router'
import { TanStackRouterDevtoolsPanel } from '@tanstack/react-router-devtools'

import TanStackQueryDevtools from '../integrations/tanstack-query/devtools'

import { PageLoader } from '@/components/ui/loader'
import { Toaster } from '@/components/ui/sonner'
import { useAuthStore } from '@/stores/auth'
import type { QueryClient } from '@tanstack/react-query'

interface MyRouterContext {
  queryClient: QueryClient
}

export const Route = createRootRouteWithContext<MyRouterContext>()({
  beforeLoad: async ({ location }) => {
    const authStore = useAuthStore.getState()
    const isAuthenticated = authStore.isAuthenticated()

    const publicRoutes = [
      '/login',
      '/register',
      '/oauth-callback',
      '/reset-password',
      '/request-password-reset',
    ]

    if (!isAuthenticated) {
      if (!publicRoutes.includes(location.pathname)) {
        throw redirect({
          to: '/login',
        })
      }
    }
  },
  pendingComponent: () => <PageLoader />,
  component: () => (
    <>
      <Outlet />
      <Toaster />
      {process.env.NODE_ENV === 'development' && (
        <TanstackDevtools
          config={{
            position: 'bottom-left',
          }}
          plugins={[
            {
              name: 'Tanstack Router',
              render: <TanStackRouterDevtoolsPanel />,
            },
            TanStackQueryDevtools,
          ]}
        />
      )}
    </>
  ),
})
