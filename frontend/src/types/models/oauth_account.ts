export type OAuthAccount = {
	id: number
	user_id: number
	provider: 'google' | 'github' | 'discord' | 'twitter'
	provider_user_id: string
	email?: string
	name?: string
	avatar_url?: string
	token_expires_at?: string
	created_at?: string
	updated_at?: string
}

export const OAuthProviders = {
	GOOGLE: 'google',
	GITHUB: 'github',
	DISCORD: 'discord',
	TWITTER: 'twitter',
} as const

export type OAuthProvider = (typeof OAuthProviders)[keyof typeof OAuthProviders]

export type OAuthAuthResponse = {
	auth_url: string
	state: string
}
