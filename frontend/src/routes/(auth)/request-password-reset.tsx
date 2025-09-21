import { Alert, AlertDescription } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Loader } from '@/components/ui/loader'
import passwordResetService from '@/services/password-reset.service'
import { zodResolver } from '@hookform/resolvers/zod'
import { IconArrowLeft, IconCheck, IconMail } from '@tabler/icons-react'
import { createFileRoute, Link } from '@tanstack/react-router'
import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { toast } from 'sonner'
import { z } from 'zod'

const requestPasswordResetSchema = z.object({
  email: z.email(),
})

type RequestPasswordResetFormData = z.infer<typeof requestPasswordResetSchema>

export const Route = createFileRoute('/(auth)/request-password-reset')({
  component: RequestPasswordResetPage,
})

function RequestPasswordResetPage() {
  const [isLoading, setIsLoading] = useState(false)
  const [isSubmitted, setIsSubmitted] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const {
    register,
    handleSubmit,
    formState: { errors },
    watch,
  } = useForm<RequestPasswordResetFormData>({
    resolver: zodResolver(requestPasswordResetSchema),
  })

  const email = watch('email')

  const onSubmit = async (data: RequestPasswordResetFormData) => {
    setIsLoading(true)
    setError(null)

    try {
      await passwordResetService.requestPasswordReset({ email: data.email })
      setIsSubmitted(true)
      toast.success('Password reset instructions sent to your email')
    } catch (err: any) {
      setError(
        err.response?.data?.error || 'Failed to send password reset email',
      )
    } finally {
      setIsLoading(false)
    }
  }

  const handleResend = async () => {
    if (!email) return

    setIsLoading(true)
    setError(null)

    try {
      await passwordResetService.requestPasswordReset({ email })
      toast.success('Password reset instructions sent again')
    } catch (err: any) {
      setError(
        err.response?.data?.error || 'Failed to resend password reset email',
      )
    } finally {
      setIsLoading(false)
    }
  }

  if (isSubmitted) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <Card className="w-full max-w-md">
          <CardHeader className="text-center">
            <div className="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-green-100">
              <IconCheck className="h-6 w-6 text-green-600" />
            </div>
            <CardTitle>Check Your Email</CardTitle>
            <CardDescription>
              We've sent password reset instructions to <strong>{email}</strong>
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <Alert>
              <IconMail className="h-4 w-4" />
              <AlertDescription>
                Click the link in the email to reset your password. The link
                will expire in 1 hour.
              </AlertDescription>
            </Alert>

            <div className="space-y-3">
              <Button
                onClick={handleResend}
                disabled={isLoading}
                variant="outline"
                className="w-full"
              >
                {isLoading ? <Loader size="sm" /> : 'Resend Email'}
              </Button>

              <Button
                onClick={() => setIsSubmitted(false)}
                variant="ghost"
                className="w-full"
              >
                Use Different Email
              </Button>

              <div className="text-center">
                <Link
                  to="/login"
                  className="text-sm text-muted-foreground hover:text-primary inline-flex items-center gap-1"
                >
                  <IconArrowLeft className="h-3 w-3" />
                  Back to Login
                </Link>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    )
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <Card className="w-full max-w-md">
        <CardHeader className="text-center">
          <div className="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-primary/10">
            <IconMail className="h-6 w-6 text-primary" />
          </div>
          <CardTitle>Reset Your Password</CardTitle>
          <CardDescription>
            Enter your email address and we'll send you a link to reset your
            password.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            {error && (
              <Alert variant="destructive">
                <AlertDescription>{error}</AlertDescription>
              </Alert>
            )}

            <div className="space-y-2">
              <Label htmlFor="email">Email Address</Label>
              <Input
                id="email"
                type="email"
                placeholder="Enter your email address"
                {...register('email')}
                disabled={isLoading}
                autoComplete="email"
              />
              {errors.email && (
                <p className="text-sm text-destructive">
                  {errors.email.message}
                </p>
              )}
            </div>

            <Button type="submit" className="w-full" disabled={isLoading}>
              {isLoading ? <Loader size="sm" /> : 'Send Reset Link'}
            </Button>

            <div className="text-center">
              <Link
                to="/login"
                className="text-sm text-muted-foreground hover:text-primary inline-flex items-center gap-1"
              >
                <IconArrowLeft className="h-3 w-3" />
                Back to Login
              </Link>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}
