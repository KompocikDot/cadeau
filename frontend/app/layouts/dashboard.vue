<script setup lang="ts">
import type { NavigationMenuItem } from "@nuxt/ui";

const route = useRoute();
const { logout } = useAuth();

const items = computed<NavigationMenuItem[]>(() => [
  {
    label: "Your requests",
    icon: "i-lucide-scroll",
    to: "/user/requests",
    active: route.path.startsWith("/user/requests"),
  },
  {
    label: "Request lists",
    icon: "i-lucide-notepad-text",
    to: "/occasions",
    active: route.path.startsWith("/occasions"),
  },
]);

const onLogout = async () => {
  try {
    await $api("http://localhost:8000/api/auth/logout/");
    logout();
  } catch (e) {
    console.log(e);
  }
};
</script>

<template>
  <UHeader>
    <template #title>
      <span class="text-lg hover:text-primary">CADEAU</span>
    </template>

    <template #right>
      <UColorModeButton />

      <UButton
        label="Logout"
        trailing-icon="i-lucide-arrow-right"
        variant="outline"
        @click="onLogout"
      />
    </template>
  </UHeader>

  <UDashboardToolbar>
    <template #left>
      <UNavigationMenu :items="items" highlight class="flex-1" />
    </template>
    <template #right> </template>
  </UDashboardToolbar>

  <UMain>
    <UContainer>
      <NuxtPage />
    </UContainer>
  </UMain>
  <Footer />
</template>
