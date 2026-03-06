<script setup>
import { onMounted, ref } from "vue"
import { GrpcWebFetchTransport } from "@protobuf-ts/grpcweb-transport"
import { UserServiceClient } from "../grpc-ts/proto/user.client"

const user = ref(null)
const errorMessage = ref("")

const transport = new GrpcWebFetchTransport({ baseUrl: "http://localhost:8080" })
const client = new UserServiceClient(transport)

onMounted(async () => {
  try {
    const call = client.getUser({ id: 1 })
    const res = await call.response
    user.value = res
  } catch (err) {
    console.error(err)
    errorMessage.value = "Gagal mengambil data user via gRPC-Web."
  }
})
</script>

<template>
  <div>
    <h1>User</h1>

    <p v-if="errorMessage">{{ errorMessage }}</p>

    <div v-if="user">
      <p>ID: {{ user.id }}</p>
      <p>Name: {{ user.name }}</p>
      <p>Age: {{ user.age }}</p>
    </div>
  </div>
</template>
