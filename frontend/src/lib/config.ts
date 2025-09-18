const baseUrl = import.meta.env.VITE_BACKEND_URL

const apiVersion = import.meta.env.VITE_BACKEND_API_VERSION

const apiUrl = `${baseUrl}/api/${apiVersion}`

export default {
	apiUrl,
	baseUrl,
}
