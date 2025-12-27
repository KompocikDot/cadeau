<script setup lang="ts">
import type { UserOccasion } from "~/types/responses";
import type { FormSubmitEvent } from "@nuxt/ui";

definePageMeta({
  layout: "dashboard",
});

const route = useRoute();
const toast = useToast();

const editModalOpen = ref(false);
const addModalOpen = ref(false);
const request = ref<UserOccasion | {}>({});
const fetchingRequest = ref(true);
const gifts = ref([]);

const fetchRequest = async () => {
  try {
    const data = await $api<UserOccasion>(
      `http://localhost:8000/api/user/me/occasions/${route.params.id}`,
    );

    request.value = data;

    const giftsData = await $api<UserOccasion>(
      `http://localhost:8000/api/user/me/occasions/${route.params.id}/gifts/`,
    );

    gifts.value = giftsData;
    fetchingRequest.value = false;
  } catch (e) {
    console.log(e);
  }
};

async function editOccasion(event: FormSubmitEvent<any>) {
  try {
    const data = await $api(
      `http://localhost:8000/api/user/me/occasions/${route.params.id}/`,
      {
        method: "PATCH",
        body: event.data,
      },
    );

    request.value = { name: data.name, id: data.id };
    editModalOpen.value = false;

    toast.add({
      title: "Success",
      description: "The request have been edited",
      color: "success",
    });
  } catch (e) {
    console.log((e as FetchError).data);
  }
}

const addGift = async (event: FormSubmitEvent<any>) => {
  try {
    await $api(
      `http://localhost:8000/api/user/me/occasions/${route.params.id}/gifts/`,
      { method: "POST", body: event.data },
    );

    gifts.value.push(event.data);
    addModalOpen.value = false;

    toast.add({
      title: "Success",
      description: "The gift have been added",
      color: "success",
    });
  } catch {
    console.log((e as FetchError).data);
  }
};

onMounted(() => fetchRequest());
</script>

<template>
  <div class="py-5">
    <div class="flex justify-between py-5">
      <div class="text-3xl">Request: {{ data?.name }}</div>
      <div class="flex gap-x-2">
        <UModal
          :dismissible="false"
          title="Add gift"
          v-model:open="addModalOpen"
        >
          <UButton label="Add gift" variant="outline" />

          <template #body>
            <RequestForm @submit="addGift" />
          </template>
        </UModal>
        <UModal
          :dismissible="false"
          title="Edit request"
          v-model:open="editModalOpen"
        >
          <UButton label="Edit request" />

          <template #body>
            <RequestForm @submit="editOccasion" />
          </template>
        </UModal>
      </div>
    </div>
    <USkeleton v-if="fetchingRequest" class="h-96" />
    <UTable :data="gifts" class="flex-1" />
  </div>
</template>
