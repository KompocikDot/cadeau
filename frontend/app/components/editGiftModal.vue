<script setup lang="ts">
import * as z from "zod";

const open = defineModel<boolean>("open", { required: true });

const schema = z.object({
  name: z.string("Invalid name").min(2, "Must be at least 2 characters long"),
  url: z.preprocess((v) => (v === "" ? undefined : v), z.url().optional()),
});

type Schema = z.output<typeof schema>;

const state = reactive<Partial<Schema>>({
  name: undefined,
  url: undefined,
});
</script>

<template>
  <UModal :dismissible="false" title="Edit gift" v-model:open="open">
    <template #body>
      <UForm :schema="schema" :state="state" class="space-y-4" @submit="">
        <UFormField label="Name" name="name">
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
