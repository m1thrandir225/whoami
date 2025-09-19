export type SuspiciousActivity = {
	id: number
	user_id: number
	activity_type: string
	ip_address: string
	user_agent: string
	description: string
	metadata: Record<string, any> | null
	severity: 'low' | 'medium' | 'high' | 'critical'
	resolved: boolean | null
	created_at: string | null
}

export type ResolveSuspiciousActivityRequest = {
	activity_id: number
}
