<script setup lang="ts">
import * as z from "zod";
import type { FormSubmitEvent } from "@nuxt/ui";
import type { UserOccasion } from "~/types/responses";
import type { FetchError } from "ofetch";

const route = useRoute();
const toast = useToast();

const open = defineModel<boolean>("open", { required: true });
const gifts = defineModel<UserOccasion[]>("gifts", { required: true });

const schema = z.object({
  name: z.string("Invalid name").min(2, "Must be at least 2 characters long"),
  url: z.url().optional().or(z.literal("")),
});

type Schema = z.output<typeof schema>;

const state = reactive<Partial<Schema>>({
  name: "",
  url: "",
});

const addGift = async (event: FormSubmitEvent<Schema>) => {
  try {
    const { id } = await $api<{ id: number }>(
      `http://localhost:8000/api/user/me/occasions/${route.params.id}/gifts/`,
      { method: "POST", body: event.data },
    );

    gifts.value?.push({ ...event.data, id: id });
    open.value = false;

    toast.add({
      title: "The gift have been added",
      color: "success",
      icon: "i-lucide-circle-check",
    });
  } catch (e) {
    console.log((e as FetchError).data);
  }
};

watch(open, (isOpen) => {
  if (!isOpen) {
    state.name = "";
    state.url = "";
  }
});
</script>

<template>
  <UModal :dismissible="false" title="Add gift" v-model:open="open">
    <UButton label="Add gift" />

    <template #body>
      <UForm
        :schema="schema"
        :state="state"
        class="space-y-4"
        @submit="addGift"
      >
        <UFormField label="Name" name="name" required>
          <UInput v-model="state.name" />
        </UFormField>
        <UFormField label="URL" name="url">
          <UInput v-model="state.url" />
        </UFormField>

        <UButton type="submit">Submit</UButton>
      </UForm>
    </template>
  </UModal>
</template>
