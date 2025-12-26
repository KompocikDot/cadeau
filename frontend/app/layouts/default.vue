<template>
  <UHeader>
    <template #title>
      <span class="text-lg hover:text-primary">CADEAU</span>
    </template>

    <template #right>
      <UColorModeButton />

      <UButton label="Login" to="/auth/login" variant="outline" />
      <UButton label="Register" to="/auth/register" />
    </template>
  </UHeader>
  <UMain>
    <UContainer>
      <NuxtPage />
    </UContainer>
  </UMain>
  <Footer />
</template>
