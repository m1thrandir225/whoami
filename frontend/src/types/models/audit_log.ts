export type AuditLog = {
	id: number
	user_id: number | null
	action: string
	resource_type: string | null
	resource_id: number | null
	ip_address: string | null
	user_agent: string | null
	details: Record<string, any> | null
	created_at: string | null
}

// Common audit actions
export const AuditActions = {
	USER_LOGIN: 'user_login',
	USER_LOGOUT: 'user_logout',
	USER_REGISTER: 'user_register',
	USER_UPDATE: 'user_update',
	USER_DEACTIVATE: 'user_deactivate',
	USER_ACTIVATE: 'user_activate',
	PASSWORD_CHANGE: 'password_change',
	PASSWORD_RESET: 'password_reset',
	EMAIL_VERIFY: 'email_verify',
	EMAIL_RESEND: 'email_resend',
	SESSION_CREATE: 'session_create',
	SESSION_REVOKE: 'session_revoke',
	SESSION_REVOKE_ALL: 'session_revoke_all',
	ACCOUNT_LOCKOUT: 'account_lockout',
	SUSPICIOUS_ACTIVITY: 'suspicious_activity',
	DATA_EXPORT: 'data_export',
	PRIVACY_SETTINGS: 'privacy_settings',
} as const

// Common resource types
export const AuditResourceTypes = {
	USER: 'user',
	SESSION: 'session',
	PASSWORD: 'password',
	EMAIL: 'email',
	ACCOUNT: 'account',
	DATA: 'data',
	PRIVACY: 'privacy',
	DEVICE: 'device',
} as const
