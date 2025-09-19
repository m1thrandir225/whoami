import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import auditService from '@/services/audit.service'
import { AuditActions, AuditResourceTypes } from '@/types/models/audit_log'
import { useQuery } from '@tanstack/react-query'
import { createFileRoute } from '@tanstack/react-router'
import { Activity, Lock, Mail, Shield, User } from 'lucide-react'
import { useState } from 'react'

export const Route = createFileRoute('/_pathlessLayout/audit-logs')({
  component: RouteComponent,
})

function RouteComponent() {
  const [filter, setFilter] = useState<'recent' | 'action' | 'resource'>(
    'recent',
  )
  const [filterValue, setFilterValue] = useState('')

  // Get audit logs based on filter
  const { data: auditLogs, isLoading } = useQuery({
    queryKey: ['audit-logs', filter, filterValue],
    queryFn: () => {
      switch (filter) {
        case 'action':
          return filterValue
            ? auditService.getAuditLogsByAction(filterValue)
            : auditService.getRecentAuditLogs()
        case 'resource':
          return filterValue
            ? auditService.getAuditLogsByResourceType(filterValue)
            : auditService.getRecentAuditLogs()
        default:
          return auditService.getRecentAuditLogs()
      }
    },
  })

  const getActionIcon = (action: string) => {
    if (action.includes('login') || action.includes('logout'))
      return <Shield className="h-4 w-4" />
    if (action.includes('user')) return <User className="h-4 w-4" />
    if (action.includes('password')) return <Lock className="h-4 w-4" />
    if (action.includes('email')) return <Mail className="h-4 w-4" />
    return <Activity className="h-4 w-4" />
  }

  const getActionColor = (action: string) => {
    if (action.includes('login')) return 'default'
    if (action.includes('logout')) return 'secondary'
    if (action.includes('register')) return 'default'
    if (action.includes('password')) return 'destructive'
    if (action.includes('suspicious')) return 'destructive'
    if (action.includes('lockout')) return 'destructive'
    return 'outline'
  }

  const formatAction = (action: string) => {
    return action
      .split('_')
      .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
      .join(' ')
  }

  if (isLoading) {
    return <div>Loading audit logs...</div>
  }

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-3xl font-bold tracking-tight">Audit Logs</h2>
        <p className="text-muted-foreground">
          Track all activity and changes to your account.
        </p>
      </div>

      {/* Filters */}
      <Card>
        <CardHeader>
          <CardTitle className="text-base">Filters</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex gap-4 items-end">
            <div className="space-y-2">
              <label className="text-sm font-medium">Filter by</label>
              <Select
                value={filter}
                onValueChange={(value: any) => setFilter(value)}
              >
                <SelectTrigger className="w-48">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="recent">Recent Activity</SelectItem>
                  <SelectItem value="action">Action Type</SelectItem>
                  <SelectItem value="resource">Resource Type</SelectItem>
                </SelectContent>
              </Select>
            </div>

            {filter === 'action' && (
              <div className="space-y-2">
                <label className="text-sm font-medium">Action</label>
                <Select value={filterValue} onValueChange={setFilterValue}>
                  <SelectTrigger className="w-48">
                    <SelectValue placeholder="Select action" />
                  </SelectTrigger>
                  <SelectContent>
                    {Object.values(AuditActions).map((action) => (
                      <SelectItem key={action} value={action}>
                        {formatAction(action)}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            )}

            {filter === 'resource' && (
              <div className="space-y-2">
                <label className="text-sm font-medium">Resource</label>
                <Select value={filterValue} onValueChange={setFilterValue}>
                  <SelectTrigger className="w-48">
                    <SelectValue placeholder="Select resource" />
                  </SelectTrigger>
                  <SelectContent>
                    {Object.values(AuditResourceTypes).map((resource) => (
                      <SelectItem key={resource} value={resource}>
                        {formatAction(resource)}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            )}

            <Button onClick={() => setFilterValue('')} variant="outline">
              Clear
            </Button>
          </div>
        </CardContent>
      </Card>
      {/* Audit Logs */}
      {!auditLogs?.audit_logs?.length ? (
        <Card>
          <CardContent className="flex items-center justify-center h-32">
            <p className="text-muted-foreground">No audit logs found.</p>
          </CardContent>
        </Card>
      ) : (
        <div className="flex flex-col gap-3 max-w-[1200px] max-h-[500px] overflow-y-auto p-6">
          {auditLogs.audit_logs.map((log) => (
            <Card key={log.id} className="max-w-full">
              <CardContent className="pt-6 w-full">
                <div className="flex items-start justify-between">
                  <div className="flex items-start gap-3">
                    {getActionIcon(log.action)}
                    <div className="flex-1">
                      <div className="flex items-center gap-2 mb-2">
                        <Badge variant={getActionColor(log.action) as any}>
                          {formatAction(log.action)}
                        </Badge>
                        {log.resource_type && (
                          <Badge variant="outline">
                            {formatAction(log.resource_type)}
                          </Badge>
                        )}
                      </div>
                      <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
                        <div>
                          <span className="font-medium">Time:</span>
                          <p className="text-muted-foreground">
                            {log.created_at
                              ? new Date(log.created_at).toLocaleString()
                              : 'Unknown'}
                          </p>
                        </div>
                        {log.ip_address && (
                          <div>
                            <span className="font-medium">IP:</span>
                            <p className="text-muted-foreground font-mono">
                              {log.ip_address}
                            </p>
                          </div>
                        )}
                        {log.resource_id && (
                          <div>
                            <span className="font-medium">Resource ID:</span>
                            <p className="text-muted-foreground">
                              {log.resource_id}
                            </p>
                          </div>
                        )}
                        {log.user_agent && (
                          <div className="col-span-2 md:col-span-1">
                            <span className="font-medium">Device:</span>
                            <p className="text-muted-foreground text-xs truncate">
                              {log.user_agent}
                            </p>
                          </div>
                        )}
                      </div>
                      {log.details && (
                        <div className="mt-3">
                          <span className="font-medium text-sm">Details:</span>
                          <pre className="text-xs text-muted-foreground mt-1 bg-muted p-2 rounded break-all text-wrap">
                            {JSON.stringify(log.details, null, 2)}
                          </pre>
                        </div>
                      )}
                    </div>
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
