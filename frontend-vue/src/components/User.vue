<script setup>
import { onMounted, ref } from "vue"

const user = ref(null)

onMounted(async () => {
  const res = await fetch("http://localhost:8080/user")
  user.value = await res.json()
})
</script>

<template>
  <div>
    <h1>User</h1>

    <div v-if="user">
      <p>ID: {{ user.id }}</p>
      <p>Name: {{ user.name }}</p>
      <p>Age: {{ user.age }}</p>
    </div>
  </div>
</template>