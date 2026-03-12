import { createRouter, createWebHistory } from "vue-router"
import Login from "../components/Login.vue"
import User from "../components/User.vue"
import { isAuthed } from "../auth"

const routes = [
  {
    path: "/",
    component: User
  },
  {
    path: "/login",
    component: Login
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to) => {
  if (to.path === "/login" && isAuthed()) {
    return "/"
  }
  if (to.path !== "/login" && !isAuthed()) {
    return "/login"
  }
  return true
})

export default router
