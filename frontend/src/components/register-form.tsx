import { oauthManager } from '@/lib/oauth'
import { cn } from '@/lib/utils'
import { OAuthProviders } from '@/types/models/oauth_account'
import { zodResolver } from '@hookform/resolvers/zod'
import {
  IconBrandGithub,
  IconBrandGoogle,
  IconLoader,
} from '@tabler/icons-react'
import { Link, useNavigate } from '@tanstack/react-router'
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
    defaultValues: {
      show_email: false,
      show_last_login: false,
      two_factor_enabled: false,
    },
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

export default RegisterFormComponent
