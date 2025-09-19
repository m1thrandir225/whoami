export type DataExport = {
	id: number
	user_id: number
	export_type: string
	status: 'pending' | 'completed' | 'failed' | 'expired'
	file_path?: string
	file_size?: number
	expires_at: string
	created_at: string | null
	completed_at?: string | null
}

export const DataExportTypes = {
	USER_DATA: 'user_data',
	AUDIT_LOGS: 'audit_logs',
	LOGIN_HISTORY: 'login_history',
	COMPLETE: 'complete',
} as const

export const DataExportStatuses = {
	PENDING: 'pending',
	COMPLETED: 'completed',
	FAILED: 'failed',
	EXPIRED: 'expired',
} as const
