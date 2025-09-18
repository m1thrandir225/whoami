import config from '@/lib/config'
import { apiRequest } from './api.service'
import type {
	RequestPasswordResetRequest,
	ResetPasswordRequest,
	VerifyResetTokenRequest,
} from '@/types/api/auth.requests'
import type { VerifyResetTokenResponse } from '@/types/api/auth.responses'
import type { GenericMessageResponse } from '@/types/api/generic.response'

const passwordResetAPIUrl = `${config.apiUrl}/password-reset`

const passwordResetService = {
	requestPasswordReset: (input: RequestPasswordResetRequest) =>
		apiRequest<void>({
			headers: undefined,
			protected: false,
			method: 'POST',
			params: undefined,
			url: `${passwordResetAPIUrl}/request`,
			data: input,
		}),
	verifyPassword: (input: VerifyResetTokenRequest) =>
		apiRequest<VerifyResetTokenResponse>({
			protected: false,
			headers: undefined,
			method: 'POST',
			data: input,
			url: `${passwordResetAPIUrl}/verify`,
			params: undefined,
		}),

	resetPassword: (input: ResetPasswordRequest) =>
		apiRequest<GenericMessageResponse>({
			protected: false,
			headers: undefined,
			method: 'POST',
			data: input,
			url: `${passwordResetAPIUrl}/reset`,
			params: undefined,
		}),
}

export default passwordResetService
