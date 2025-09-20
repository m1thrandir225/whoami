import * as React from 'react'

import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from '@/components/ui/sidebar'
import { Link } from '@tanstack/react-router'
import { NavUser } from './nav-user'
import { Logo } from './ui/logo'

// This is sample data.
const data = {
  navMain: [
    {
      title: 'Me',
      url: '/me',
    },
    {
      title: 'Devices',
      url: '/devices',
    },
    {
      title: 'Audit Logs',
      url: '/audit-logs',
    },
    {
      title: 'Sessions',
      url: '/sessions',
    },
    {
      title: 'Security',
      url: '/security',
    },
  ],
}

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  return (
    <Sidebar variant="floating" {...props}>
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem className="w-full block">
            <Link
              to="/"
              className="w-full self-center flex items-center justify-center py-4 hover:opacity-60"
            >
              <Logo variant="full" className="self-center" />
            </Link>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent>
        <SidebarGroup>
          <SidebarMenu className="gap-2">
            {data.navMain.map((item) => (
              <SidebarMenuItem key={item.title}>
                <SidebarMenuButton asChild>
                  <Link
                    to={item.url}
                    className="font-medium transition-all ease-out duration-100"
                    activeProps={{
                      className:
                        'bg-sidebar-accent text-sidebar-accent-foreground border !font-semibold ',
                    }}
                  >
                    {item.title}
                  </Link>
                </SidebarMenuButton>
              </SidebarMenuItem>
            ))}
          </SidebarMenu>
        </SidebarGroup>
      </SidebarContent>
      <SidebarFooter>
        <NavUser />
      </SidebarFooter>
    </Sidebar>
  )
}
