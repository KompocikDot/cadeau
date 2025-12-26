<script setup lang="ts">
import * as z from "zod";
import type { FormSubmitEvent } from "@nuxt/ui";
import type { FetchError } from "ofetch";

const emit = defineEmits(["submit"]);

const schema = z.object({
  name: z.string("Invalid name").min(2, "Must be at least 8 characters long"),
});

type Schema = z.output<typeof schema>;

const state = reactive<Partial<Schema>>({
  name: undefined,
});

const receivers = [1, 2];
</script>

<template>
  <UForm
    :schema="schema"
    :state="state"
    class="space-y-4"
    @submit="emit('submit', $event)"
  >
    <UFormField label="Name" name="name">
      <UInput v-model="state.name" />
    </UFormField>

    <UButton type="submit">Submit</UButton>
  </UForm>
</template>
