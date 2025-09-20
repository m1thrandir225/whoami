import {
  Outlet,
  createRootRouteWithContext,
  redirect,
} from '@tanstack/react-router'
import { TanStackRouterDevtoolsPanel } from '@tanstack/react-router-devtools'
import { TanstackDevtools } from '@tanstack/react-devtools'

import TanStackQueryDevtools from '../integrations/tanstack-query/devtools'

import type { QueryClient } from '@tanstack/react-query'
import { useAuthStore } from '@/stores/auth'
import { Toaster } from '@/components/ui/sonner'
import { PageLoader } from '@/components/ui/loader'

interface MyRouterContext {
  queryClient: QueryClient
}

export const Route = createRootRouteWithContext<MyRouterContext>()({
  beforeLoad: async () => {
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
          search: {
            redirect: location.href,
          },
        })
      }
    }
  },
  pendingComponent: () => <PageLoader />,
  component: () => (
    <>
      <Outlet />
      <Toaster />
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
    </>
  ),
})
