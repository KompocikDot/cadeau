<script setup lang="ts">
import * as z from "zod";
import type { FetchError } from "ofetch";
import type { UserOccasion, Guest } from "~/types/responses";
import type { FormSubmitEvent } from "@nuxt/ui";

const route = useRoute();
const toast = useToast();

const open = defineModel<boolean>("open", { required: true });
const request = defineModel<UserOccasion>("request", { required: true });
const loadingUsers = ref(false);
const guests = ref<Guest[]>([]);

const schema = z.object({
  name: z.string("Invalid name").min(2, "Must be at least 2 characters long"),
  guests: z.array(z.int()),
});

type Schema = z.output<typeof schema>;

const state = reactive<Partial<Schema>>({
  name: "",
  guests: [],
});

async function editOccasion(event: FormSubmitEvent<Schema>) {
  try {
    await $api(
      `http://localhost:8000/api/user/me/occasions/${route.params.id}/`,
      {
        method: "PATCH",
        body: event.data,
      },
    );

    request.value = {
      name: event.data.name,
      id: request.value.id,
      guests: event.data.value
        .filter((u) => event.data.guests.includes(u.id))
        .map((u) => ({ id: u.id, username: u.username })),
    };
    open.value = false;

    toast.add({
      title: "The request have been edited",
      color: "success",
      icon: "i-lucide-circle-check",
    });
  } catch (e) {
    console.log(e as FetchError);
  }
}

const fetchUsers = async () => {
  try {
    const users = await $api(`http://localhost:8000/api/users/`);
    guests.value = Array.from(new Set([...users, ...state.guests]));
  } catch (e) {}
};

watch(open, async (isOpen) => {
  if (isOpen) {
    state.name = request.value.name;
    state.guests = request.value.guests;
    loadingUsers.value = true;
    await fetchUsers();
    loadingUsers.value = false;
  }
});
</script>

<template>
  <UModal :dismissible="false" title="Edit request" v-model:open="open">
    <template #body>
      <UForm
        :schema="schema"
        :state="state"
        class="space-y-4"
        @submit="editOccasion"
      >
        <UFormField label="Name" name="name">
          <UInput v-model="state.name" class="w-full" />
        </UFormField>

        <UFormField label="Guests" name="guests">
          <USelectMenu
            multiple
            v-model="state.guests"
            :items="guests"
            valueKey="id"
            labelKey="username"
            class="w-full"
          />
        </UFormField>

        <UButton type="submit">Submit</UButton>
      </UForm>
    </template>
  </UModal>
</template>
