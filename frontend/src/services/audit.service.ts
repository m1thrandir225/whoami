import config from '@/lib/config'
import { apiRequest } from './api.service'
import type { AuditLog } from '@/types/models/audit_log'
import type { GenericMessageResponse } from '@/types/api/generic.response'

const auditAPIUrl = `${config.apiUrl}/audit`

const auditService = {
	getAuditLogsByUserID: (userId: number) =>
		apiRequest<AuditLog[]>({
			headers: undefined,
			protected: true,
			method: 'GET',
			params: undefined,
			url: `${auditAPIUrl}/user/${userId}`,
		}),

	getAuditLogsByAction: (action: string) =>
		apiRequest<AuditLog[]>({
			headers: undefined,
			protected: true,
			method: 'GET',
			params: undefined,
			url: `${auditAPIUrl}/action/${action}`,
		}),

	getAuditLogsByResourceType: (resourceType: string) =>
		apiRequest<AuditLog[]>({
			headers: undefined,
			protected: true,
			method: 'GET',
			params: undefined,
			url: `${auditAPIUrl}/resource/${resourceType}`,
		}),

	getAuditLogsByResourceID: (resourceType: string, resourceId: number) =>
		apiRequest<AuditLog[]>({
			headers: undefined,
			protected: true,
			method: 'GET',
			params: undefined,
			url: `${auditAPIUrl}/resource/${resourceType}/${resourceId}`,
		}),

	getAuditLogsByIP: (ip: string) =>
		apiRequest<AuditLog[]>({
			headers: undefined,
			protected: true,
			method: 'GET',
			params: undefined,
			url: `${auditAPIUrl}/ip/${ip}`,
		}),

	getAuditLogsByDateRange: (startDate: string, endDate: string) =>
		apiRequest<AuditLog[]>({
			headers: undefined,
			protected: true,
			method: 'GET',
			params: { start_date: startDate, end_date: endDate },
			url: `${auditAPIUrl}/date-range`,
		}),

	getRecentAuditLogs: () =>
		apiRequest<AuditLog[]>({
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
