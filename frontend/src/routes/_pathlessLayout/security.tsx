import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import exportsService from '@/services/exports.service'
import securityService from '@/services/security.service'
import { DataExportTypes } from '@/types/models/data_export'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { createFileRoute } from '@tanstack/react-router'
import {
  AlertTriangle,
  CheckCircle,
  Clock,
  Download,
  FileText,
  Shield,
  ShieldAlert,
} from 'lucide-react'
import { useState } from 'react'
import { toast } from 'sonner'

export const Route = createFileRoute('/_pathlessLayout/security')({
  component: RouteComponent,
})

function RouteComponent() {
  const queryClient = useQueryClient()
  const [selectedExportType, setSelectedExportType] = useState<string>('')

  // Get suspicious activities
  const { data: suspiciousActivities, isLoading: activitiesLoading } = useQuery(
    {
      queryKey: ['suspicious-activities'],
      queryFn: securityService.getSuspiciousActivities,
    },
  )

  // Get data exports
  const { data: dataExports, isLoading: exportsLoading } = useQuery({
    queryKey: ['data-exports'],
    queryFn: exportsService.getDataExports,
  })

  // Resolve suspicious activity mutation
  const resolveActivityMutation = useMutation({
    mutationFn: (activityId: number) =>
      securityService.resolveSuspiciousActivity({ activity_id: activityId }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['suspicious-activities'] })
      toast.success('Activity resolved successfully')
    },
    onError: (error) => {
      toast.error(error.message || 'Failed to resolve activity')
    },
  })

  // Request data export mutation
  const requestExportMutation = useMutation({
    mutationFn: (exportType: string) =>
      exportsService.requestDataExport(exportType),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['data-exports'] })
      toast.success('Data export requested successfully')
      setSelectedExportType('')
    },
    onError: (error) => {
      toast.error(error.message || 'Failed to request data export')
    },
  })

  // Download export mutation
  const downloadExportMutation = useMutation({
    mutationFn: (id: number) => exportsService.downloadDataExport(id),
    onSuccess: (blob, variables) => {
      const url = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `export-${variables}.json`
      document.body.appendChild(a)
      a.click()
      window.URL.revokeObjectURL(url)
      document.body.removeChild(a)
      toast.success('Export downloaded successfully')
    },
    onError: (error) => {
      toast.error(error.message || 'Failed to download export')
    },
  })

  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case 'critical':
        return 'destructive'
      case 'high':
        return 'destructive'
      case 'medium':
        return 'secondary'
      case 'low':
        return 'outline'
      default:
        return 'outline'
    }
  }

  const getSeverityIcon = (severity: string) => {
    switch (severity) {
      case 'critical':
      case 'high':
        return <ShieldAlert className="h-4 w-4" />
      default:
        return <AlertTriangle className="h-4 w-4" />
    }
  }

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'completed':
        return <CheckCircle className="h-4 w-4 text-green-500" />
      case 'pending':
        return <Clock className="h-4 w-4 text-yellow-500" />
      case 'failed':
        return <AlertTriangle className="h-4 w-4 text-red-500" />
      default:
        return <FileText className="h-4 w-4" />
    }
  }

  const formatExportType = (type: string) => {
    return type
      .split('_')
      .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
      .join(' ')
  }

  const formatFileSize = (bytes: number) => {
    const sizes = ['Bytes', 'KB', 'MB', 'GB']
    if (bytes === 0) return '0 Byte'
    const i = Math.floor(Math.log(bytes) / Math.log(1024))
    return Math.round((bytes / Math.pow(1024, i)) * 100) / 100 + ' ' + sizes[i]
  }

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-3xl font-bold tracking-tight">Security Center</h2>
        <p className="text-muted-foreground">
          Monitor security events and manage your data exports.
        </p>
      </div>

      {/* Suspicious Activities */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Shield className="h-5 w-5" />
            Suspicious Activities
          </CardTitle>
          <CardDescription>
            Review and resolve security alerts for your account.
          </CardDescription>
        </CardHeader>
        <CardContent>
          {activitiesLoading ? (
            <div>Loading activities...</div>
          ) : !suspiciousActivities?.length ? (
            <div className="flex items-center justify-center h-24">
              <p className="text-muted-foreground">
                No suspicious activities detected.
              </p>
            </div>
          ) : (
            <div className="space-y-4">
              {suspiciousActivities.map((activity) => (
                <div
                  key={activity.id}
                  className="flex items-start justify-between p-4 border rounded-lg"
                >
                  <div className="flex items-start gap-3">
                    {getSeverityIcon(activity.severity)}
                    <div className="flex-1">
                      <div className="flex items-center gap-2 mb-2">
                        <Badge
                          variant={getSeverityColor(activity.severity) as any}
                        >
                          {activity.severity.toUpperCase()}
                        </Badge>
                        <span className="font-medium">
                          {activity.activity_type}
                        </span>
                        {activity.resolved && (
                          <Badge variant="secondary">
                            <CheckCircle className="w-3 h-3 mr-1" />
                            Resolved
                          </Badge>
                        )}
                      </div>
                      <p className="text-sm text-muted-foreground mb-2">
                        {activity.description}
                      </p>
                      <div className="grid grid-cols-2 gap-4 text-xs text-muted-foreground">
                        <div>
                          <span className="font-medium">IP:</span>{' '}
                          {activity.ip_address}
                        </div>
                        <div>
                          <span className="font-medium">Time:</span>{' '}
                          {activity.created_at
                            ? new Date(activity.created_at).toLocaleString()
                            : 'Unknown'}
                        </div>
                      </div>
                    </div>
                  </div>
                  {!activity.resolved && (
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() =>
                        resolveActivityMutation.mutate(activity.id)
                      }
                      disabled={resolveActivityMutation.isPending}
                    >
                      {resolveActivityMutation.isPending
                        ? 'Resolving...'
                        : 'Resolve'}
                    </Button>
                  )}
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      {/* Data Exports */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle className="flex items-center gap-2">
                <Download className="h-5 w-5" />
                Data Exports
              </CardTitle>
              <CardDescription>
                Request and download copies of your account data.
              </CardDescription>
            </div>
            <div className="flex gap-2">
              <Select
                value={selectedExportType}
                onValueChange={setSelectedExportType}
              >
                <SelectTrigger className="w-48">
                  <SelectValue placeholder="Select export type" />
                </SelectTrigger>
                <SelectContent>
                  {Object.values(DataExportTypes).map((type) => (
                    <SelectItem key={type} value={type}>
                      {formatExportType(type)}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              <Button
                onClick={() =>
                  selectedExportType &&
                  requestExportMutation.mutate(selectedExportType)
                }
                disabled={
                  !selectedExportType || requestExportMutation.isPending
                }
              >
                {requestExportMutation.isPending
                  ? 'Requesting...'
                  : 'Request Export'}
              </Button>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          {exportsLoading ? (
            <div>Loading exports...</div>
          ) : !dataExports?.length ? (
            <div className="flex items-center justify-center h-24">
              <p className="text-muted-foreground">No data exports found.</p>
            </div>
          ) : (
            <div className="space-y-4">
              {dataExports.map((exportItem) => (
                <div
                  key={exportItem.id}
                  className="flex items-center justify-between p-4 border rounded-lg"
                >
                  <div className="flex items-center gap-3">
                    {getStatusIcon(exportItem.status)}
                    <div>
                      <div className="flex items-center gap-2 mb-1">
                        <span className="font-medium">
                          {formatExportType(exportItem.export_type)}
                        </span>
                        <Badge variant="outline">{exportItem.status}</Badge>
                      </div>
                      <div className="text-sm text-muted-foreground">
                        <span>
                          Created:{' '}
                          {new Date(
                            exportItem.created_at || '',
                          ).toLocaleString()}
                        </span>
                        {exportItem.file_size && (
                          <span className="ml-4">
                            Size: {formatFileSize(exportItem.file_size)}
                          </span>
                        )}
                        <span className="ml-4">
                          Expires:{' '}
                          {new Date(exportItem.expires_at).toLocaleDateString()}
                        </span>
                      </div>
                    </div>
                  </div>
                  {exportItem.status === 'completed' && (
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() =>
                        downloadExportMutation.mutate(exportItem.id)
                      }
                      disabled={downloadExportMutation.isPending}
                    >
                      <Download className="h-4 w-4 mr-1" />
                      {downloadExportMutation.isPending
                        ? 'Downloading...'
                        : 'Download'}
                    </Button>
                  )}
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
