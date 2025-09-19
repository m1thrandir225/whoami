import config from '@/lib/config'
import type {
	LoginRequest,
	RefreshTokenRequest,
	RegisterRequest,
} from '@/types/api/auth.requests'
import { apiRequest } from './api.service'
import type {
	RegisterResponse,
	LoginResponse,
	RefreshTokenResponse,
} from '@/types/api/auth.responses'
import type { User } from '@/types/models/user'
import type { GenericMessageResponse } from '@/types/api/generic.response'

const authAPIUrl = `${config.apiUrl}`

const authService = {
	register: (input: RegisterRequest) =>
		apiRequest<RegisterResponse>({
			headers: undefined,
			protected: false,
			method: 'POST',
			params: undefined,
			url: `${authAPIUrl}/register`,
			data: input,
		}),
	login: (input: LoginRequest) =>
		apiRequest<LoginResponse>({
			headers: undefined,
			protected: false,
			method: 'POST',
			params: undefined,
			url: `${authAPIUrl}/login`,
			data: input,
		}),
	refreshToken: (input: RefreshTokenRequest) =>
		apiRequest<RefreshTokenResponse>({
			headers: undefined,
			protected: false,
			method: 'POST',
			params: undefined,
			url: `${authAPIUrl}/refresh-token`,
			data: input,
		}),
	getCurrentUser: () =>
		apiRequest<User>({
			headers: undefined,
			protected: true,
			method: 'GET',
			params: undefined,
			url: `${authAPIUrl}/me`,
		}),
	logout: () =>
		apiRequest<GenericMessageResponse>({
			headers: undefined,
			protected: true,
			method: 'POST',
			params: undefined,
			url: `${authAPIUrl}/logout`,
		}),
}

export default authService
