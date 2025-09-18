import type { PrivacySettings } from './privacy_settings'

export type User = {
	id: number
	email: string
	username: string
	password: string
	email_verified: boolean
	role: string
	active: boolean
	privacy_settings: PrivacySettings
	last_login_at: string | null
	password_changed_at: string
	updated_at: string
	created_at: string
}
