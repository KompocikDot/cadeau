<script setup lang="ts">
import type { UserOccasion, Gift } from "~/types/responses";
import type { FormSubmitEvent } from "@nuxt/ui";
import type { TableColumn } from "@nuxt/ui";
import { h, resolveComponent } from "#imports";
import type { Row } from "@tanstack/vue-table";
import type { FetchError } from "ofetch";

definePageMeta({
  layout: "dashboard",
});

const route = useRoute();
const toast = useToast();

const UDropdownMenu = resolveComponent("UDropdownMenu");
const ULink = resolveComponent("ULink");
const UButton = resolveComponent("UButton");

const editRequestModalOpen = ref(false);
const editGiftModalOpen = ref(false);
const addModalOpen = ref(false);
const request = ref<UserOccasion>({ name: "", id: 0 });
const fetchingRequest = ref(true);
const gifts = ref<Gift[]>([]);

const removeRequest = async () => {
  try {
    await $api(
      `http://localhost:8000/api/user/me/occasions/${route.params.id}/`,
      { method: "DELETE" },
    );

    toast.add({
      title: "Request successfully removed",
      color: "success",
      icon: "i-lucide-circle-check",
    });
    navigateTo("/user/requests");
  } catch (e) {
    console.log((e as FetchError).data);
  }
};

const removeGift = async (giftId: number) => {
  try {
    await $api(
      `http://localhost:8000/api/user/me/occasions/${route.params.id}/gifts/${giftId}/`,
      { method: "DELETE" },
    );

    gifts.value = gifts.value.filter((g) => g.id !== giftId);
    toast.add({
      title: "Gift successfully removed",
      color: "success",
      icon: "i-lucide-circle-check",
    });
  } catch (e) {
    console.log((e as FetchError).data);
  }
};

function getRowItems(row: Row<Gift>) {
  return [
    {
      type: "label",
      label: "Actions",
    },
    { type: "separator" },
    {
      label: "Edit gift details",
      onSelect() {
        editGiftModalOpen.value = true;
      },
    },
    {
      label: "Delete gift",
      onSelect() {
        removeGift(row.original.id);
      },
    },
  ];
}

const columns: TableColumn<Gift>[] = [
  { accessorKey: "name", header: "Name" },
  {
    accessorKey: "url",
    header: "URL",
    cell: ({ row }) => {
      const url = row.original.url;
      if (!url) {
        return "-";
      }

      return h(
        ULink,
        { to: url, class: "hover:text-primary click:text-primary" },
        [new URL(url).host],
      );
    },
  },
  {
    id: "actions",
    cell: ({ row }) => {
      return h(
        "div",
        { class: "text-right" },
        h(
          UDropdownMenu,
          {
            content: {
              align: "end",
            },
            items: getRowItems(row),
            "aria-label": "Actions dropdown",
          },
          () =>
            h(UButton, {
              icon: "i-lucide-ellipsis-vertical",
              color: "neutral",
              variant: "ghost",
              class: "ml-auto",
              "aria-label": "Actions dropdown",
            }),
        ),
      );
    },
  },
];

const fetchRequest = async () => {
  try {
    const [requestData, giftsData] = await Promise.all([
      $api<UserOccasion>(
        `http://localhost:8000/api/user/me/occasions/${route.params.id}`,
      ),
      $api<Gift[]>(
        `http://localhost:8000/api/user/me/occasions/${route.params.id}/gifts/`,
      ),
    ]);

    request.value = requestData;
    gifts.value = giftsData;

    fetchingRequest.value = false;
  } catch (e) {
    console.log(e);
  }
};

onMounted(() => fetchRequest());
</script>

<template>
  <div class="py-5">
    <div class="flex justify-between py-5">
      <div class="text-3xl">Request: {{ request?.name }}</div>
      <div class="flex gap-x-2">
        <AddGiftModal v-model:open="addModalOpen" v-model:gifts="gifts" />
        <EditRequestModal
          v-model:open="editRequestModalOpen"
          v-model:request="request"
        />
        <UDropdownMenu
          :content="{ align: 'end' }"
          aria-label="Actions dropdown"
          :items="[
            {
              type: 'label',
              label: 'Actions',
            },
            { type: 'separator' },
            {
              label: 'Edit request details',
              onSelect() {
                editRequestModalOpen = true;
              },
            },
            {
              label: 'Delete request',
              onSelect() {
                removeRequest();
              },
            },
          ]"
        >
          <UButton
            aria-label="Actions dropdown"
            icon="i-lucide-ellipsis-vertical"
            color="neutral"
            variant="ghost"
            class="ml-auto"
          />
        </UDropdownMenu>
      </div>
    </div>
    <USkeleton v-if="fetchingRequest" class="h-96" />
    <UTable :data="gifts" :columns="columns" class="flex-1" />
  </div>
  <EditGiftModal v-model:open="editGiftModalOpen" />
</template>
