<script setup lang="ts">
import type { NavigationMenuItem } from "@nuxt/ui";
import type { UserOccasion } from "../../types/responses";
import type { FetchError } from "ofetch";

definePageMeta({
  layout: "dashboard",
});

const route = useRoute();
const toast = useToast();

const requests = ref<UserOccasion[]>([]);
const modalOpen = ref(false);
const fetchingRequests = ref(true);

const items: NavigationMenuItem[] = [
  {
    label: "Your requests",
    icon: "i-lucide-scroll",
    active: route.path.startsWith("/dashboard"),
  },
  {
    label: "Request lists",
    icon: "i-lucide-notepad-text",
  },
];

const handleSuccess = (data: { name: string; giftReceiver: number }) => {
  requests.value.unshift({ Name: data.name, Id: 0 });
  modalOpen.value = false;

  toast.add({
    title: "Success",
    description: "The request have been created",
    color: "success",
  });
};

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

onMounted(() => fetchRequests());
</script>

<template>
  <UDashboardToolbar>
    <template #left>
      <UNavigationMenu :items="items" highlight class="flex-1" />
    </template>
    <template #right> </template>
  </UDashboardToolbar>
  <UContainer>
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
            <RequestForm @success="handleSuccess" />
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
        :items="requests.map((r) => ({ title: r.Name, id: r.Id }))"
        orientation="vertical"
        class="w-full data-[orientation=vertical]:h-120 border border-default rounded-sm"
      >
        <UPageCard
          v-bind="item"
          :variant="index % 2 === 0 ? 'soft' : 'outline'"
          class="rounded-none"
        />
      </UScrollArea>
    </div>
  </UContainer>
</template>
