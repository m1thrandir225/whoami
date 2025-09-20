import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { Monitor, Smartphone, Tablet, Shield, Trash2 } from 'lucide-react'
import { toast } from 'sonner'

import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import sessionService from '@/services/session.service'
import { useAuthStore } from '@/stores/auth'
import { PageLoader } from '@/components/ui/loader'
import { createFileRoute } from '@tanstack/react-router'
import type { UserDevice } from '@/types/models/user_device'

export const Route = createFileRoute('/_pathlessLayout/sessions')({
  component: RouteComponent,
})

function RouteComponent() {
  const queryClient = useQueryClient()
  const authStore = useAuthStore()
  const currentAccessToken = authStore.accessToken

  const { data, isLoading } = useQuery({
    queryKey: ['user-sessions'],
    queryFn: () => sessionService.getUserSessions(),
  })

  const revokeSessionMutation = useMutation({
    mutationFn: (sessionId: string) => sessionService.revokeSession(sessionId),
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

  const getDeviceName = (deviceInfo: UserDevice) => {
    return (
      deviceInfo.device_name || `${deviceInfo.device_type || 'Unknown'} Device`
    )
  }

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString()
  }

  const isCurrentSession = (session: any) => {
    return session.token === currentAccessToken
  }

  if (isLoading) {
    return <PageLoader text="Loading sessions..." />
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">Active Sessions</h1>
          <p className="text-muted-foreground">
            Manage your active login sessions across devices
          </p>
        </div>
        <Button
          variant="destructive"
          onClick={() => revokeAllSessionsMutation.mutate()}
          disabled={revokeAllSessionsMutation.isPending}
        >
          <Trash2 className="w-4 h-4 mr-2" />
          Revoke All Sessions
        </Button>
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
            <Card key={session.id}>
              <CardHeader>
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    {getDeviceIcon(session.user_agent)}
                    <div>
                      <CardTitle className="text-base">
                        {getDeviceName(session.device_info)}
                      </CardTitle>
                      <CardDescription>
                        {getLocation(session.ip_address)}
                      </CardDescription>
                    </div>
                  </div>
                  <div className="flex items-center gap-2">
                    {isCurrentSession(session) ? (
                      <Badge variant="secondary">
                        <Shield className="w-3 h-3 mr-1" />
                        Current Session
                      </Badge>
                    ) : (
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => revokeSessionMutation.mutate(session.id)}
                        disabled={revokeSessionMutation.isPending}
                      >
                        <Trash2 className="w-3 h-3 mr-1" />
                        Revoke
                      </Button>
                    )}
                  </div>
                </div>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-2 gap-4 text-sm">
                  <div>
                    <p className="text-muted-foreground">Created</p>
                    <p>{formatDate(session.created_at)}</p>
                  </div>
                  <div>
                    <p className="text-muted-foreground">Last Active</p>
                    <p>{formatDate(session.last_active)}</p>
                  </div>
                  <div>
                    <p className="text-muted-foreground">Device Type</p>
                    <p className="capitalize">
                      {session.device_info.device_type || 'Unknown'}
                    </p>
                  </div>
                  <div>
                    <p className="text-muted-foreground">Browser</p>
                    <p>
                      {session.device_info.user_agent?.split(' ')[0] ||
                        'Unknown'}
                    </p>
                  </div>
                </div>

                {/* Debug info - remove in production */}
                <details className="mt-4">
                  <summary className="text-xs text-muted-foreground cursor-pointer">
                    Session Details
                  </summary>
                  <div className="mt-2 text-xs bg-muted p-2 rounded">
                    <p>
                      <strong>Session ID:</strong> {session.id}
                    </p>
                    <p>
                      <strong>Access Token:</strong>{' '}
                      {session.token.slice(0, 20)}...
                    </p>
                    <p>
                      <strong>Refresh Token:</strong>{' '}
                      {session.refresh_token?.slice(0, 20)}...
                    </p>
                  </div>
                </details>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  )
}
