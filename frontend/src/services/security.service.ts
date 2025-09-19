import config from '@/lib/config'
import { apiRequest } from './api.service'
import type {
	SuspiciousActivity,
	ResolveSuspiciousActivityRequest,
} from '@/types/models/suspicious_activity'
import type { GenericMessageResponse } from '@/types/api/generic.response'

const securityAPIUrl = `${config.apiUrl}/security`

const securityService = {
	getSuspiciousActivities: () =>
		apiRequest<SuspiciousActivity[]>({
			headers: undefined,
			protected: true,
			method: 'GET',
			params: undefined,
			url: `${securityAPIUrl}/activities`,
		}),

	resolveSuspiciousActivity: (data: ResolveSuspiciousActivityRequest) =>
		apiRequest<GenericMessageResponse>({
			headers: undefined,
			protected: true,
			method: 'POST',
			params: undefined,
			url: `${securityAPIUrl}/activities/resolve`,
			data,
		}),

	cleanupExpiredLockouts: () =>
		apiRequest<GenericMessageResponse>({
			headers: undefined,
			protected: true,
			method: 'POST',
			params: undefined,
			url: `${securityAPIUrl}/cleanup`,
		}),
}

export default securityService
