import type { AuthStore } from '@/types/stores/auth'
import { create } from 'zustand'

export const useAuthStore = create<AuthStore>()((set, get) => ({
	accessToken: null,
	refreshToken: null,
	accessTokenExpiresAt: null,
	refreshTokenExpiresAt: null,
	User: null,
	register: (data) => {
		get().setUser(data.user)
		get().setAccessToken(data.access_token, data.access_token_expires_at)
		get().setRefreshToken(data.refresh_token, data.refresh_token_expires_at)
	},
	login: (data) => {
		get().setUser(data.user)
		get().setAccessToken(data.access_token, data.access_token_expires_at)
		get().setRefreshToken(data.refresh_token, data.refresh_token_expires_at)
	},
	logout: () => {
		set({
			accessToken: null,
			refreshToken: null,
			accessTokenExpiresAt: null,
			refreshTokenExpiresAt: null,
			User: null,
		})
	},
	checkAuth: () => {
		const accessToken = get().accessToken
		const accessTokenExpiresAt = get().accessTokenExpiresAt
		const refreshToken = get().refreshToken
		const refreshTokenExpiresAt = get().refreshTokenExpiresAt
		if (
			!accessToken ||
			!accessTokenExpiresAt ||
			!refreshToken ||
			!refreshTokenExpiresAt
		) {
			return false
		}
		if (refreshTokenExpiresAt.getTime() < Date.now()) {
			return false
		}
		return true
	},
	setUser: (newUser) => {
		set({ User: newUser })
	},
	setAccessToken: (token, expiresAt) => {
		const expiresAtDate = new Date(expiresAt)
		set({ accessToken: token, accessTokenExpiresAt: expiresAtDate })
	},
	setRefreshToken: (token, expiresAt) => {
		const expiresAtDate = new Date(expiresAt)
		set({ refreshToken: token, refreshTokenExpiresAt: expiresAtDate })
	},
}))
