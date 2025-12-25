<script setup lang="ts">
import * as z from 'zod'
import type { FormSubmitEvent, AuthFormField } from '@nuxt/ui'

const toast = useToast()

const fields: AuthFormField[] = [{
  name: 'username',
  type: 'text',
  label: 'Username',
  placeholder: 'Enter your username',
  required: true
}, {
  name: 'password',
  label: 'Password',
  type: 'password',
  placeholder: 'Enter your password',
  required: true
}, {
  name: 'remember',
  label: 'Remember me',
  type: 'checkbox'
}]

const schema = z.object({
  username: z.string('Invalid username'),
  password: z.string('Password is required').min(8, 'Must be at least 8 characters')
})

type Schema = z.output<typeof schema>

async function onSubmit(payload: FormSubmitEvent<Schema>) {
	try {

await $fetch("http://localhost:8000/auth/login/", {method: 'POST', body: payload.data, credentials: 'include'})

		navigateTo('/dashboard')
		console.log("should redirect")
	} catch (e) {console.log(e)}
	
}
</script>

<template>
  <div class="flex flex-col items-center justify-center gap-4 p-4">
    <UPageCard class="w-full max-w-md">
      <UAuthForm
        :schema="schema"
        title="Login"
        description="Enter your credentials to access your account."
        icon="i-lucide-user"
        :fields="fields"
	:submit="{label: 'Login'}"
        @submit="onSubmit"
			>
				 <template #description>
          Don't have an account? <ULink to="/auth/register" class="text-primary font-medium">Sign up</ULink>.
        </template>
			</UAuthForm>
    </UPageCard>
  </div>
</template>
