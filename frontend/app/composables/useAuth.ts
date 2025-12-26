export const useAuth = () => {
  const user = useState<any>("user", () => null);
  const isLoggedIn = computed(() => !!user.value);

  async function fetchUser() {
    try {
      const data = await $api("http://localhost:8000/api/user/me/");
      user.value = data;
    } catch (error) {
      user.value = null;
    }
  }

  async function login(credentials: { username: string; password: string }) {
    await $api("http://localhost:8000/api/auth/login/", {
      method: "POST",
      body: credentials,
    });

    await fetchUser();

    return navigateTo("/user/requests");
  }

  function logout() {
    user.value = null;
    return navigateTo("/auth/login");
  }

  return { user, isLoggedIn, fetchUser, login, logout };
};
