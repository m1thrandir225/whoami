import type { AxiosInstance } from 'axios'
import axios, { AxiosError } from 'axios'
import { useAuthStore } from '@/stores/auth'
import type { ApiRequestOptions, MultipartRequestOptions } from '@/types/api'

import config from '@/lib/config'
import { buildFormData } from '@/lib/form'
import authService from './auth.service'

interface FailedQueuePromise {
	resolve: (token: string | null) => void
	reject: (error: unknown) => void
}

const createApiInstance = (): AxiosInstance => {
	const api = axios.create({
		baseURL: config.apiUrl,
		headers: {
			'Content-Type': 'application/json',
		},
	})

	let isRefreshing = false
	let failedQueue: FailedQueuePromise[] = []

	const processQueue = (
		error: unknown | null,
		token: string | null = null,
	): void => {
		failedQueue.forEach((prom) => {
			if (error) {
				prom.reject(error)
			} else {
				prom.resolve(token)
			}
		})

		failedQueue = []
	}

	/**
	 * Request Interceptor
	 * Used for adding JWT token if the request is protected
	 */
	api.interceptors.request.use((config) => {
		const accessToken = useAuthStore.getState().accessToken

		const isProtected = config.headers?.protected !== false

		if (isProtected && accessToken) {
			config.headers.Authorization = `Bearer ${accessToken}`
		}

		if (config.headers?.protected !== undefined) {
			delete config.headers.protected
		}
		return config
	})

	/**
	 * Used for refreshing authentication token via Cookie.
	 */
	api.interceptors.response.use(
		(response) => response,
		async (error) => {
			const originalRequest = error.config

			if (!originalRequest || originalRequest._isRetry) {
				return Promise.reject(error)
			}

			const status = error.response?.status

			if (
				status === 401 &&
				!originalRequest.url?.includes('/auth/refresh-token')
			) {
				originalRequest._isRetry = true

				if (isRefreshing) {
					return new Promise((resolve, reject) => {
						failedQueue.push({ resolve, reject })
					})
						.then((token) => {
							originalRequest.headers['Authorization'] = `Bearer ${token}`
							return api(originalRequest)
						})
						.catch((err) => Promise.reject(err))
				}
				isRefreshing = true

				try {
					const authStore = useAuthStore.getState()
					const canRefresh = authStore.checkAuth()

					if (!canRefresh) {
						throw new Error('User is logged out or refresh token expired')
					}

					const refreshToken = authStore.refreshToken

					if (!refreshToken) {
						throw new Error('Missing refresh token')
					}

					const newTokens = await authService.refreshToken({
						refresh_token: refreshToken,
					})
					if (!newTokens.access_token) {
						throw new Error('Refresh endpoint did not return an access token')
					}

					authStore.setAccessToken(newTokens.access_token, newTokens.expires_at)

					originalRequest.headers['Authorization'] =
						`Bearer ${newTokens.access_token}`

					processQueue(null, newTokens.access_token)

					return api(originalRequest)
				} catch (refreshError) {
					processQueue(refreshError, null)

					useAuthStore.getState().logout()

					return Promise.reject(refreshError)
				} finally {
					isRefreshing = false
				}
			}
			return Promise.reject(error)
		},
	)

	return api
}

const api = createApiInstance()

/**
 * A generic request abstraction for json requests
 * T is the expected type/interface outcome
 */
export const apiRequest = async <T>(config: ApiRequestOptions) => {
	try {
		const response = await api.request<T>({
			url: config.url,
			method: config.method,
			headers: {
				...config.headers,
			},
			params: config.params,
			data: config.data,
			withCredentials: config.withCredentials,
			responseType: config.responseType,
		})

		return response.data
	} catch (e: unknown) {
		if (e instanceof AxiosError) {
			throw new Error(e.response?.data.error)
		} else {
			throw new Error('Unknown error happened.')
		}
	}
}

/**
 * A generic request abstraction for multipart requests
 * T is the type/interface of the multipart request
 * R is the type/interface of the expected outcome
 */

export const multipartApiRequest = async <
	T extends Record<string, unknown>,
	R = unknown,
>(
	config: MultipartRequestOptions<T>,
): Promise<R> => {
	try {
		const formData = buildFormData(config.data)

		const response = await api.request<R>({
			url: config.url,
			method: config.method,
			data: formData,
			headers: {
				...config.headers,
				'Content-Type': 'multipart/form-data',
				protected: config.protected,
			},
			withCredentials: config.withCredentials,
			params: config.params,
		})

		return response.data
	} catch (e: unknown) {
		if (e instanceof AxiosError) {
			throw new Error(e.response?.data.error)
		} else {
			throw new Error('Unknown error happened.')
		}
	}
}

export default api
