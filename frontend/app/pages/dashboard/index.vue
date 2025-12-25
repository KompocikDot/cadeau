<script setup lang="ts">
import type { NavigationMenuItem } from "@nuxt/ui";

const route = useRoute();

const items: NavigationMenuItem[][] = [
  [
    {
      label: "Your requests",
      icon: "i-lucide-scroll",
      active: route.path.startsWith("/dashboard"),
    },
    {
      label: "Request lists",
      icon: "i-lucide-notepad-text",
    },
  ],
  [
    {
      label: "Logout",
      icon: "i-lucide-log-out",
      to: "https://github.com/nuxt/ui/issues",
    },
  ],
];

definePageMeta({
  layout: "dashboard",
});
</script>

<template>
  <UDashboardToolbar>
    <UNavigationMenu :items="items" highlight class="flex-1" />
  </UDashboardToolbar>
  <UContainer>
    <div class="py-5">
      <div class="flex justify-between">
        <div class="text-3xl">Your requests</div>
        <UModal>
          <UButton label="Create new" />

          <template #content>
            <Placeholder class="h-48 m-4" />
          </template>
        </UModal>
      </div>
      <ul>
        <li>aa</li>
        <li>aa</li>
        <li>aa</li>
        <li>aa</li>
        <li>aa</li>
        <li>aa</li>
      </ul>
    </div>
  </UContainer>
</template>
