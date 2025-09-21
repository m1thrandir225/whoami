import { SetPasswordForm } from '@/components/set-password-form'
import { Badge } from '@/components/ui/badge'
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
import { Loader, PageLoader } from '@/components/ui/loader'
import { Separator } from '@/components/ui/separator'
import { Switch } from '@/components/ui/switch'
import authService from '@/services/auth.service'
import oauthService from '@/services/oauth.service'
import userService from '@/services/user.service'
import type {
  UpdatePasswordRequest,
  UpdateUserRequest,
} from '@/types/api/auth.requests'
import type { PrivacySettings } from '@/types/models/privacy_settings'
import { hasPassword } from '@/types/models/user'
import { IconBrandGithub, IconBrandGoogle } from '@tabler/icons-react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import { toast } from 'sonner'

export const Route = createFileRoute('/_pathlessLayout/me')({
  component: RouteComponent,
  head: () => ({
    meta: [
      {
        title: 'whoami - Me',
      },
    ],
  }),
  pendingComponent: () => <PageLoader />,
})

function RouteComponent() {
  const queryClient = useQueryClient()
  const [isEditingProfile, setIsEditingProfile] = useState(false)
  const [isChangingPassword, setIsChangingPassword] = useState(false)

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
      {oauthAccounts && oauthAccounts.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle>Connected Accounts</CardTitle>
            <CardDescription>
              Link your social accounts for quick sign-in and enhanced security.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            {oauthAccounts.map((account) => (
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
      )}

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
