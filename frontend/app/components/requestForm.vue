<script setup lang="ts">
import * as z from "zod";
import type { FormSubmitEvent } from "@nuxt/ui";

const emit = defineEmits(["success"]);

const schema = z.object({
  name: z.string("Invalid name").min(2, "Must be at least 8 characters long"),
  giftReceiver: z.number("receiver ID is required"),
});

type Schema = z.output<typeof schema>;

const state = reactive<Partial<Schema>>({
  name: undefined,
  giftReceiver: undefined,
});

const receivers = [1, 2];
async function onSubmit(event: FormSubmitEvent<Schema>) {
  try {
    await $fetch("http://localhost:8000/api/occasions/", {
      method: "POST",
      credentials: "include",
      body: event.data,
    });

    emit("success", event.data);
  } catch (e: FetchError) {
    console.log(e.data);
  }
}
</script>

<template>
  <UForm :schema="schema" :state="state" class="space-y-4" @submit="onSubmit">
    <UFormField label="Name" name="name">
      <UInput v-model="state.name" />
    </UFormField>

    <UFormField label="Receiver" name="giftReceiver">
      <USelectMenu v-model="state.giftReceiver" :items="receivers" />
    </UFormField>

    <UButton type="submit">Submit</UButton>
  </UForm>
</template>
