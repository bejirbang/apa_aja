<script setup>
import { ref } from "vue"

const user = ref(null)
const searchId = ref("")
const newUser = ref({ name: "", age: "" })
const message = ref("")
const loading = ref(false)

// Fungsi untuk mencari User berdasarkan ID
const fetchUser = async () => {
  if (!searchId.value) return
  loading.value = true
  message.value = ""
  try {
    const res = await fetch(`http://localhost:8080/user?id=${searchId.value}`)
    if (!res.ok) throw new Error("User tidak ditemukan")
    user.value = await res.json()
  } catch (err) {
    user.value = null
    message.value = err.message
  } finally {
    loading.value = false
  }
}

// Fungsi untuk membuat User baru
const createUser = async () => {
  if (!newUser.value.name || !newUser.value.age) return
  loading.value = true
  message.value = ""
  try {
    const res = await fetch("http://localhost:8080/user/create", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        name: newUser.value.name,
        age: parseInt(newUser.value.age)
      }),
    })
    if (!res.ok) throw new Error("Gagal membuat user")
    const data = await res.json()
    message.value = `User ${data.name} berhasil dibuat dengan ID: ${data.id}!`
    newUser.value = { name: "", age: "" }
  } catch (err) {
    message.value = err.message
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="user-manager">
    <!-- Section Search -->
    <div class="section">
      <h3>🔍 Cari User</h3>
      <div class="input-group">
        <input v-model="searchId" type="number" placeholder="Masukkan User ID..." @keyup.enter="fetchUser" />
        <button @click="fetchUser" :disabled="loading">Cari</button>
      </div>
      
      <div v-if="user" class="result-card">
        <p><strong>ID:</strong> {{ user.id }}</p>
        <p><strong>Nama:</strong> {{ user.name }}</p>
        <p><strong>Umur:</strong> {{ user.age }}</p>
      </div>
    </div>

    <hr />

    <!-- Section Create -->
    <div class="section">
      <h3>➕ Tambah User Baru</h3>
      <div class="form">
        <input v-model="newUser.name" type="text" placeholder="Nama Lengkap" />
        <input v-model="newUser.age" type="number" placeholder="Umur" />
        <button class="btn-create" @click="createUser" :disabled="loading">Simpan User</button>
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

input:focus {
  border-color: #4facfe;
  box-shadow: 0 0 0 2px rgba(79, 172, 254, 0.2);
}

button {
  padding: 10px 20px;
  background: #4facfe;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-weight: bold;
  transition: all 0.3s ease;
}

button:hover {
  background: #00f2fe;
  transform: translateY(-1px);
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

.btn-create:hover {
  background: #27ae60;
}

.result-card {
  margin-top: 15px;
  padding: 15px;
  background: #fff;
  border-left: 4px solid #4facfe;
  border-radius: 4px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.05);
}

.form {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.status-message {
  margin-top: 15px;
  font-weight: bold;
  color: #e67e22;
  text-align: center;
}

hr {
  border: 0;
  border-top: 1px solid #eee;
  margin: 20px 0;
}
</style>