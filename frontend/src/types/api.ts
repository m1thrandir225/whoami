type RequestOptions = {
	url: string
	protected: boolean | undefined
	headers: Record<string, string> | undefined
	params: Record<string, string> | undefined
	withCredentials?: boolean
}

export type ApiRequestOptions = RequestOptions & {
	method: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE'
	data?: Record<string, unknown>
	responseType?: 'json' | 'blob' | 'stream' | 'arraybuffer'
}

export type MultipartRequestOptions<T extends Record<string, unknown>> =
	RequestOptions & {
		method: 'GET' | 'POST'
		data: T
	}
