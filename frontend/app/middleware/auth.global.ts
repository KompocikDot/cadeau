// middleware/auth.global.ts
export default defineNuxtRouteMiddleware((to) => {
  const { isLoggedIn } = useAuth();

  const publicRoutes = ["/auth/login", "/auth/register", "/"];

  if (!isLoggedIn.value && !publicRoutes.includes(to.path)) {
    return navigateTo("/auth/login");
  }

  if (isLoggedIn.value && publicRoutes.includes(to.path)) {
    return navigateTo("/user/requests");
  }
});
