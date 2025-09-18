import type { User } from '../models/user'
import type { UserDevice } from '../models/user_device'

export type RegisterResponse = {
	user: User
	access_token: string
	access_token_expires_at: string
	refresh_token: string
	refresh_token_expires_at: string
	device: UserDevice | null
}

export type LoginResponse = {
	user: User
	access_token: string
	access_token_expires_at: string
	refresh_token: string
	refresh_token_expires_at: string
	device: UserDevice | null
}

export type RefreshTokenResponse = {
	access_token: string
	expires_at: string
}

export type LogoutResponse = {
	message: string
}
export type VerifyResetTokenResponse = {
	valid: boolean
	expires_at: string
}
