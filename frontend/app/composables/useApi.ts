// composables/useApi.ts
import type { UseFetchOptions } from 'nuxt/app'

export const useApi = <T>(url: string, options: UseFetchOptions<T> = {}) => {
	const router = useRouter()

	return $fetch<T>(url, {
		...options,
		baseURL: 'http://localhost:8080',
		credentials: 'include',

		async onResponseError({ response }) {
			if (response.status === 401) {
				await navigateTo('/auth/login')
			}
		}
	})
}
