import oauthService from '@/services/oauth.service'
import type { OAuthProvider } from '@/types/models/oauth_account'

export class OAuthManager {
	private static instance: OAuthManager
	private oauthWindows: Map<string, Window | null> = new Map()

	private constructor() {}

	static getInstance(): OAuthManager {
		if (!OAuthManager.instance) {
			OAuthManager.instance = new OAuthManager()
		}
		return OAuthManager.instance
	}

	/**
	 * Initiate OAuth login flow
	 */
	async loginWithOAuth(provider: OAuthProvider): Promise<void> {
		try {
			const { auth_url, state } =
				await oauthService.initiateOAuthLogin(provider)

			// Store state for validation
			sessionStorage.setItem(`oauth_state_${provider}`, state)

			// Add popup parameter to the auth URL
			const url = new URL(auth_url)
			url.searchParams.set('popup', 'true')

			// Open OAuth window
			const oauthWindow = this.openOAuthWindow(url.toString(), provider)
			this.oauthWindows.set(provider, oauthWindow)

			// Listen for callback
			return this.waitForCallback(provider, oauthWindow)
		} catch (error) {
			console.error('OAuth login failed:', error)
			throw error
		}
	}

	/**
	 * Link OAuth account to existing user
	 */
	async linkAccount(provider: OAuthProvider): Promise<void> {
		try {
			const { auth_url, state } = await oauthService.linkOAuthAccount(provider)

			// Store state for validation
			sessionStorage.setItem(`oauth_state_${provider}`, state)

			// Add popup parameter to the auth URL
			const url = new URL(auth_url)
			url.searchParams.set('popup', 'true')

			// Open OAuth window
			const oauthWindow = this.openOAuthWindow(url.toString(), provider)
			this.oauthWindows.set(provider, oauthWindow)

			// Listen for callback
			return this.waitForLinkCallback(provider, oauthWindow)
		} catch (error) {
			console.error('OAuth link failed:', error)
			throw error
		}
	}

	private openOAuthWindow(authUrl: string, provider: string): Window | null {
		const width = 600
		const height = 700
		const left = window.screenX + (window.outerWidth - width) / 2
		const top = window.screenY + (window.outerHeight - height) / 2

		return window.open(
			authUrl,
			`oauth_${provider}`,
			`width=${width},height=${height},left=${left},top=${top},resizable=yes,scrollbars=yes`,
		)
	}

	private waitForCallback(
		provider: OAuthProvider,
		oauthWindow: Window | null,
	): Promise<void> {
		return new Promise((resolve, reject) => {
			if (!oauthWindow) {
				reject(new Error('Failed to open OAuth window'))
				return
			}

			const checkClosed = setInterval(() => {
				if (oauthWindow.closed) {
					clearInterval(checkClosed)
					window.removeEventListener('message', handleMessage)
					reject(new Error('OAuth window was closed'))
				}
			}, 1000)

			const handleMessage = async (event: MessageEvent) => {
				// Verify origin for security
				if (event.origin !== window.location.origin) {
					return
				}

				if (event.data.type === 'OAUTH_SUCCESS') {
					clearInterval(checkClosed)
					window.removeEventListener('message', handleMessage)
					oauthWindow.close()
					this.oauthWindows.delete(provider)

					try {
						const { authData } = event.data

						// Store tokens and user data
						const authStore = (
							await import('@/stores/auth')
						).useAuthStore.getState()
						authStore.setUser(authData.user)
						authStore.setTokens(
							authData.access_token,
							authData.refresh_token,
							authData.access_token_expires_at,
							authData.refresh_token_expires_at,
						)

						// Clean up
						sessionStorage.removeItem(`oauth_state_${provider}`)

						resolve()
					} catch (error) {
						reject(error)
					}
				} else if (event.data.type === 'OAUTH_ERROR') {
					clearInterval(checkClosed)
					window.removeEventListener('message', handleMessage)
					oauthWindow.close()
					this.oauthWindows.delete(provider)
					reject(new Error(event.data.error || 'OAuth authentication failed'))
				}
			}

			window.addEventListener('message', handleMessage)
		})
	}

	private waitForLinkCallback(
		provider: OAuthProvider,
		oauthWindow: Window | null,
	): Promise<void> {
		return new Promise((resolve, reject) => {
			if (!oauthWindow) {
				reject(new Error('Failed to open OAuth window'))
				return
			}

			const checkClosed = setInterval(() => {
				if (oauthWindow.closed) {
					clearInterval(checkClosed)
					window.removeEventListener('message', handleMessage)
					reject(new Error('OAuth window was closed'))
				}
			}, 1000)

			const handleMessage = (event: MessageEvent) => {
				// Verify origin for security
				if (event.origin !== window.location.origin) {
					return
				}

				if (event.data.type === 'OAUTH_LINK_SUCCESS') {
					clearInterval(checkClosed)
					window.removeEventListener('message', handleMessage)
					oauthWindow.close()
					this.oauthWindows.delete(provider)

					// Clean up
					sessionStorage.removeItem(`oauth_state_${provider}`)

					resolve()
				} else if (event.data.type === 'OAUTH_ERROR') {
					clearInterval(checkClosed)
					window.removeEventListener('message', handleMessage)
					oauthWindow.close()
					this.oauthWindows.delete(provider)
					reject(new Error(event.data.error || 'OAuth linking failed'))
				}
			}

			window.addEventListener('message', handleMessage)
		})
	}

	/**
	 * Clean up OAuth windows and listeners
	 */
	cleanup(): void {
		this.oauthWindows.forEach((window) => {
			if (window && !window.closed) {
				window.close()
			}
		})
		this.oauthWindows.clear()
	}
}

export const oauthManager = OAuthManager.getInstance()
