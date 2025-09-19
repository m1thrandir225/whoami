import type { LoginResponse, RegisterResponse } from '../api/auth.responses'
import type { User } from '../models/user'

export type AuthStore = State & Actions

type State = {
	accessToken: string | null
	accessTokenExpiresAt: Date | null
	refreshToken: string | null
	refreshTokenExpiresAt: Date | null
	User: User | null
}

type Actions = {
	isAuthenticated: () => boolean
	register: (data: RegisterResponse) => void
	setTokens: (
		accessToken: string,
		refreshToken: string,
		accessTokenExpiresAt: string,
		refreshTokenExpiresAt: string,
	) => void
	login: (data: LoginResponse) => void
	logout: () => void
	checkAuth: () => boolean
	setUser: (user: User) => void
	setAccessToken: (token: string, expiresAt: string) => void
	setRefreshToken: (token: string, expiresAt: string) => void
}
