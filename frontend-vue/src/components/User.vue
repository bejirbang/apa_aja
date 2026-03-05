<script setup>
import { ref, onMounted } from "vue"

const user = ref(null)

onMounted(async () => {
  const res = await fetch("http://localhost:8080/user")
  user.value = await res.json()
})
</script>

<template>
  <div>
    <h2>User Data</h2>

    <div v-if="user">
      <p>ID: {{ user.id }}</p>
      <p>Name: {{ user.name }}</p>
      <p>Age: {{ user.age }}</p>
    </div>
  </div>
</template>