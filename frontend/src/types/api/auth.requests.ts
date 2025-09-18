import type { PrivacySettings } from '../models/privacy_settings'

export type RegisterRequest = {
	email: string
	password: string
	username?: string
	privacy_settings?: PrivacySettings
}

export type LoginRequest = {
	email: string
	password: string
}

export type UpdateUserRequest = {
	email: string
	username: string
}

export type RefreshTokenRequest = {
	refresh_token: string
}

export type UpdatePasswordRequest = {
	current_password: string
	new_password: string
}

export type ResolveSuspiciousActivityRequest = {
	activity_id: number
}

export type VerifyEmailRequest = {
	token: string
}

export type RequestPasswordResetRequest = {
	email: string
}

export type ResetPasswordRequest = {
	token: string
	new_password: string
}

export type VerifyResetTokenRequest = {
	token: string
}

export type RevokeSessionRequest = {
	token: string
}

export type RevokeAllSessionsRequest = {
	reason: string
}

export type ResetRateLimitRequest = {
	type: string
	ip?: string
	user_id: number
}
