<script setup>
import { ref } from "vue"
import { useRouter } from "vue-router"
import { GrpcWebFetchTransport } from "@protobuf-ts/grpcweb-transport"
import { UserServiceClient } from "../grpc-ts/proto/user.client"
import { clearCurrentUser, clearToken, getAuthMeta, getCurrentUser } from "../auth"

const user = ref(null)
const searchId = ref("")
const newUser = ref({ name: "", age: "" })
const message = ref("")
const loading = ref(false)
const currentUser = ref(getCurrentUser())
const router = useRouter()

const transport = new GrpcWebFetchTransport({ baseUrl: "http://localhost:8080" })
const client = new UserServiceClient(transport)

const logout = () => {
  clearToken()
  clearCurrentUser()
  currentUser.value = null
  router.push("/login")
}

const fetchUser = async () => {
  if (!searchId.value) return
  loading.value = true
  message.value = ""
  user.value = null

  try {
    const call = client.getUser({ id: Number(searchId.value) }, { meta: getAuthMeta() })
    user.value = await call.response
  } catch (err) {
    console.error(err)
    message.value = "Gagal mengambil user."
    if (String(err).includes("Unauthenticated")) {
      logout()
    }
  } finally {
    loading.value = false
  }
}

const createUser = async () => {
  if (!newUser.value.name || !newUser.value.age) return
  loading.value = true
  message.value = ""

  try {
    const call = client.createUser({
      name: newUser.value.name,
      age: Number(newUser.value.age),
    }, { meta: getAuthMeta() })
    const created = await call.response
    message.value = `User ${created.name} berhasil dibuat dengan ID: ${created.id}`
    newUser.value = { name: "", age: "" }
  } catch (err) {
    console.error(err)
    message.value = "Gagal membuat user."
    if (String(err).includes("Unauthenticated")) {
      logout()
    }
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="user-manager">
    <div v-if="currentUser" class="section auth-section">
      <div class="result-card">
        <p><strong>Nama:</strong> {{ currentUser.name }}</p>
        <p><strong>Email:</strong> {{ currentUser.email }}</p>
        <button class="btn-logout" @click="logout">Logout</button>
      </div>
      <hr />
    </div>

    <div class="section">
      <h3>Cari User</h3>
      <div class="input-group">
        <input v-model="searchId" type="number" placeholder="Masukkan ID user" @keyup.enter="fetchUser" />
        <button @click="fetchUser" :disabled="loading">Cari</button>
      </div>

      <div v-if="user" class="result-card">
        <p><strong>ID:</strong> {{ user.id }}</p>
        <p><strong>Nama:</strong> {{ user.name }}</p>
        <p><strong>Umur:</strong> {{ user.age }}</p>
      </div>
    </div>

    <hr />

    <div class="section">
      <h3>Tambah User Baru</h3>
      <div class="form">
        <input v-model="newUser.name" type="text" placeholder="Nama" />
        <input v-model="newUser.age" type="number" placeholder="Umur" />
        <button class="btn-create" @click="createUser" :disabled="loading">Simpan</button>
      </div>
    </div>

    <p v-if="message" class="status-message">{{ message }}</p>
  </div>
</template>

<style scoped>
.user-manager {
  text-align: left;
}

.section {
  margin-bottom: 20px;
}

h3 {
  font-size: 18px;
  color: #444;
  margin-bottom: 12px;
}

.input-group {
  display: flex;
  gap: 10px;
}

input {
  flex: 1;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 8px;
  outline: none;
  font-size: 14px;
}

button {
  padding: 10px 20px;
  background: #4facfe;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-weight: bold;
}

button:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.btn-create {
  width: 100%;
  margin-top: 10px;
  background: #2ecc71;
}

.btn-logout {
  margin-top: 10px;
  background: #ff6b6b;
}

.result-card {
  margin-top: 15px;
  padding: 15px;
  background: #fff;
  border-left: 4px solid #4facfe;
  border-radius: 4px;
}

.form {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.status-message {
  margin-top: 15px;
  font-weight: bold;
}

hr {
  border: 0;
  border-top: 1px solid #eee;
  margin: 20px 0;
}
</style>
