<script setup lang="ts">
import type { UserOccasion } from "../../types/responses";

const route = useRoute();

const { data, status, error } = await useFetch<UserOccasion>(
  `http://localhost:8000/api/user/me/occasions/${route.params.id}/`,
  { credentials: "include" },
);

if (error.value) {
  throw createError({ status: 404, statusText: "Page Not Found" });
}
</script>

<template>
  <div v-if="status === 'success'">{{ data?.Name }}</div>
</template>
