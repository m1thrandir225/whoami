import config from '@/lib/config'
import type { GenericMessageResponse } from '@/types/api/generic.response'
import type { SessionResponse } from '@/types/models/session'
import { apiRequest } from './api.service'

const sessionAPIUrl = `${config.apiUrl}/sessions`

const sessionService = {
	getUserSessions: () =>
		apiRequest<SessionResponse>({
			headers: undefined,
			protected: true,
			method: 'GET',
			params: undefined,
			url: sessionAPIUrl,
		}),

	revokeSession: (token: string) =>
		apiRequest<GenericMessageResponse>({
			headers: undefined,
			protected: true,
			method: 'DELETE',
			params: undefined,
			url: `${sessionAPIUrl}/${token}`,
		}),

	revokeAllSessions: () =>
		apiRequest<GenericMessageResponse>({
			headers: undefined,
			protected: true,
			method: 'DELETE',
			params: undefined,
			url: sessionAPIUrl,
			data: { reason: 'User requested revocation of all sessions' },
		}),
}

export default sessionService
