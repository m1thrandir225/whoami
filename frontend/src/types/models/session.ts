import type { UserDevice } from './user_device'

export type SessionResponse = {
	sessions: Session[]
}

export type Session = {
	user_id: number
	token: string
	device_info: UserDevice
	ip_address: string
	user_agent: string
	created_at: string
	last_active: string
	is_active: boolean
}
