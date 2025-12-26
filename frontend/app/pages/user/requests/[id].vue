<script setup lang="ts">
import type { UserOccasion } from "~/types/responses";

definePageMeta({
  layout: "dashboard",
});

const route = useRoute();
const toast = useToast();

const modalOpen = ref(false);
const request = ref<UserOccasion>({});
const fetchingRequest = ref(true);

const handleSuccess = (data: { name: string; ID: number }) => {
  request.value = { Name: data.name, Id: data.ID };
  modalOpen.value = false;

  toast.add({
    title: "Success",
    description: "The request have been created",
    color: "success",
  });
};

const fetchRequest = async () => {
  try {
    const data = await $api<UserOccasion>(
      `http://localhost:8000/api/user/me/occasions/${route.params.id}`,
      { credentials: "include" },
    );

    request.value = data;
    fetchingRequest.value = false;
  } catch (e) {
    console.log(e);
  }
};

async function onSubmit(event: FormSubmitEvent<Schema>) {
  try {
    const data = await $api("http://localhost:8000/api/user/me/occasions/", {
      method: "POST",
      body: event.data,
    });

    handleSuccess({ ...event.data, Id: data.Id });
  } catch (e) {
    console.log((e as FetchError).data);
  }
}

onMounted(() => fetchRequest());
</script>

<template>
  <div class="py-5">
    <div class="flex justify-between py-5">
      <div class="text-3xl">Request: {{ data?.Name }}</div>
      <UModal
        :dismissible="false"
        title="Edit request"
        v-model:open="modalOpen"
      >
        <UButton label="Edit request" />

        <template #body>
          <RequestForm @submit="onSubmit" />
        </template>
      </UModal>
    </div>
    <USkeleton v-if="fetchingRequest" class="h-96" />
  </div>
</template>
