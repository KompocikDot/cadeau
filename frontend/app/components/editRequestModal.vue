<script setup lang="ts">
import * as z from "zod";
import type { FetchError } from "ofetch";
import type { UserOccasion } from "~/types/responses";
import type { FormSubmitEvent } from "@nuxt/ui";

const route = useRoute();
const toast = useToast();

const open = defineModel<boolean>("open", { required: true });
const request = defineModel<UserOccasion | {}>("request", { required: true });

const schema = z.object({
  name: z.string("Invalid name").min(2, "Must be at least 2 characters long"),
});

type Schema = z.output<typeof schema>;

const state = reactive<Partial<Schema>>({
  name: undefined,
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

    request.value = { name: event.data.name, id: route.params.id };
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
          <UInput v-model="state.name" />
        </UFormField>

        <UButton type="submit">Submit</UButton>
      </UForm>
    </template>
  </UModal>
</template>
