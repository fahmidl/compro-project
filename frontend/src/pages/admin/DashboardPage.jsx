import { useEffect, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { getContent } from '../../services/api'

function DashboardPage() {
  const [content, setContent] = useState(null)
  const [loading, setLoading] = useState(true)
  const navigate = useNavigate()
  const [role, setRole] = useState('')

  useEffect(() => {
    const token = localStorage.getItem('token')
    if (!token) {
      navigate('/admin/login')
      return
    }
    setRole(localStorage.getItem('role') || '')

    getContent()
      .then((res) => setContent(res.data))
      .catch(console.error)
      .finally(() => setLoading(false))
  }, [navigate])

  const handleLogout = () => {
    localStorage.removeItem('token')
    localStorage.removeItem('role')
    navigate('/admin/login')
  }

  if (loading) return <div className="loading">Loading...</div>

  return (
    <div className="admin-layout">
      <header className="admin-header">
        <h2>CMS Dashboard</h2>
        <nav>
          <Link to="/">View Site</Link>
          <button className="btn btn-danger" onClick={handleLogout}>
            Logout
          </button>
        </nav>
      </header>
      <div className="admin-body">
        <h2 style={{ marginBottom: 20 }}>Manage Sections</h2>
        <div className="dashboard-grid">
          <div className="dashboard-card">
            <h3>Hero Section</h3>
            <p>
              Status: {content?.hero?.visible ? '✅ Visible' : '❌ Hidden'}
            </p>
            <Link to="/admin/hero" className="btn btn-primary">
              Edit Hero
            </Link>
          </div>
          <div className="dashboard-card">
            <h3>About Section</h3>
            <p>
              Status: {content?.about?.visible ? '✅ Visible' : '❌ Hidden'}
            </p>
            <Link to="/admin/about" className="btn btn-primary">
              Edit About
            </Link>
          </div>
          <div className="dashboard-card">
            <h3>Contact Section</h3>
            <p>
              Status: {content?.contact?.visible ? '✅ Visible' : '❌ Hidden'}
            </p>
            <Link to="/admin/contact" className="btn btn-primary">
              Edit Contact
            </Link>
          </div>
          <div className="dashboard-card">
            <h3>News</h3>
            <p>Manage news posts</p>
            <Link to="/admin/news" className="btn btn-primary">
              Manage News
            </Link>
          </div>
          {role === 'admin' && (
            <div className="dashboard-card">
              <h3>Users</h3>
              <p>Manage users and roles</p>
              <Link to="/admin/users" className="btn btn-primary">
                Manage Users
              </Link>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

export default DashboardPage
