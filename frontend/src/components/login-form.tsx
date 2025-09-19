import React from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { IconLoader } from '@tabler/icons-react'

import * as z from 'zod'
import { Card, CardContent, CardHeader, CardTitle } from './ui/card'
import { Form, FormControl, FormItem, FormLabel, FormMessage } from './ui/form'
import { cn } from '@/lib/utils'
import { FormField } from './ui/form'
import { Input } from './ui/input'
import { Button } from './ui/button'
import { Link } from '@tanstack/react-router'

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
    </div>
  )
}

export default LoginForm
