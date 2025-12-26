<script setup lang="ts">
import * as z from "zod";
import type { FormSubmitEvent, AuthFormField } from "@nuxt/ui";

const toast = useToast();

const fields: AuthFormField[] = [
  {
    name: "username",
    type: "text",
    label: "Username",
    placeholder: "Enter your username",
    required: true,
  },
  {
    name: "password",
    label: "Password",
    type: "password",
    placeholder: "Enter your password",
    required: true,
  },
  {
    name: "remember",
    label: "Remember me",
    type: "checkbox",
  },
];

const schema = z.object({
  username: z.string("Invalid username"),
  password: z
    .string("Password is required")
    .min(8, "Must be at least 8 characters"),
});

type Schema = z.output<typeof schema>;

const onSubmit = async (payload: FormSubmitEvent<Schema>) => {
  try {
    await $fetch("http://localhost:8000/api/auth/register/", {
      method: "POST",
      body: payload.data,
      credentials: "include",
    });

    navigateTo("/auth/login");
    console.log("should redirect");
  } catch (e) {
    console.log(e);
  }
};
</script>

<template>
  <div class="flex flex-col items-center justify-center gap-4 p-4">
    <UPageCard class="w-full max-w-md">
      <UAuthForm
        :schema="schema"
        :fields="fields"
        title="Welcome!"
        icon="i-lucide-hand"
        :submit="{ label: 'Create an account' }"
        @submit="onSubmit"
      >
        <template #description>
          Already have an account?
          <ULink to="/auth/login" class="text-primary font-medium"
            >Sign in</ULink
          >.
        </template>
        <template #footer>
          By signing in, you agree to our
          <ULink to="#" class="text-primary font-medium">Terms of Service</ULink
          >.
        </template>
      </UAuthForm>
    </UPageCard>
  </div>
</template>
