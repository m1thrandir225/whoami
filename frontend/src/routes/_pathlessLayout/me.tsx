import { createFileRoute } from '@tanstack/react-router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Switch } from '@/components/ui/switch'
import { Separator } from '@/components/ui/separator'
import { Badge } from '@/components/ui/badge'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '@/components/ui/alert-dialog'
import { useState } from 'react'
import { toast } from 'sonner'
import { Github, Unlink } from 'lucide-react'
import authService from '@/services/auth.service'
import userService from '@/services/user.service'
import oauthService from '@/services/oauth.service'
import { oauthManager } from '@/lib/oauth'
import { OAuthProviders } from '@/types/models/oauth_account'
import type {
  UpdateUserRequest,
  UpdatePasswordRequest,
} from '@/types/api/auth.requests'
import type { PrivacySettings } from '@/types/models/privacy_settings'
import { hasPassword } from '@/types/models/user'
import { SetPasswordForm } from '@/components/set-password-form'
import { Loader } from '@/components/ui/loader'
import { IconBrandGoogle, IconBrandGithub } from '@tabler/icons-react'

export const Route = createFileRoute('/_pathlessLayout/me')({
  component: RouteComponent,
})

function RouteComponent() {
  const queryClient = useQueryClient()
  const [isEditingProfile, setIsEditingProfile] = useState(false)
  const [isChangingPassword, setIsChangingPassword] = useState(false)
  const [oauthLoading, setOauthLoading] = useState<string | null>(null)

  // Get current user data
  const { data: user, isLoading } = useQuery({
    queryKey: ['current-user'],
    queryFn: authService.getCurrentUser,
  })

  // Get OAuth accounts
  const { data: oauthAccounts } = useQuery({
    queryKey: ['oauth-accounts'],
    queryFn: oauthService.getOAuthAccounts,
  })

  // Profile update mutation
  const profileMutation = useMutation({
    mutationFn: (data: UpdateUserRequest) =>
      userService.updateUser(user?.id || 0, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['current-user'] })
      setIsEditingProfile(false)
      toast.success('Profile updated successfully')
    },
    onError: (error) => {
      toast.error(error.message || 'Failed to update profile')
    },
  })

  // Privacy settings mutation
  const privacyMutation = useMutation({
    mutationFn: (data: PrivacySettings) =>
      userService.updateUserPrivacySettings(user?.id || 0, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['current-user'] })
      toast.success('Privacy settings updated successfully')
    },
    onError: (error) => {
      toast.error(error.message || 'Failed to update privacy settings')
    },
  })

  // Password change mutation
  const passwordMutation = useMutation({
    mutationFn: (data: UpdatePasswordRequest) =>
      userService.updatePassword(data),
    onSuccess: () => {
      setIsChangingPassword(false)
      toast.success('Password changed successfully')
    },
    onError: (error) => {
      toast.error(error.message || 'Failed to change password')
    },
  })

  // Unlink OAuth account mutation
  const unlinkOAuthMutation = useMutation({
    mutationFn: (provider: string) => oauthService.unlinkOAuthAccount(provider),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['oauth-accounts'] })
      toast.success('OAuth account unlinked successfully')
    },
    onError: (error) => {
      toast.error(error.message || 'Failed to unlink OAuth account')
    },
  })

  const handleProfileSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    const formData = new FormData(e.currentTarget)
    profileMutation.mutate({
      email: formData.get('email') as string,
      username: formData.get('username') as string,
    })
  }

  const handlePasswordSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    const formData = new FormData(e.currentTarget)
    passwordMutation.mutate({
      current_password: formData.get('current_password') as string,
      new_password: formData.get('new_password') as string,
    })
  }

  const handlePrivacyChange = (key: keyof PrivacySettings, value: boolean) => {
    if (!user) return

    privacyMutation.mutate({
      ...user.privacy_settings,
      [key]: value,
    })
  }

  const handleLinkOAuth = async (provider: string) => {
    setOauthLoading(provider)
    try {
      await oauthManager.linkAccount(provider as any)
      queryClient.invalidateQueries({ queryKey: ['oauth-accounts'] })
      toast.success(`${provider} account linked successfully!`)
    } catch (error: any) {
      toast.error(error.message || `Failed to link ${provider} account`)
    } finally {
      setOauthLoading(null)
    }
  }

  const getProviderIcon = (provider: string) => {
    switch (provider) {
      case 'google':
        return (
          <svg className="h-4 w-4" viewBox="0 0 24 24">
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
        )
      case 'github':
        return <Github className="h-4 w-4" />
      default:
        return null
    }
  }

  const isProviderLinked = (provider: string) => {
    return oauthAccounts?.some((account) => account.provider === provider)
  }

  const handleUnlinkAccount = async (provider: string) => {
    if (!user || !hasPassword(user)) {
      toast.error('You must set a password before disconnecting OAuth accounts')
      return
    }

    try {
      await oauthService.unlinkOAuthAccount(provider)
      toast.success(`${provider} account disconnected successfully`)
      // Refresh OAuth accounts
      // The original code had oauthAccountsQuery.refetch(), but oauthAccountsQuery is not defined.
      // Assuming the intent was to invalidate the query if it were defined.
      // For now, we'll just toast the success message.
      queryClient.invalidateQueries({ queryKey: ['oauth-accounts'] })
    } catch (error: any) {
      toast.error(error.response?.data?.error || 'Failed to disconnect account')
    }
  }

  if (isLoading) {
    return <div>Loading...</div>
  }

  if (!user) {
    return <div>Error loading user data</div>
  }

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-3xl font-bold tracking-tight">Profile</h2>
        <p className="text-muted-foreground">
          Manage your account settings and privacy preferences.
        </p>
      </div>

      {/* Profile Information */}
      <Card>
        <CardHeader>
          <CardTitle>Profile Information</CardTitle>
          <CardDescription>
            Your basic account information and status.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          {isEditingProfile ? (
            <form onSubmit={handleProfileSubmit} className="space-y-4">
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="email">Email</Label>
                  <Input
                    id="email"
                    name="email"
                    type="email"
                    defaultValue={user.email}
                    required
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="username">Username</Label>
                  <Input
                    id="username"
                    name="username"
                    defaultValue={user.username}
                    required
                  />
                </div>
              </div>
              <div className="flex gap-2">
                <Button type="submit" disabled={profileMutation.isPending}>
                  {profileMutation.isPending ? 'Saving...' : 'Save Changes'}
                </Button>
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => setIsEditingProfile(false)}
                >
                  Cancel
                </Button>
              </div>
            </form>
          ) : (
            <div className="space-y-4">
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <Label>Email</Label>
                  <div className="flex items-center gap-2 mt-1">
                    <span>{user.email}</span>
                    {user.email_verified ? (
                      <Badge variant="secondary">Verified</Badge>
                    ) : (
                      <Badge variant="destructive">Unverified</Badge>
                    )}
                  </div>
                </div>
                <div>
                  <Label>Username</Label>
                  <p className="mt-1">{user.username}</p>
                </div>
                <div>
                  <Label>Role</Label>
                  <p className="mt-1 capitalize">{user.role}</p>
                </div>
                <div>
                  <Label>Status</Label>
                  <div className="mt-1">
                    {user.active ? (
                      <Badge variant="secondary">Active</Badge>
                    ) : (
                      <Badge variant="destructive">Inactive</Badge>
                    )}
                  </div>
                </div>
                <div>
                  <Label>Member Since</Label>
                  <p className="mt-1">
                    {new Date(user.created_at).toLocaleDateString()}
                  </p>
                </div>
                <div>
                  <Label>Last Login</Label>
                  <p className="mt-1">
                    {user.last_login_at
                      ? new Date(user.last_login_at).toLocaleString()
                      : 'Never'}
                  </p>
                </div>
              </div>
              <Button onClick={() => setIsEditingProfile(true)}>
                Edit Profile
              </Button>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Connected Accounts */}
      <Card>
        <CardHeader>
          <CardTitle>Connected Accounts</CardTitle>
          <CardDescription>
            Link your social accounts for quick sign-in and enhanced security.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          {oauthAccounts?.map((account) => (
            <div
              key={account.provider}
              className="flex items-center justify-between p-4 border rounded-lg"
            >
              <div className="flex items-center space-x-3">
                <div className="w-8 h-8 bg-gray-100 rounded-full flex items-center justify-center">
                  {account.provider === 'google' ? (
                    <IconBrandGoogle className="h-4 w-4" />
                  ) : (
                    <IconBrandGithub className="h-4 w-4" />
                  )}
                </div>
                <div>
                  <p className="font-medium capitalize">{account.provider}</p>
                  <p className="text-sm text-muted-foreground">
                    Connected{' '}
                    {new Date(account.created_at!).toLocaleDateString()}
                  </p>
                </div>
              </div>
              <Button
                variant="outline"
                size="sm"
                onClick={() => handleUnlinkAccount(account.provider)}
                disabled={!hasPassword(user) || unlinkOAuthMutation.isPending}
              >
                {unlinkOAuthMutation.isPending ? (
                  <Loader size="sm" />
                ) : (
                  'Disconnect'
                )}
              </Button>
            </div>
          ))}
        </CardContent>
      </Card>

      {/* Privacy Settings */}
      <Card>
        <CardHeader>
          <CardTitle>Privacy Settings</CardTitle>
          <CardDescription>
            Control what information is visible to others.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex items-center justify-between">
            <div>
              <Label>Show Email</Label>
              <p className="text-sm text-muted-foreground">
                Allow others to see your email address
              </p>
            </div>
            <Switch
              checked={user.privacy_settings.show_email}
              onCheckedChange={(value) =>
                handlePrivacyChange('show_email', value)
              }
              disabled={privacyMutation.isPending}
            />
          </div>
          <Separator />
          <div className="flex items-center justify-between">
            <div>
              <Label>Show Last Login</Label>
              <p className="text-sm text-muted-foreground">
                Display when you were last active
              </p>
            </div>
            <Switch
              checked={user.privacy_settings.show_last_login}
              onCheckedChange={(value) =>
                handlePrivacyChange('show_last_login', value)
              }
              disabled={privacyMutation.isPending}
            />
          </div>
          <Separator />
          <div className="flex items-center justify-between">
            <div>
              <Label>Two-Factor Authentication</Label>
              <p className="text-sm text-muted-foreground">
                Enable additional security for your account
              </p>
            </div>
            <Switch
              checked={user.privacy_settings.two_factor_enabled}
              onCheckedChange={(value) =>
                handlePrivacyChange('two_factor_enabled', value)
              }
              disabled={privacyMutation.isPending}
            />
          </div>
        </CardContent>
      </Card>

      {/* Password Management */}
      {!hasPassword(user) ? (
        <Card>
          <CardHeader>
            <CardTitle>Set Password</CardTitle>
            <CardDescription>
              You signed up with OAuth. Set a password to enable email/password
              login.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <SetPasswordForm
              onSuccess={() => {
                // Refresh user data or show success message
                toast.success('Password set successfully!')
              }}
            />
          </CardContent>
        </Card>
      ) : (
        <Card>
          <CardHeader>
            <CardTitle>Password</CardTitle>
            <CardDescription>
              Change your password to keep your account secure.
            </CardDescription>
          </CardHeader>
          <CardContent>
            {isChangingPassword ? (
              <form onSubmit={handlePasswordSubmit} className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="current_password">Current Password</Label>
                  <Input
                    id="current_password"
                    name="current_password"
                    type="password"
                    required
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="new_password">New Password</Label>
                  <Input
                    id="new_password"
                    name="new_password"
                    type="password"
                    required
                    minLength={8}
                  />
                </div>
                <div className="flex gap-2">
                  <Button type="submit" disabled={passwordMutation.isPending}>
                    {passwordMutation.isPending
                      ? 'Changing...'
                      : 'Change Password'}
                  </Button>
                  <Button
                    type="button"
                    variant="outline"
                    onClick={() => setIsChangingPassword(false)}
                  >
                    Cancel
                  </Button>
                </div>
              </form>
            ) : (
              <div className="space-y-4">
                <div>
                  <Label>Password last changed</Label>
                  <p className="mt-1 text-sm text-muted-foreground">
                    {new Date(user.password_changed_at!).toLocaleDateString()}
                  </p>
                </div>
                <Button onClick={() => setIsChangingPassword(true)}>
                  Change Password
                </Button>
              </div>
            )}
          </CardContent>
        </Card>
      )}
    </div>
  )
}
