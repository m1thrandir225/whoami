import RegisterForm from '@/components/register-form'
import authService from '@/services/auth.service'
import type { RegisterRequest } from '@/types/api/auth.requests'
import { IconZoom } from '@tabler/icons-react'
import { useMutation } from '@tanstack/react-query'
import { createFileRoute, Link } from '@tanstack/react-router'

export const Route = createFileRoute('/(auth)/register')({
  component: RouteComponent,
})

function RouteComponent() {
  const { mutateAsync, status } = useMutation({
    mutationKey: ['register'],
    mutationFn: (values: RegisterRequest) => authService.register(values),
  })

  return (
    <div className="flex min-h-svh flex-col items-center justify-center gap-6 bg-muted p-6 md:p-10">
      <div className="flex w-full max-w-sm flex-col gap-6">
        <Link
          to="/"
          className="flex items-center gap-2 self-center font-medium font-mono"
        >
          <IconZoom />
          whoami
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
