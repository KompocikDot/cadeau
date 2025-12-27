<script setup lang="ts">
import type { UserOccasion } from "~/types/responses";
import type { FetchError } from "ofetch";
import type { FormSubmitEvent } from "@nuxt/ui";

definePageMeta({
  layout: "dashboard",
});

const toast = useToast();

const requests = ref<UserOccasion[]>([]);
const modalOpen = ref(false);
const fetchingRequests = ref(true);

const fetchRequests = async () => {
  try {
    const data = await $api<UserOccasion[]>(
      "http://localhost:8000/api/user/me/occasions/",
      { credentials: "include" },
    );

    requests.value = data;
    fetchingRequests.value = false;
  } catch (e) {
    console.log(e);
  }
};

async function createOccasion(event: FormSubmitEvent<Schema>) {
  try {
    const data = await $api<{ id: number }>(
      "http://localhost:8000/api/user/me/occasions/",
      {
        method: "POST",
        body: event.data,
      },
    );

    requests.value.unshift({ name: event.data.name, id: data.id });
    modalOpen.value = false;

    toast.add({
      title: "Success",
      description: "The request have been created",
      color: "success",
    });
  } catch (e) {
    console.log((e as FetchError).data);
  }
}

onMounted(() => fetchRequests());
</script>

<template>
  <div class="py-5">
    <div class="flex justify-between py-5">
      <div class="text-3xl">Your requests</div>
      <UModal
        :dismissible="false"
        title="Create new request"
        v-model:open="modalOpen"
      >
        <UButton label="Create new" />

        <template #body>
          <RequestForm @submit="createOccasion" />
        </template>
      </UModal>
    </div>
    <USkeleton v-if="fetchingRequests" class="h-96" />
    <UEmpty
      v-else-if="requests.length === 0"
      title="No requests found"
      description="It looks like you haven't added any requests. Create one to get started."
    />
    <UScrollArea
      v-else
      v-slot="{ item, index }"
      :items="requests.map((r) => ({ title: r.name, id: r.id }))"
      orientation="vertical"
      class="w-full data-[orientation=vertical]:h-120 border border-default rounded-sm"
    >
      <UPageCard
        v-bind="item"
        :variant="index % 2 === 0 ? 'soft' : 'outline'"
        class="rounded-none"
        :to="{ name: 'user-requests-id', params: { id: item.id } }"
      />
    </UScrollArea>
  </div>
</template>
