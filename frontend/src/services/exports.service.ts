import config from '@/lib/config'
import { apiRequest } from './api.service'
import type { DataExport } from '@/types/models/data_export'
import type { GenericMessageResponse } from '@/types/api/generic.response'

const exportsAPIUrl = `${config.apiUrl}/exports`

const exportsService = {
	requestDataExport: (exportType: string) =>
		apiRequest<DataExport>({
			headers: undefined,
			protected: true,
			method: 'POST',
			params: undefined,
			url: exportsAPIUrl,
			data: { export_type: exportType },
		}),

	getDataExports: () =>
		apiRequest<DataExport[]>({
			headers: undefined,
			protected: true,
			method: 'GET',
			params: undefined,
			url: exportsAPIUrl,
		}),

	getDataExport: (id: number) =>
		apiRequest<DataExport>({
			headers: undefined,
			protected: true,
			method: 'GET',
			params: undefined,
			url: `${exportsAPIUrl}/${id}`,
		}),

	downloadDataExport: (id: number) =>
		apiRequest<Blob>({
			headers: undefined,
			protected: true,
			method: 'GET',
			params: undefined,
			url: `${exportsAPIUrl}/${id}/download`,
			responseType: 'blob',
		}),

	deleteDataExport: (id: number) =>
		apiRequest<GenericMessageResponse>({
			headers: undefined,
			protected: true,
			method: 'DELETE',
			params: undefined,
			url: `${exportsAPIUrl}/${id}`,
		}),
}

export default exportsService
