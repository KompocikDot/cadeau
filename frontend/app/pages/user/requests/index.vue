<script setup lang="ts">
import * as z from "zod";
import type { UserOccasion } from "~/types/responses";
import type { FetchError } from "ofetch";
import type { FormSubmitEvent } from "@nuxt/ui";

definePageMeta({
  layout: "dashboard",
});

const toast = useToast();
const schema = z.object({
  name: z.string("Invalid name").min(2, "Must be at least 2 characters long"),
  guests: z.array(z.number()),
});

type Schema = z.output<typeof schema>;

const state = reactive<Partial<Schema>>({
  name: undefined,
  guests: [],
});

const requests = ref<UserOccasion[]>([]);
const modalOpen = ref(false);
const fetchingRequests = ref(true);
const fetchingGuests = ref([]);
const guests = ref([]);

const fetchRequests = async () => {
  try {
    const data = await $api<UserOccasion[]>(
      "http://localhost:8000/api/user/me/occasions/",
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
      title: "The request have been created",
      color: "success",
      icon: "i-lucide-circle-check",
    });
  } catch (e) {
    console.log((e as FetchError).data);
  }
}

const loadUsersList = async () => {
  try {
    fetchingGuests.value = true;
    guests.value = await $api(`http://localhost:8000/api/users/`);
    fetchingGuests.value = false;
  } catch (e) {}
};

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
        <UButton label="Create new" @click="loadUsersList" />

        <template #body>
          <UForm
            :schema="schema"
            :state="state"
            class="space-y-4"
            @submit="createOccasion"
          >
            <UFormField label="Name" name="name">
              <UInput v-model="state.name" class="w-full" />
            </UFormField>

            <UFormField label="Guests" name="guests">
              <USelectMenu
                :items="guests"
                v-model="state.guests"
                multiple
                valueKey="id"
                labelKey="username"
                class="w-full"
              />
            </UFormField>

            <UButton type="submit">Submit</UButton>
          </UForm>
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
