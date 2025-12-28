<script setup lang="ts">
import * as z from "zod";
import type { Gift } from "~/types/responses";
import type { FormSubmitEvent } from "@nuxt/ui";

const route = useRoute();
const toast = useToast();

type EditedGift = {
  id: number;
  name: string;
  url: string;
};

const props = defineProps({
  editedGift: { type: Object as PropType<EditedGift>, required: true },
});

const open = defineModel<boolean>("open", { required: true });
const gifts = defineModel<Gift[]>("gifts", { required: true });

const schema = z.object({
  name: z.string("Invalid name").min(2, "Must be at least 2 characters long"),
  url: z.preprocess((v) => (v === "" ? undefined : v), z.url().optional()),
});

type Schema = z.output<typeof schema>;

const state = reactive<Partial<Schema>>({
  name: "",
  url: "",
});

const updateGift = async (event: FormSubmitEvent<Schema>) => {
  try {
    await $api<{ id: number }>(
      `http://localhost:8000/api/user/me/occasions/${route.params.id}/gifts/${props.editedGift.id}/`,
      { method: "PATCH", body: event.data },
    );

    gifts.value = gifts.value.map((gift) =>
      gift.id === updatedId ? { ...gift, ...event.data } : gift,
    );

    toast.add({
      title: "The gift have been edited",
      color: "success",
      icon: "i-lucide-circle-check",
    });

    open.value = false;
  } catch (e) {
    console.log(e);
  }
};

watch(open, (isOpen) => {
  if (isOpen) {
    state.name = props.editedGift.name;
    state.url = props.editedGift.url;
  }
});
</script>

<template>
  <UModal :dismissible="false" title="Edit gift" v-model:open="open">
    <template #body>
      <UForm
        :schema="schema"
        :state="state"
        class="space-y-4"
        @submit="updateGift"
      >
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
