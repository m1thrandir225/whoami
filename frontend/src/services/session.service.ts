import config from '@/lib/config'
import { apiRequest } from './api.service'
import type { Session } from '@/types/models/session'
import type { GenericMessageResponse } from '@/types/api/generic.response'

const sessionAPIUrl = `${config.apiUrl}/sessions`

const sessionService = {
	getUserSessions: () =>
		apiRequest<Session[]>({
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
