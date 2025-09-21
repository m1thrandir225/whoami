import LoginForm from '@/components/login-form'
import { Logo } from '@/components/ui/logo'
import authService from '@/services/auth.service'
import { useAuthStore } from '@/stores/auth'
import type { LoginRequest } from '@/types/api/auth.requests'
import { useMutation } from '@tanstack/react-query'
import { createFileRoute, Link, useNavigate } from '@tanstack/react-router'
import { toast } from 'sonner'

export const Route = createFileRoute('/(auth)/login')({
  component: RouteComponent,
  head: () => ({
    meta: [
      {
        title: 'whoami - Login',
      },
    ],
  }),
})

function RouteComponent() {
  const authStore = useAuthStore()
  const navigate = useNavigate()
  const { mutateAsync, status } = useMutation({
    mutationKey: ['login'],
    mutationFn: (values: LoginRequest) => authService.login(values),
    onSuccess: (data) => {
      authStore.login(data)

      navigate({ to: '/' })
    },
    onError: (error) => {
      toast.error(error.message)
    },
  })

  return (
    <div className="flex min-h-svh flex-col items-center justify-center gap-6 bg-muted p-6 md:p-10">
      <div className="flex w-full max-w-sm flex-col gap-6">
        <Link
          to="/login"
          className="flex items-center gap-2 self-center font-medium font-mono"
        >
          <Logo variant="full" />
        </Link>
        <LoginForm
          handleSubmit={async (values) => {
            await mutateAsync(values)
          }}
          isLoading={status === 'pending'}
        />
      </div>
    </div>
  )
}
