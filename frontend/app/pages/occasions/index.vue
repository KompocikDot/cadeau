<script setup lang="ts">
definePageMeta({
  layout: "dashboard",
});

const loadInvites = async () => {
  try {
    await $api();
  } catch (e) {}
};

onMounted(() => {
  loadInvites();
});
</script>

<template></template>
