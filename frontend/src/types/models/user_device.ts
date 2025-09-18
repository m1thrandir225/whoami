export type UserDevice = {
	id: number
	user_id: number
	device_id: string
	device_name: string
	device_type: string
	user_agent: string
	ip_address: string
	trusted: boolean
	last_used_at: string
	created_at?: string
}
