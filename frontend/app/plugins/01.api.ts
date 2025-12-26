export default defineNuxtPlugin((nuxtApp) => {
  const config = useRuntimeConfig();

  const $customFetch = $fetch.create({
    baseURL: (config.baseUrl as string) ?? "https://api.nuxt.com",
    credentials: "include",
    async onResponseError({ response }) {
      if (response.status === 401) {
        await nuxtApp.runWithContext(() => navigateTo("/auth/login"));
      }
    },
  });

  return {
    provide: {
      customFetch: $customFetch,
    },
  };
});
