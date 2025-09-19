export type Session = {
	user_id: number
	token: string
	device_info: Record<string, string>
	ip_address: string
	user_agent: string
	created_at: string
	last_active: string
	is_active: boolean
}
