import axios from 'axios'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || '/api',
})

// Add auth token to requests
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Handle 401 responses
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      if (window.location.pathname.startsWith('/admin')) {
        window.location.href = '/admin/login'
      }
    }
    return Promise.reject(error)
  }
)

export const getContent = () => api.get('/content')

export const updateContent = (data) => api.put('/content', data)

export const login = (username, password) =>
  api.post('/auth/login', { username, password })

export const seedAdmin = (username, password) =>
  api.post('/auth/seed', { username, password })

export const uploadImage = (file) => {
  const formData = new FormData()
  formData.append('image', file)
  return api.post('/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
}

// News
export const getNews = () => api.get('/news')
export const getNewsPost = (id) => api.get(`/news/${id}`)
export const createNews = (data) => api.post('/news', data)
export const updateNews = (id, data) => api.put(`/news/${id}`, data)
export const deleteNews = (id) => api.delete(`/news/${id}`)

// User management
export const listUsers = () => api.get('/auth/users')
export const createUser = (data) => api.post('/auth/users', data)
export const deleteUser = (id) => api.delete(`/auth/users/${id}`)

export default api
