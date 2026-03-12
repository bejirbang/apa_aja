<script setup>
import { onMounted, ref } from "vue"
import { useRouter } from "vue-router"
import { GrpcWebFetchTransport } from "@protobuf-ts/grpcweb-transport"
import { UserServiceClient } from "../grpc-ts/proto/user.client"
import { setCurrentUser, setToken } from "../auth"

const router = useRouter()
const authMessage = ref("")
const googleReady = ref(false)
const googleClientId = import.meta.env.VITE_GOOGLE_CLIENT_ID

const transport = new GrpcWebFetchTransport({ baseUrl: "http://localhost:8080" })
const client = new UserServiceClient(transport)

const initGoogle = () => {
  if (!googleClientId) {
    authMessage.value = "VITE_GOOGLE_CLIENT_ID belum diatur."
    return false
  }

  if (!window.google || !window.google.accounts || !window.google.accounts.id) {
    return false
  }

  window.google.accounts.id.initialize({
    client_id: googleClientId,
    callback: async (response) => {
      if (!response || !response.credential) {
        authMessage.value = "Gagal mendapatkan token Google."
        return
      }

      try {
        const call = client.loginWithGoogle({ idToken: response.credential })
        const result = await call.response
        setToken(result.accessToken)
        setCurrentUser(result.user)
        authMessage.value = "Login berhasil."
        router.push("/")
      } catch (err) {
        console.error(err)
        authMessage.value = "Login gagal."
      }
    },
  })

  const target = document.getElementById("google-signin")
  if (target) {
    window.google.accounts.id.renderButton(target, {
      theme: "outline",
      size: "large",
      shape: "pill",
      text: "signin_with",
    })
  }

  googleReady.value = true
  return true
}

onMounted(() => {
  if (initGoogle()) return
  const timer = setInterval(() => {
    if (initGoogle()) {
      clearInterval(timer)
    }
  }, 500)
})
</script>

<template>
  <div class="login">
    <h3>Login Google</h3>
    <div id="google-signin" class="google-signin"></div>
    <p v-if="!googleReady" class="status-message">Memuat tombol Google...</p>
    <p v-if="authMessage" class="status-message">{{ authMessage }}</p>
  </div>
</template>

<style scoped>
.login {
  text-align: left;
}

.google-signin {
  margin-bottom: 10px;
}

.status-message {
  margin-top: 10px;
  font-weight: bold;
}
</style>
