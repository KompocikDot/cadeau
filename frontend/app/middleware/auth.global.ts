const allowed = ["/auth/login", "/auth/register", "/"];

export default defineNuxtRouteMiddleware((to) => {
  if (allowed.includes(to.path)) {
    return;
  }

  // if (useCookie('token').value === undefined) {
  // 	return navigateTo('/auth/login')
  // }
});
