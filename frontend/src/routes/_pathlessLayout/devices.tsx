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
import {
  Trash2,
  Shield,
  ShieldCheck,
  Smartphone,
  Monitor,
  Tablet,
} from 'lucide-react'
import { toast } from 'sonner'
import devicesService from '@/services/devices.service'

export const Route = createFileRoute('/_pathlessLayout/devices')({
  component: RouteComponent,
})

function RouteComponent() {
  const queryClient = useQueryClient()

  // Get user devices
  const { data: devices, isLoading } = useQuery({
    queryKey: ['user-devices'],
    queryFn: devicesService.getUserDevices,
  })

  // Trust device mutation
  const trustDeviceMutation = useMutation({
    mutationFn: (id: number) => devicesService.markDeviceAsTrusted(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['user-devices'] })
      toast.success('Device marked as trusted')
    },
    onError: (error) => {
      toast.error(error.message || 'Failed to trust device')
    },
  })

  // Delete device mutation
  const deleteDeviceMutation = useMutation({
    mutationFn: (id: number) => devicesService.deleteUserDevice(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['user-devices'] })
      toast.success('Device removed successfully')
    },
    onError: (error) => {
      toast.error(error.message || 'Failed to remove device')
    },
  })

  // Delete all devices mutation
  const deleteAllDevicesMutation = useMutation({
    mutationFn: devicesService.deleteAllUserDevices,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['user-devices'] })
      toast.success('All devices removed successfully')
    },
    onError: (error) => {
      toast.error(error.message || 'Failed to remove all devices')
    },
  })

  const getDeviceIcon = (deviceType: string) => {
    switch (deviceType.toLowerCase()) {
      case 'mobile':
      case 'smartphone':
        return <Smartphone className="h-4 w-4" />
      case 'tablet':
        return <Tablet className="h-4 w-4" />
      default:
        return <Monitor className="h-4 w-4" />
    }
  }

  if (isLoading) {
    return <div>Loading devices...</div>
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-3xl font-bold tracking-tight">Trusted Devices</h2>
          <p className="text-muted-foreground">
            Manage devices that have accessed your account.
          </p>
        </div>
        <AlertDialog>
          <AlertDialogTrigger asChild>
            <Button variant="destructive" disabled={!devices?.length}>
              Remove All Devices
            </Button>
          </AlertDialogTrigger>
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle>Remove All Devices</AlertDialogTitle>
              <AlertDialogDescription>
                This will remove all devices from your trusted devices list.
                You'll need to verify new devices when signing in.
              </AlertDialogDescription>
            </AlertDialogHeader>
            <AlertDialogFooter>
              <AlertDialogCancel>Cancel</AlertDialogCancel>
              <AlertDialogAction
                onClick={() => deleteAllDevicesMutation.mutate()}
                disabled={deleteAllDevicesMutation.isPending}
              >
                {deleteAllDevicesMutation.isPending
                  ? 'Removing...'
                  : 'Remove All'}
              </AlertDialogAction>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>
      </div>

      {!devices?.length ? (
        <Card>
          <CardContent className="flex items-center justify-center h-32">
            <p className="text-muted-foreground">No devices found.</p>
          </CardContent>
        </Card>
      ) : (
        <div className="grid gap-4">
          {devices.map((device) => (
            <Card key={device.id}>
              <CardHeader>
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    {getDeviceIcon(device.device_type)}
                    <div>
                      <CardTitle className="text-base">
                        {device.device_name || 'Unknown Device'}
                      </CardTitle>
                      <CardDescription>
                        {device.device_type} • IP: {device.ip_address}
                      </CardDescription>
                    </div>
                  </div>
                  <div className="flex items-center gap-2">
                    {device.trusted ? (
                      <Badge variant="secondary">
                        <ShieldCheck className="w-3 h-3 mr-1" />
                        Trusted
                      </Badge>
                    ) : (
                      <div className="flex gap-2">
                        <Badge variant="outline">Untrusted</Badge>
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => trustDeviceMutation.mutate(device.id)}
                          disabled={trustDeviceMutation.isPending}
                        >
                          <Shield className="h-4 w-4 mr-1" />
                          Trust
                        </Button>
                      </div>
                    )}
                    <AlertDialog>
                      <AlertDialogTrigger asChild>
                        <Button variant="ghost" size="sm">
                          <Trash2 className="h-4 w-4" />
                        </Button>
                      </AlertDialogTrigger>
                      <AlertDialogContent>
                        <AlertDialogHeader>
                          <AlertDialogTitle>Remove Device</AlertDialogTitle>
                          <AlertDialogDescription>
                            This will remove the device from your trusted
                            devices list. You'll need to verify this device
                            again when signing in.
                          </AlertDialogDescription>
                        </AlertDialogHeader>
                        <AlertDialogFooter>
                          <AlertDialogCancel>Cancel</AlertDialogCancel>
                          <AlertDialogAction
                            onClick={() =>
                              deleteDeviceMutation.mutate(device.id)
                            }
                            disabled={deleteDeviceMutation.isPending}
                          >
                            {deleteDeviceMutation.isPending
                              ? 'Removing...'
                              : 'Remove'}
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
                    <span className="font-medium">Device ID:</span>
                    <p className="text-muted-foreground font-mono text-xs">
                      {device.device_id}
                    </p>
                  </div>
                  <div>
                    <span className="font-medium">Last Used:</span>
                    <p className="text-muted-foreground">
                      {new Date(device.last_used_at).toLocaleString()}
                    </p>
                  </div>
                  <div className="col-span-2">
                    <span className="font-medium">User Agent:</span>
                    <p className="text-muted-foreground text-xs mt-1 break-all">
                      {device.user_agent}
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
