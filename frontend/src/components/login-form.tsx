import { zodResolver } from '@hookform/resolvers/zod'
import { IconLoader } from '@tabler/icons-react'
import React from 'react'
import { useForm } from 'react-hook-form'

import { oauthManager } from '@/lib/oauth'
import { cn } from '@/lib/utils'
import { OAuthProviders } from '@/types/models/oauth_account'
import { IconBrandGithub, IconBrandGoogle } from '@tabler/icons-react'
import { Link, useNavigate } from '@tanstack/react-router'
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

const loginFormSchema = z.object({
  email: z.email(),
  password: z.string(),
  username: z.string().optional(),
})

type LoginFormSchemaType = z.infer<typeof loginFormSchema>

interface ComponentProps extends React.ComponentPropsWithoutRef<'div'> {
  handleSubmit: (values: LoginFormSchemaType) => Promise<void>
  isLoading: boolean
}

const LoginForm: React.FC<ComponentProps> = (props) => {
  const { className, handleSubmit, isLoading } = props
  const form = useForm<LoginFormSchemaType>({
    resolver: zodResolver(loginFormSchema),
    defaultValues: {
      username: undefined,
    },
  })

  async function onSubmit(values: LoginFormSchemaType) {
    await handleSubmit(values)
  }

  const navigate = useNavigate()
  const [oauthLoading, setOauthLoading] = React.useState<string | null>(null)

  const handleOAuthLogin = async (provider: string) => {
    setOauthLoading(provider)
    try {
      await oauthManager.loginWithOAuth(provider as any)
      toast.success(`Logged in with ${provider} successfully!`)
      navigate({ to: '/' })
    } catch (error: any) {
      toast.error(error.message || `Failed to login with ${provider}`)
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
                            <Input
                              {...field}
                              type="email"
                              placeholder="m@example.com"
                            />
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

                  <Button type="submit" className="w-full" disabled={isLoading}>
                    {isLoading ? (
                      <IconLoader className="animate-spin" />
                    ) : (
                      <span>Login</span>
                    )}
                  </Button>
                </div>
                <div className="text-center">
                  <Link
                    to="/request-password-reset"
                    className="text-sm text-muted-foreground hover:text-primary underline underline-offset-4"
                  >
                    Forgot your password?
                  </Link>
                </div>
                <div className="text-center text-sm">
                  Don&apos;t have an account?{' '}
                  <Link to="/register" className="underline underline-offset-4">
                    Register
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
            <IconBrandGoogle className="mr-2 h-4 w-4" />
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
            <IconBrandGithub className="mr-2 h-4 w-4" />
          )}
          GitHub
        </Button>
      </div>
    </div>
  )
}

export default LoginForm
