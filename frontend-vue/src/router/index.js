import { createRouter, createWebHistory } from "vue-router"
import User from "../components/User.vue"

const routes = [
  {
    path: "/",
    component: User
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router