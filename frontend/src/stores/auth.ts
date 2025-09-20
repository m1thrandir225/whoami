import type { AuthStore } from '@/types/stores/auth'
import { create } from 'zustand'
import { persist } from 'zustand/middleware'
export const useAuthStore = create<AuthStore>()(
	persist(
		(set, get) => ({
			accessToken: null,
			refreshToken: null,
			accessTokenExpiresAt: null,
			refreshTokenExpiresAt: null,
			User: null,
			isAuthenticated: () => {
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
				return true
			},
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
			setTokens: (
				accessToken,
				refreshToken,
				accessTokenExpiresAt,
				refreshTokenExpiresAt,
			) => {
				const accessTokenExpiresAtDate = new Date(accessTokenExpiresAt)
				const refreshTokenExpiresAtDate = new Date(refreshTokenExpiresAt)
				set({
					accessToken,
					refreshToken,
					accessTokenExpiresAt: accessTokenExpiresAtDate,
					refreshTokenExpiresAt: refreshTokenExpiresAtDate,
				})
			},
			setAccessToken: (token, expiresAt) => {
				const expiresAtDate = new Date(expiresAt)
				set({ accessToken: token, accessTokenExpiresAt: expiresAtDate })
			},
			setRefreshToken: (token, expiresAt) => {
				const expiresAtDate = new Date(expiresAt)
				set({ refreshToken: token, refreshTokenExpiresAt: expiresAtDate })
			},
		}),
		{
			name: 'auth',
			onRehydrateStorage: () => (state) => {
				if (state) {
					if (
						state.accessTokenExpiresAt &&
						typeof state.accessTokenExpiresAt === 'string'
					) {
						state.accessTokenExpiresAt = new Date(state.accessTokenExpiresAt)
					}

					if (
						state.refreshTokenExpiresAt &&
						typeof state.refreshTokenExpiresAt === 'string'
					) {
						state.refreshTokenExpiresAt = new Date(state.refreshTokenExpiresAt)
					}
				}
			},
		},
	),
)
