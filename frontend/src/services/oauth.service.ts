import config from '@/lib/config'
import { apiRequest } from './api.service'
import type {
	OAuthAccount,
	OAuthAuthResponse,
} from '@/types/models/oauth_account'
import type { LoginResponse } from '@/types/api/auth.responses'
import type { GenericMessageResponse } from '@/types/api/generic.response'

const oauthAPIUrl = `${config.apiUrl}/oauth`

const oauthService = {
	// Initiate OAuth login flow
	initiateOAuthLogin: (provider: string) =>
		apiRequest<OAuthAuthResponse>({
			headers: undefined,
			protected: false,
			method: 'GET',
			params: undefined,
			url: `${oauthAPIUrl}/login/${provider}`,
		}),

	// Exchange temporary token for auth data
	exchangeTempToken: (token: string) =>
		apiRequest<LoginResponse>({
			headers: undefined,
			protected: false,
			method: 'POST',
			params: undefined,
			url: `${oauthAPIUrl}/exchange`,
			data: { token },
		}),

	// Link OAuth account to existing user
	linkOAuthAccount: (provider: string) =>
		apiRequest<OAuthAuthResponse>({
			headers: undefined,
			protected: true,
			method: 'POST',
			params: undefined,
			url: `${oauthAPIUrl}/link`,
			data: { provider },
		}),

	// Get linked OAuth accounts
	getOAuthAccounts: () =>
		apiRequest<OAuthAccount[]>({
			headers: undefined,
			protected: true,
			method: 'GET',
			params: undefined,
			url: `${oauthAPIUrl}/accounts`,
		}),

	// Unlink OAuth account
	unlinkOAuthAccount: (provider: string) =>
		apiRequest<GenericMessageResponse>({
			headers: undefined,
			protected: true,
			method: 'DELETE',
			params: undefined,
			url: `${oauthAPIUrl}/unlink/${provider}`,
		}),
}

export default oauthService
