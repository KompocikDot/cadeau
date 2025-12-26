<script setup lang="ts">
import type { TimelineItem, ButtonProps } from "@nuxt/ui";

const route = useRoute();

const links: ButtonProps[] = [
  {
    label: "Get started",
    to: "/auth/login",
  },
  {
    label: "Learn more",
    to: "#dupa",
    color: "neutral",
    variant: "subtle",
    trailingIcon: "i-lucide-arrow-right",
  },
];

const items: TimelineItem[] = [
  {
    title: "Add Your Wishes",
    description:
      "Paste links to the exact items you want from your favorite shops.",
    icon: "i-lucide-list-plus",
  },
  {
    title: "Share the Link",
    description:
      "Send your list to family and friends. Everyone sees the same real-time list.",
    icon: "i-lucide-share-2",
  },
  {
    title: "Ghost Mode",
    description: "Guests mark gifts as bought so others don't buy them again.",
    icon: "i-lucide-ghost",
  },
  {
    title: "Surprise!",
    description:
      "You only see the original list. You won’t know what’s been bought until you unwrap it.",
    icon: "i-lucide-party-popper",
  },
];
</script>

<template>
  <UPageHero
    title="The Wishlist That Works for Everyone."
    description="Create your list with direct links. Your guests coordinate secretly to avoid duplicates, keeping the surprise alive until the big day."
    :links="links"
  />
  <div
    class="text-xl w-full flex justify-center py-10 text-6xl sm:text-6xl text-pretty tracking-tight font-bold text-highlighted"
  >
    How that works?
  </div>
  <section id="step" class="w-full flex justify-center">
    <UTimeline size="lg" id="dupa" :items="items" />
  </section>
</template>
