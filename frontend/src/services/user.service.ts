import config from '@/lib/config'
import type {
	UpdatePasswordRequest,
	UpdateUserRequest,
} from '@/types/api/auth.requests'
import { apiRequest } from './api.service'
import type { PrivacySettings } from '@/types/models/privacy_settings'
import type { GenericMessageResponse } from '@/types/api/generic.response'

const userAPIUrl = `${config.apiUrl}/user`

const userService = {
	updateUser: (userId: number, input: UpdateUserRequest) =>
		apiRequest<void>({
			headers: undefined,
			protected: true,
			method: 'PUT',
			params: undefined,
			url: `${userAPIUrl}/${userId}`,
			data: input,
		}),
	updateUserPrivacySettings: (userId: number, input: PrivacySettings) =>
		apiRequest<void>({
			headers: undefined,
			protected: true,
			method: 'PUT',
			params: undefined,
			url: `${userAPIUrl}/${userId}/privacy-settings`,
			data: input,
		}),
	updatePassword: (input: UpdatePasswordRequest) =>
		apiRequest<GenericMessageResponse>({
			headers: undefined,
			protected: true,
			method: 'POST',
			params: undefined,
			url: `${userAPIUrl}/update-password`,
			data: input,
		}),
	deactivateUser: (userId: number) =>
		apiRequest<void>({
			headers: undefined,
			protected: true,
			method: 'POST',
			params: undefined,
			url: `${userAPIUrl}/${userId}/deactivate`,
		}),
	activateUser: (userId: number) =>
		apiRequest<void>({
			headers: undefined,
			protected: true,
			method: 'POST',
			params: undefined,
			url: `${userAPIUrl}/${userId}/activate`,
		}),

	setPassword: (password: string) =>
		apiRequest<GenericMessageResponse>({
			headers: undefined,
			protected: true,
			method: 'POST',
			params: undefined,
			url: `${userAPIUrl}/set-password`,
			data: { password },
		}),
}

export default userService
