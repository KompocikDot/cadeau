import type { UseFetchOptions } from "nuxt/app";
import type { NitroFetchRequest, NitroFetchOptions } from "nitropack";

export function useApiFetch<T>(
  url: string | (() => string),
  options: UseFetchOptions<T> = {},
) {
  return useFetch(url, {
    ...options,
    $fetch: useNuxtApp().$customFetch,
  });
}

export function $api<T = unknown>(
  request: NitroFetchRequest,
  opts?: NitroFetchOptions<NitroFetchRequest>,
) {
  const { $customFetch } = useNuxtApp();
  return $customFetch<T>(request, opts);
}
