import { createFileRoute } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Shield, Activity, Monitor, AlertTriangle, Download, Users, Clock, CheckCircle } from 'lucide-react'
import { Link } from '@tanstack/react-router'
import authService from '@/services/auth.service'
import sessionService from '@/services/session.service'
import devicesService from '@/services/devices.service'
import securityService from '@/services/security.service'
import auditService from '@/services/audit.service'

export const Route = createFileRoute('/_pathlessLayout/')({
  component: RouteComponent,
})

function RouteComponent() {
  // Get dashboard data
  const { data: user } = useQuery({
    queryKey: ['current-user'],
    queryFn: authService.getCurrentUser,
  })

  const { data: sessions } = useQuery({
    queryKey: ['user-sessions'],
    queryFn: sessionService.getUserSessions,
  })

  const { data: devices } = useQuery({
    queryKey: ['user-devices'],
    queryFn: devicesService.getUserDevices,
  })

  const { data: suspiciousActivities } = useQuery({
    queryKey: ['suspicious-activities'],
    queryFn: securityService.getSuspiciousActivities,
  })

  const { data: recentAuditLogs } = useQuery({
    queryKey: ['recent-audit-logs'],
    queryFn: auditService.getRecentAuditLogs,
  })

  const activeSessions = sessions?.filter(s => s.is_active).length || 0
  const trustedDevices = devices?.filter(d => d.trusted).length || 0
  const unresolvedActivities = suspiciousActivities?.filter(a => !a.resolved).length || 0

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-3xl font-bold tracking-tight">Dashboard</h2>
        <p className="text-muted-foreground">
          Welcome back! Here's an overview of your account security and activity.
        </p>
      </div>

      {/* Security Overview Cards */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Active Sessions</CardTitle>
            <Monitor className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{activeSessions}</div>
            <p className="text-xs text-muted-foreground">
              Across all your devices
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Trusted Devices</CardTitle>
            <Shield className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{trustedDevices}</div>
            <p className="text-xs text-muted-foreground">
              Verified for quick access
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Security Alerts</CardTitle>
            <AlertTriangle className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{unresolvedActivities}</div>
            <p className="text-xs text-muted-foreground">
              {unresolvedActivities === 0 ? 'All clear!' : 'Require attention'}
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Account Status</CardTitle>
            <Users className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="flex items-center gap-2">
              {user?.active ? (
                <>
                  <CheckCircle className="h-5 w-5 text-green-500" />
                  <span className="text-sm font-medium">Active</span>
                </>
              ) : (
                <>
                  <AlertTriangle className="h-5 w-5 text-red-500" />
                  <span className="text-sm font-medium">Inactive</span>
                </>
              )}
            </div>
            <p className="text-xs text-muted-foreground">
              {user?.email_verified ? 'Email verified' : 'Email pending verification'}
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Quick Actions */}
      <Card>
        <CardHeader>
          <CardTitle>Quick Actions</CardTitle>
          <CardDescription>
            Common tasks and account management options.
          </CardDescription>
        </CardHeader>
        <CardContent className="grid gap-2 md:grid-cols-2 lg:grid-cols-4">
          <Button asChild variant="outline" className="justify-start">
            <Link to="/me">
              <Users className="mr-2 h-4 w-4" />
              Edit Profile
            </Link>
          </Button>
          <Button asChild variant="outline" className="justify-start">
            <Link to="/sessions">
              <Monitor className="mr-2 h-4 w-4" />
              Manage Sessions
            </Link>
          </Button>
          <Button asChild variant="outline" className="justify-start">
            <Link to="/devices">
              <Shield className="mr-2 h-4 w-4" />
              Trusted Devices
            </Link>
          </Button>
          <Button asChild variant="outline" className="justify-start">
            <Link to="/security">
              <Download className="mr-2 h-4 w-4" />
              Export Data
            </Link>
          </Button>
        </CardContent>
      </Card>

      {/* Recent Activity */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle>Recent Activity</CardTitle>
              <CardDescription>
                Your latest account activity and security events.
              </CardDescription>
            </div>
            <Button asChild variant="outline" size="sm">
              <Link to="/audit-logs">
                View All
              </Link>
            </Button>
          </div>
        </CardHeader>
        <CardContent>
          {!recentAuditLogs?.length ? (
            <div className="flex items-center justify-center h-24">
              <p className="text-muted-foreground">No recent activity.</p>
            </div>
          ) : (
            <div className="space-y-3">
              {recentAuditLogs.slice(0, 5).map((log) => (
                <div key={log.id} className="flex items-center gap-3 p-3 border rounded-lg">
                  <Activity className="h-4 w-4 text-muted-foreground" />
                  <div className="flex-1">
                    <div className="flex items-center gap-2 mb-1">
                      <span className="text-sm font-medium">
                        {log.action.split('_').map(word =>
                          word.charAt(0).toUpperCase() + word.slice(1)
                        ).join(' ')}
                      </span>
                      {log.resource_type && (
                        <Badge variant="outline" className="text-xs">
                          {log.resource_type}
                        </Badge>
                      )}
                    </div>
                    <div className="flex items-center gap-4 text-xs text-muted-foreground">
                      <span>
                        <Clock className="inline h-3 w-3 mr-1" />
                        {log.created_at
                          ? new Date(log.created_at).toLocaleString()
                          : 'Unknown'
                        }
                      </span>
                      {log.ip_address && (
                        <span>IP: {log.ip_address}</span>
                      )}
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      {/* Security Alerts */}
      {unresolvedActivities > 0 && (
        <Card className="border-red-200 bg-red-50">
          <CardHeader>
            <div className="flex items-center gap-2">
              <AlertTriangle className="h-5 w-5 text-red-500" />
              <CardTitle className="text-red-700">Security Alerts</CardTitle>
            </div>
            <CardDescription className="text-red-600">
              You have {unresolvedActivities} unresolved security alert{unresolvedActivities !== 1 ? 's' : ''}.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Button asChild>
              <Link to="/security">
                Review Security Alerts
              </Link>
            </Button>
          </CardContent>
        </Card>
      )}
    </div>
  )
}
