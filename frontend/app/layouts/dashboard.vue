<script setup lang="ts">
const { logout } = useAuth();

const onLogout = async () => {
  try {
    await $fetch("http://localhost:8000/api/auth/logout/", {
      credentials: "include",
    });

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
  <UMain>
    <NuxtPage />
  </UMain>
  <Footer />
</template>
