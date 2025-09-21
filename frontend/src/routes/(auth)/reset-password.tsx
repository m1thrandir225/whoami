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
import {
  IconEye,
  IconEyeOff,
  IconLock,
  IconMail,
  IconShield,
} from '@tabler/icons-react'
import { createFileRoute, useNavigate, useSearch } from '@tanstack/react-router'
import { useEffect, useState } from 'react'
import { useForm } from 'react-hook-form'
import { toast } from 'sonner'
import { z } from 'zod'

const resetPasswordSchema = z
  .object({
    password: z.string().min(8, 'Password must be at least 8 characters'),
    confirmPassword: z.string(),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: "Passwords don't match",
    path: ['confirmPassword'],
  })

const otpSchema = z.object({
  otp: z.string().length(6, 'OTP must be 6 digits'),
})

type ResetPasswordFormData = z.infer<typeof resetPasswordSchema>
type OTPFormData = z.infer<typeof otpSchema>

export const Route = createFileRoute('/(auth)/reset-password')({
  component: ResetPasswordPage,
  validateSearch: z.object({
    token: z.string().optional(),
  }),
})

function ResetPasswordPage() {
  const navigate = useNavigate()
  const search = useSearch({ from: '/(auth)/reset-password' })
  const [step, setStep] = useState<
    'verify-token' | 'verify-otp' | 'reset-password'
  >('verify-token')
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [showPassword, setShowPassword] = useState(false)
  const [showConfirmPassword, setShowConfirmPassword] = useState(false)
  const [token, setToken] = useState<string | null>(null)
  const [canResend, setCanResend] = useState(true)
  const [resendCooldown, setResendCooldown] = useState(0)

  const otpForm = useForm<OTPFormData>({
    resolver: zodResolver(otpSchema),
  })

  const passwordForm = useForm<ResetPasswordFormData>({
    resolver: zodResolver(resetPasswordSchema),
  })

  useEffect(() => {
    if (search.token) {
      setToken(search.token)
      verifyToken(search.token)
    }
  }, [search.token])

  useEffect(() => {
    if (resendCooldown > 0) {
      const timer = setTimeout(() => {
        setResendCooldown(resendCooldown - 1)
      }, 1000)
      return () => clearTimeout(timer)
    } else {
      setCanResend(true)
    }
  }, [resendCooldown])

  const verifyToken = async (token: string) => {
    setIsLoading(true)
    setError(null)

    try {
      await passwordResetService.verifyPassword({ token })
      setStep('verify-otp')
      toast.success('Verification code sent to your email')
    } catch (err: any) {
      setError(err.response?.data?.error || 'Invalid or expired reset token')
    } finally {
      setIsLoading(false)
    }
  }

  const verifyOTP = async (data: OTPFormData) => {
    if (!token) return

    setIsLoading(true)
    setError(null)

    try {
      await passwordResetService.verifyResetOTP({ token, otp: data.otp })
      setStep('reset-password')
      toast.success('Code verified successfully')
    } catch (err: any) {
      setError(err.response?.data?.error || 'Invalid verification code')
    } finally {
      setIsLoading(false)
    }
  }

  const resetPassword = async (data: ResetPasswordFormData) => {
    if (!token) return

    setIsLoading(true)
    setError(null)

    try {
      await passwordResetService.resetPassword({
        token,
        new_password: data.password,
      })

      toast.success('Password reset successfully!')
      navigate({ to: '/login' })
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to reset password')
    } finally {
      setIsLoading(false)
    }
  }

  const resendOTP = async () => {
    if (!token || !canResend) return

    setIsLoading(true)
    setError(null)

    try {
      await passwordResetService.verifyPassword({ token })
      toast.success('New verification code sent to your email')

      setCanResend(false)
      setResendCooldown(60)
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to resend code')
    } finally {
      setIsLoading(false)
    }
  }

  if (!token) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <Card className="w-full max-w-md">
          <CardHeader className="text-center">
            <CardTitle>Invalid Reset Link</CardTitle>
            <CardDescription>
              This password reset link is invalid or missing a token.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Button
              onClick={() => navigate({ to: '/login' })}
              className="w-full"
            >
              Back to Login
            </Button>
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
            {step === 'verify-otp' ? (
              <IconShield className="h-6 w-6 text-primary" />
            ) : step === 'reset-password' ? (
              <IconLock className="h-6 w-6 text-primary" />
            ) : (
              <IconMail className="h-6 w-6 text-primary" />
            )}
          </div>
          <CardTitle>
            {step === 'verify-token' && 'Verifying Reset Link'}
            {step === 'verify-otp' && 'Enter Verification Code'}
            {step === 'reset-password' && 'Reset Your Password'}
          </CardTitle>
          <CardDescription>
            {step === 'verify-token' &&
              'Please wait while we verify your reset link...'}
            {step === 'verify-otp' &&
              'Enter the 6-digit code sent to your email'}
            {step === 'reset-password' && 'Enter your new password'}
          </CardDescription>
        </CardHeader>
        <CardContent>
          {error && (
            <Alert variant="destructive" className="mb-4">
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          )}

          {step === 'verify-token' && (
            <div className="text-center">
              <Loader size="lg" />
              <p className="mt-4 text-sm text-muted-foreground">
                Verifying your reset link...
              </p>
            </div>
          )}

          {step === 'verify-otp' && (
            <form
              onSubmit={otpForm.handleSubmit(verifyOTP)}
              className="space-y-4"
            >
              <div className="space-y-2">
                <Label htmlFor="otp">Verification Code</Label>
                <Input
                  id="otp"
                  type="text"
                  placeholder="000000"
                  maxLength={6}
                  {...otpForm.register('otp')}
                  disabled={isLoading}
                  className="text-center text-lg tracking-widest"
                />
                {otpForm.formState.errors.otp && (
                  <p className="text-sm text-destructive">
                    {otpForm.formState.errors.otp.message}
                  </p>
                )}
              </div>

              <Button type="submit" className="w-full" disabled={isLoading}>
                {isLoading ? <Loader size="sm" /> : 'Verify Code'}
              </Button>

              <div className="text-center">
                <Button
                  type="button"
                  variant="link"
                  onClick={resendOTP}
                  disabled={isLoading || !canResend}
                  className="text-sm"
                >
                  {!canResend
                    ? `Resend in ${resendCooldown}s`
                    : "Didn't receive the code? Resend"}
                </Button>
              </div>
            </form>
          )}

          {step === 'reset-password' && (
            <form
              onSubmit={passwordForm.handleSubmit(resetPassword)}
              className="space-y-4"
            >
              <div className="space-y-2">
                <Label htmlFor="password">New Password</Label>
                <div className="relative">
                  <Input
                    id="password"
                    type={showPassword ? 'text' : 'password'}
                    {...passwordForm.register('password')}
                    placeholder="Enter your new password"
                    disabled={isLoading}
                  />
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    className="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent"
                    onClick={() => setShowPassword(!showPassword)}
                    disabled={isLoading}
                  >
                    {showPassword ? (
                      <IconEyeOff className="h-4 w-4" />
                    ) : (
                      <IconEye className="h-4 w-4" />
                    )}
                  </Button>
                </div>
                {passwordForm.formState.errors.password && (
                  <p className="text-sm text-destructive">
                    {passwordForm.formState.errors.password.message}
                  </p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="confirmPassword">Confirm Password</Label>
                <div className="relative">
                  <Input
                    id="confirmPassword"
                    type={showConfirmPassword ? 'text' : 'password'}
                    {...passwordForm.register('confirmPassword')}
                    placeholder="Confirm your new password"
                    disabled={isLoading}
                  />
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    className="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent"
                    onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                    disabled={isLoading}
                  >
                    {showConfirmPassword ? (
                      <IconEyeOff className="h-4 w-4" />
                    ) : (
                      <IconEye className="h-4 w-4" />
                    )}
                  </Button>
                </div>
                {passwordForm.formState.errors.confirmPassword && (
                  <p className="text-sm text-destructive">
                    {passwordForm.formState.errors.confirmPassword.message}
                  </p>
                )}
              </div>

              <Button type="submit" className="w-full" disabled={isLoading}>
                {isLoading ? <Loader size="sm" /> : 'Reset Password'}
              </Button>
            </form>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
