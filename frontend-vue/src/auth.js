const TOKEN_KEY = "access_token"
const USER_KEY = "current_user"

export const getToken = () => localStorage.getItem(TOKEN_KEY) || ""

export const setToken = (token) => {
  if (token) {
    localStorage.setItem(TOKEN_KEY, token)
  }
}

export const clearToken = () => {
  localStorage.removeItem(TOKEN_KEY)
}

export const isAuthed = () => Boolean(getToken())

export const getAuthMeta = () => {
  const token = getToken()
  if (!token) return {}
  return { authorization: `Bearer ${token}` }
}

export const setCurrentUser = (user) => {
  if (user) {
    localStorage.setItem(USER_KEY, JSON.stringify(user))
  }
}

export const getCurrentUser = () => {
  const raw = localStorage.getItem(USER_KEY)
  if (!raw) return null
  try {
    return JSON.parse(raw)
  } catch {
    return null
  }
}

export const clearCurrentUser = () => {
  localStorage.removeItem(USER_KEY)
}
