// Forwards /api/* from the Nuxt server to the Go backend so the browser only
// talks to one origin and no CORS configuration is needed on the Go side.
// The target is read at runtime from NUXT_API_URL.
export default defineEventHandler(async (event): Promise<unknown> => {
  const { apiUrl } = useRuntimeConfig()
  const path = event.context.params?.path ?? ''
  return await $fetch(`${apiUrl}/api/${path}`)
})
