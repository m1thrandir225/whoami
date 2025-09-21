import RegisterForm from '@/components/register-form'
import { Logo } from '@/components/ui/logo'
import authService from '@/services/auth.service'
import { useAuthStore } from '@/stores/auth'
import type { RegisterRequest } from '@/types/api/auth.requests'
import { useMutation } from '@tanstack/react-query'
import { createFileRoute, Link, useNavigate } from '@tanstack/react-router'
import { toast } from 'sonner'

export const Route = createFileRoute('/(auth)/register')({
  component: RouteComponent,
  head: () => ({
    meta: [
      {
        title: 'whoami - Register',
      },
    ],
  }),
})

function RouteComponent() {
  const authStore = useAuthStore()
  const navigate = useNavigate()
  const { mutateAsync, status } = useMutation({
    mutationKey: ['register'],
    mutationFn: (values: RegisterRequest) => authService.register(values),
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
          to="/register"
          className="flex items-center gap-2 self-center font-medium font-mono"
        >
          <Logo variant="full" />
        </Link>
        <RegisterForm
          submitValues={async (values) => {
            await mutateAsync({
              email: values.email,
              password: values.password,
              username: values.username,
              privacy_settings: {
                show_email: values.show_email,
                show_last_login: values.show_last_login,
                two_factor_enabled: values.two_factor_enabled,
              },
            })
          }}
          isLoading={status === 'pending'}
        />
      </div>
    </div>
  )
}
