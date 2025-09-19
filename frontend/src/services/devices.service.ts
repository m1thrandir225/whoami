import config from '@/lib/config'
import { apiRequest } from './api.service'
import type { UserDevice } from '@/types/models/user_device'
import type { GenericMessageResponse } from '@/types/api/generic.response'

const devicesAPIUrl = `${config.apiUrl}/devices`

const devicesService = {
	getUserDevices: () =>
		apiRequest<UserDevice[]>({
			headers: undefined,
			protected: true,
			method: 'GET',
			params: undefined,
			url: devicesAPIUrl,
		}),

	getUserDevice: (id: number) =>
		apiRequest<UserDevice>({
			headers: undefined,
			protected: true,
			method: 'GET',
			params: undefined,
			url: `${devicesAPIUrl}/${id}`,
		}),

	updateUserDevice: (id: number, data: Partial<UserDevice>) =>
		apiRequest<UserDevice>({
			headers: undefined,
			protected: true,
			method: 'PUT',
			params: undefined,
			url: `${devicesAPIUrl}/${id}`,
			data,
		}),

	deleteUserDevice: (id: number) =>
		apiRequest<GenericMessageResponse>({
			headers: undefined,
			protected: true,
			method: 'DELETE',
			params: undefined,
			url: `${devicesAPIUrl}/${id}`,
		}),

	deleteAllUserDevices: () =>
		apiRequest<GenericMessageResponse>({
			headers: undefined,
			protected: true,
			method: 'DELETE',
			params: undefined,
			url: devicesAPIUrl,
		}),

	markDeviceAsTrusted: (id: number) =>
		apiRequest<UserDevice>({
			headers: undefined,
			protected: true,
			method: 'PATCH',
			params: undefined,
			url: `${devicesAPIUrl}/${id}/trust`,
		}),
}

export default devicesService
