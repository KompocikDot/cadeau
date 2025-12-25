<script setup lang="ts">
import type { NavigationMenuItem } from '@nuxt/ui'

const route = useRoute()

const items = computed<NavigationMenuItem[]>(() => [
  {
    label: 'Login',
    to: '/auth/login',
    active: route.path.startsWith('/auth/login')
  },
  {
    label: 'Register',
    to: '/auth/register',
    active: route.path.startsWith('/auth/register')
  },
  {
    label: 'Dashboard',
    to: '/dashboard',
    active: route.path.startsWith('/dashboard')
  },
])
</script>

<template>
  <UHeader>
    <template #title>
	<span class="text-lg hover:text-primary">CADEAU</span>
    </template>

    <UNavigationMenu :items="items" />

		<template #body>
           <UNavigationMenu orientation="vertical" :items="links" />
      </template>


    <template #right>
      <UColorModeButton />

      <UTooltip text="Open on GitHub" :kbds="['meta', 'G']">
        <UButton
          color="neutral"
          variant="ghost"
          to="https://github.com/KompocikDot/cadeau"
          target="_blank"
          icon="i-simple-icons-github"
          aria-label="GitHub"
        />
      </UTooltip>
    </template>

  </UHeader>
</template>

