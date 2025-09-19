import { oauthManager } from '@/lib/oauth'
import { cn } from '@/lib/utils'
import { OAuthProviders } from '@/types/models/oauth_account'
import { zodResolver } from '@hookform/resolvers/zod'
import { IconLoader } from '@tabler/icons-react'
import { createFileRoute, Link, useNavigate } from '@tanstack/react-router'
import { Github } from 'lucide-react'
import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { toast } from 'sonner'
import * as z from 'zod'
import { Button } from './ui/button'
import { Card, CardContent, CardHeader, CardTitle } from './ui/card'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from './ui/form'
import { Input } from './ui/input'
import { Separator } from './ui/separator'
import { Switch } from './ui/switch'

const registerFormSchema = z.object({
  email: z.email(),
  password: z.string(),
  username: z.string().optional(),
  show_email: z.boolean(),
  show_last_login: z.boolean(),
  two_factor_enabled: z.boolean(),
})

type RegisterFormSchemaType = z.infer<typeof registerFormSchema>

interface ComponentProps extends React.ComponentPropsWithoutRef<'div'> {
  submitValues: (values: RegisterFormSchemaType) => Promise<void>
  isLoading: boolean
}

function RegisterFormComponent(props: ComponentProps) {
  const { className, submitValues, isLoading } = props
  const form = useForm<RegisterFormSchemaType>({
    resolver: zodResolver(registerFormSchema),
  })
  const navigate = useNavigate()
  const [oauthLoading, setOauthLoading] = useState<string | null>(null)

  async function onSubmit(values: RegisterFormSchemaType) {
    await submitValues(values)
  }

  const handleOAuthLogin = async (provider: string) => {
    setOauthLoading(provider)
    try {
      await oauthManager.loginWithOAuth(provider as any)
      toast.success(`Signed up with ${provider} successfully!`)
      navigate({ to: '/' })
    } catch (error: any) {
      toast.error(error.message || `Failed to sign up with ${provider}`)
    } finally {
      setOauthLoading(null)
    }
  }

  return (
    <div className={cn('flex flex-col gap-6', className)}>
      <Card>
        <CardHeader className="text-center">
          <CardTitle className="text-xl">Welcome back</CardTitle>
        </CardHeader>
        <CardContent>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)}>
              <div className="grid gap-6">
                <div className="grid gap-6">
                  <div className="grid gap-2">
                    <FormField
                      control={form.control}
                      name="email"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel> Email </FormLabel>
                          <FormControl>
                            <Input {...field} placeholder="m@example.com" />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </div>
                  <div className="grid gap-2">
                    <FormField
                      control={form.control}
                      name="username"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel> Username </FormLabel>
                          <FormControl>
                            <Input {...field} placeholder="james_bond" />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </div>
                  <div className="grid gap-2">
                    <FormField
                      control={form.control}
                      name="password"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel className="flex items-center">
                            Password
                          </FormLabel>
                          <FormControl>
                            <Input
                              {...field}
                              type="password"
                              placeholder="********"
                            />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </div>
                  <div className="grid gap-2">
                    <h3 className="text-md font-bold">Privacy Settings</h3>
                    <FormField
                      control={form.control}
                      name="show_email"
                      render={({ field }) => (
                        <FormItem className="flex flex-row items-center justify-between">
                          <FormLabel className="flex items-center">
                            Show Email
                          </FormLabel>
                          <FormControl>
                            <Switch
                              checked={field.value}
                              onCheckedChange={field.onChange}
                            />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                    <FormField
                      control={form.control}
                      name="show_last_login"
                      render={({ field }) => (
                        <FormItem className="flex flex-row items-center justify-between">
                          <FormLabel className="flex items-center">
                            Show Last Login
                          </FormLabel>
                          <FormControl>
                            <Switch
                              checked={field.value}
                              onCheckedChange={field.onChange}
                            />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                    <FormField
                      control={form.control}
                      name="two_factor_enabled"
                      render={({ field }) => (
                        <FormItem className="flex flex-row items-center justify-between">
                          <FormLabel className="flex items-center">
                            Enable Two Factor Authentication
                          </FormLabel>
                          <FormControl>
                            <Switch
                              checked={field.value}
                              onCheckedChange={field.onChange}
                            />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </div>
                  <Button type="submit" className="w-full" disabled={isLoading}>
                    {isLoading ? (
                      <IconLoader className="animate-spin" />
                    ) : (
                      <span>Register</span>
                    )}
                  </Button>
                </div>
                <div className="text-center text-sm">
                  Already have an account?{' '}
                  <Link to="/login" className="underline underline-offset-4">
                    Login
                  </Link>
                </div>
              </div>
            </form>
          </Form>
        </CardContent>
      </Card>

      <div className="relative">
        <div className="absolute inset-0 flex items-center">
          <Separator className="w-full" />
        </div>
        <div className="relative flex justify-center text-xs uppercase">
          <span className="bg-background px-2 text-muted-foreground">
            Or continue with
          </span>
        </div>
      </div>

      <div className="grid gap-2">
        <Button
          variant="outline"
          type="button"
          disabled={isLoading || oauthLoading !== null}
          onClick={() => handleOAuthLogin(OAuthProviders.GOOGLE)}
        >
          {oauthLoading === OAuthProviders.GOOGLE ? (
            <div className="mr-2 h-4 w-4 animate-spin rounded-full border-2 border-b-transparent" />
          ) : (
            <svg className="mr-2 h-4 w-4" viewBox="0 0 24 24">
              <path
                d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
                fill="#4285F4"
              />
              <path
                d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
                fill="#34A853"
              />
              <path
                d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
                fill="#FBBC05"
              />
              <path
                d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
                fill="#EA4335"
              />
            </svg>
          )}
          Google
        </Button>

        <Button
          variant="outline"
          type="button"
          disabled={isLoading || oauthLoading !== null}
          onClick={() => handleOAuthLogin(OAuthProviders.GITHUB)}
        >
          {oauthLoading === OAuthProviders.GITHUB ? (
            <div className="mr-2 h-4 w-4 animate-spin rounded-full border-2 border-b-transparent" />
          ) : (
            <Github className="mr-2 h-4 w-4" />
          )}
          GitHub
        </Button>
      </div>
    </div>
  )
}

export default RegisterFormComponent
