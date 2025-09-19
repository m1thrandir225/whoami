import config from '@/lib/config'
import type { GenericMessageResponse } from '@/types/api/generic.response'
import type { AuditLogResponse } from '@/types/models/audit_log'
import { apiRequest } from './api.service'

const auditAPIUrl = `${config.apiUrl}/audit`

const auditService = {
	getAuditLogsByUserID: (userId: number) =>
		apiRequest<AuditLogResponse>({
			headers: undefined,
			protected: true,
			method: 'GET',
			params: undefined,
			url: `${auditAPIUrl}/user/${userId}`,
		}),

	getAuditLogsByAction: (action: string) =>
		apiRequest<AuditLogResponse>({
			headers: undefined,
			protected: true,
			method: 'GET',
			params: undefined,
			url: `${auditAPIUrl}/action/${action}`,
		}),

	getAuditLogsByResourceType: (resourceType: string) =>
		apiRequest<AuditLogResponse>({
			headers: undefined,
			protected: true,
			method: 'GET',
			params: undefined,
			url: `${auditAPIUrl}/resource/${resourceType}`,
		}),

	getAuditLogsByResourceID: (resourceType: string, resourceId: number) =>
		apiRequest<AuditLogResponse>({
			headers: undefined,
			protected: true,
			method: 'GET',
			params: undefined,
			url: `${auditAPIUrl}/resource/${resourceType}/${resourceId}`,
		}),

	getAuditLogsByIP: (ip: string) =>
		apiRequest<AuditLogResponse>({
			headers: undefined,
			protected: true,
			method: 'GET',
			params: undefined,
			url: `${auditAPIUrl}/ip/${ip}`,
		}),

	getAuditLogsByDateRange: (startDate: string, endDate: string) =>
		apiRequest<AuditLogResponse>({
			headers: undefined,
			protected: true,
			method: 'GET',
			params: { start_date: startDate, end_date: endDate },
			url: `${auditAPIUrl}/date-range`,
		}),

	getRecentAuditLogs: () =>
		apiRequest<AuditLogResponse>({
			headers: undefined,
			protected: true,
			method: 'GET',
			params: undefined,
			url: `${auditAPIUrl}/recent`,
		}),

	cleanupOldAuditLogs: () =>
		apiRequest<GenericMessageResponse>({
			headers: undefined,
			protected: true,
			method: 'POST',
			params: undefined,
			url: `${auditAPIUrl}/cleanup`,
		}),
}

export default auditService
