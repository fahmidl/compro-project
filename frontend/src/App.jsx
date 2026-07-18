import { Routes, Route } from 'react-router-dom'
import Layout from './components/Layout'
import HomePage from './pages/HomePage'
import AboutPage from './pages/AboutPage'
import ContactPage from './pages/ContactPage'
import NewsPage from './pages/NewsPage'
import NewsDetailPage from './pages/NewsDetailPage'
import LoginPage from './pages/admin/LoginPage'
import DashboardPage from './pages/admin/DashboardPage'
import HeroEditor from './pages/admin/HeroEditor'
import AboutEditor from './pages/admin/AboutEditor'
import ContactEditor from './pages/admin/ContactEditor'
import NewsEditor from './pages/admin/NewsEditor'
import UserManagement from './pages/admin/UserManagement'

function App() {
  return (
    <Routes>
      <Route path="/" element={<Layout />}>
        <Route index element={<HomePage />} />
        <Route path="about" element={<AboutPage />} />
        <Route path="contact" element={<ContactPage />} />
        <Route path="news" element={<NewsPage />} />
        <Route path="news/:id" element={<NewsDetailPage />} />
      </Route>
      <Route path="/admin/login" element={<LoginPage />} />
      <Route path="/admin" element={<DashboardPage />} />
      <Route path="/admin/hero" element={<HeroEditor />} />
      <Route path="/admin/about" element={<AboutEditor />} />
      <Route path="/admin/contact" element={<ContactEditor />} />
      <Route path="/admin/news" element={<NewsEditor />} />
      <Route path="/admin/users" element={<UserManagement />} />
    </Routes>
  )
}

export default App
