import { Routes, Route } from 'react-router-dom'
import Layout from './components/Layout'
import HomePage from './pages/HomePage'
import AboutPage from './pages/AboutPage'
import ContactPage from './pages/ContactPage'
import LoginPage from './pages/admin/LoginPage'
import DashboardPage from './pages/admin/DashboardPage'
import HeroEditor from './pages/admin/HeroEditor'
import AboutEditor from './pages/admin/AboutEditor'
import ContactEditor from './pages/admin/ContactEditor'

function App() {
  return (
    <Routes>
      <Route path="/" element={<Layout />}>
        <Route index element={<HomePage />} />
        <Route path="about" element={<AboutPage />} />
        <Route path="contact" element={<ContactPage />} />
      </Route>
      <Route path="/admin/login" element={<LoginPage />} />
      <Route path="/admin" element={<DashboardPage />} />
      <Route path="/admin/hero" element={<HeroEditor />} />
      <Route path="/admin/about" element={<AboutEditor />} />
      <Route path="/admin/contact" element={<ContactEditor />} />
    </Routes>
  )
}

export default App
