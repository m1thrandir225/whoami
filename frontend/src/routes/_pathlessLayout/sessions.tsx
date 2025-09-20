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
import { Trash2, Shield, Smartphone, Monitor, Tablet } from 'lucide-react'
import { toast } from 'sonner'
import sessionService from '@/services/session.service'
import { useAuthStore } from '@/stores/auth'

export const Route = createFileRoute('/_pathlessLayout/sessions')({
  component: RouteComponent,
})

function RouteComponent() {
  const queryClient = useQueryClient()
  const accessToken = useAuthStore((state) => state.accessToken)

  const { data, isLoading } = useQuery({
    queryKey: ['user-sessions'],
    queryFn: () => sessionService.getUserSessions(),
  })

  const revokeSessionMutation = useMutation({
    mutationFn: (token: string) => sessionService.revokeSession(token),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['user-sessions'] })
      toast.success('Session revoked successfully')
    },
    onError: (error) => {
      toast.error(error.message || 'Failed to revoke session')
    },
  })

  const revokeAllSessionsMutation = useMutation({
    mutationFn: () => sessionService.revokeAllSessions(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['user-sessions'] })
      toast.success('All sessions revoked successfully')
    },
    onError: (error) => {
      toast.error(error.message || 'Failed to revoke all sessions')
    },
  })

  const getDeviceIcon = (userAgent: string) => {
    const ua = userAgent.toLowerCase()
    if (
      ua.includes('mobile') ||
      ua.includes('android') ||
      ua.includes('iphone')
    ) {
      return <Smartphone className="h-4 w-4" />
    }
    if (ua.includes('tablet') || ua.includes('ipad')) {
      return <Tablet className="h-4 w-4" />
    }
    return <Monitor className="h-4 w-4" />
  }

  const getLocation = (ipAddress: string) => {
    // In a real app, you'd use a geolocation service
    return `IP: ${ipAddress}`
  }

  if (isLoading) {
    return <div>Loading sessions...</div>
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-3xl font-bold tracking-tight">Active Sessions</h2>
          <p className="text-muted-foreground">
            Manage your active sessions across all devices.
          </p>
        </div>
        <AlertDialog>
          <AlertDialogTrigger asChild>
            <Button variant="destructive" disabled={!data?.sessions?.length}>
              Revoke All Sessions
            </Button>
          </AlertDialogTrigger>
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle>Revoke All Sessions</AlertDialogTitle>
              <AlertDialogDescription>
                This will sign you out of all devices and locations. You'll need
                to sign in again.
              </AlertDialogDescription>
            </AlertDialogHeader>
            <AlertDialogFooter>
              <AlertDialogCancel>Cancel</AlertDialogCancel>
              <AlertDialogAction
                onClick={() => revokeAllSessionsMutation.mutate()}
                disabled={revokeAllSessionsMutation.isPending}
              >
                {revokeAllSessionsMutation.isPending
                  ? 'Revoking...'
                  : 'Revoke All'}
              </AlertDialogAction>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>
      </div>

      {!data?.sessions?.length ? (
        <Card>
          <CardContent className="flex items-center justify-center h-32">
            <p className="text-muted-foreground">No active sessions found.</p>
          </CardContent>
        </Card>
      ) : (
        <div className="grid gap-4">
          {data.sessions.map((session) => (
            <Card key={session.token}>
              <CardHeader>
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    {getDeviceIcon(session.user_agent)}
                    <div>
                      <CardTitle className="text-base">
                        {session.device_info.device_name || 'Unknown Device'}
                      </CardTitle>
                      <CardDescription>
                        {getLocation(session.device_info.ip_address)}
                      </CardDescription>
                    </div>
                  </div>
                  <div className="flex items-center gap-2">
                    {session.token === accessToken ? (
                      <Badge variant="secondary">
                        <Shield className="w-3 h-3 mr-1" />
                        Current Session
                      </Badge>
                    ) : (
                      <Badge variant="outline">Inactive</Badge>
                    )}
                    <AlertDialog>
                      <AlertDialogTrigger asChild>
                        <Button
                          variant="ghost"
                          size="sm"
                          disabled={session.token === accessToken}
                        >
                          <Trash2 className="h-4 w-4" />
                        </Button>
                      </AlertDialogTrigger>
                      <AlertDialogContent>
                        <AlertDialogHeader>
                          <AlertDialogTitle>Revoke Session</AlertDialogTitle>
                          <AlertDialogDescription>
                            This will sign out this device. You'll need to sign
                            in again on this device.
                          </AlertDialogDescription>
                        </AlertDialogHeader>
                        <AlertDialogFooter>
                          <AlertDialogCancel>Cancel</AlertDialogCancel>
                          <AlertDialogAction
                            onClick={() =>
                              revokeSessionMutation.mutate(session.token)
                            }
                            disabled={revokeSessionMutation.isPending}
                          >
                            {revokeSessionMutation.isPending
                              ? 'Revoking...'
                              : 'Revoke'}
                          </AlertDialogAction>
                        </AlertDialogFooter>
                      </AlertDialogContent>
                    </AlertDialog>
                  </div>
                </div>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-2 gap-4 text-sm">
                  <div>
                    <span className="font-medium">Started:</span>
                    <p className="text-muted-foreground">
                      {new Date(session.created_at).toLocaleString()}
                    </p>
                  </div>
                  <div>
                    <span className="font-medium">Last Active:</span>
                    <p className="text-muted-foreground">
                      {new Date(session.last_active).toLocaleString()}
                    </p>
                  </div>
                  <div className="col-span-2">
                    <span className="font-medium">User Agent:</span>
                    <p className="text-muted-foreground text-xs mt-1 break-all">
                      {session.device_info.user_agent}
                    </p>
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  )
}
