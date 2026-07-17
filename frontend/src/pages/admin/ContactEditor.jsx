import { useEffect, useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { getContent, updateContent } from '../../services/api'

function ContactEditor() {
  const [content, setContent] = useState(null)
  const [form, setForm] = useState({
    title: '',
    address: '',
    phone: '',
    email: '',
    mapEmbedUrl: '',
    visible: true,
  })
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [toast, setToast] = useState(null)
  const navigate = useNavigate()

  useEffect(() => {
    const token = localStorage.getItem('token')
    if (!token) { navigate('/admin/login'); return }

    getContent()
      .then((res) => {
        setContent(res.data)
        setForm({
          title: res.data.contact.title || '',
          address: res.data.contact.address || '',
          phone: res.data.contact.phone || '',
          email: res.data.contact.email || '',
          mapEmbedUrl: res.data.contact.mapEmbedUrl || '',
          visible: res.data.contact.visible,
        })
      })
      .catch(console.error)
      .finally(() => setLoading(false))
  }, [navigate])

  const showToast = (message, type = 'success') => {
    setToast({ message, type })
    setTimeout(() => setToast(null), 3000)
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    setSaving(true)
    try {
      const updated = { ...content, contact: { ...form } }
      await updateContent(updated)
      showToast('Contact section saved!')
    } catch {
      showToast('Failed to save', 'error')
    } finally {
      setSaving(false)
    }
  }

  if (loading) return <div className="loading">Loading...</div>

  return (
    <div className="admin-layout">
      <header className="admin-header">
        <h2>Edit Contact Section</h2>
        <nav>
          <Link to="/admin">Dashboard</Link>
          <Link to="/">View Site</Link>
        </nav>
      </header>
      <div className="admin-body">
        <form onSubmit={handleSubmit} className="admin-card">
          <h3>Contact Content</h3>
          <div className="form-group">
            <label>Title</label>
            <input
              type="text"
              value={form.title}
              onChange={(e) => setForm({ ...form, title: e.target.value })}
              required
            />
          </div>
          <div className="form-group">
            <label>Address</label>
            <input
              type="text"
              value={form.address}
              onChange={(e) => setForm({ ...form, address: e.target.value })}
            />
          </div>
          <div className="form-group">
            <label>Phone</label>
            <input
              type="text"
              value={form.phone}
              onChange={(e) => setForm({ ...form, phone: e.target.value })}
            />
          </div>
          <div className="form-group">
            <label>Email</label>
            <input
              type="email"
              value={form.email}
              onChange={(e) => setForm({ ...form, email: e.target.value })}
            />
          </div>
          <div className="form-group">
            <label>Map Embed URL</label>
            <input
              type="text"
              value={form.mapEmbedUrl}
              onChange={(e) => setForm({ ...form, mapEmbedUrl: e.target.value })}
              placeholder="https://www.google.com/maps/embed?pb=..."
            />
          </div>
          <div className="form-group checkbox-group">
            <input
              type="checkbox"
              checked={form.visible}
              onChange={(e) => setForm({ ...form, visible: e.target.checked })}
            />
            <label>Visible on site</label>
          </div>
          <button type="submit" className="btn btn-primary" disabled={saving}>
            {saving ? 'Saving...' : 'Save Changes'}
          </button>
        </form>
      </div>
      {toast && <div className={`toast toast-${toast.type}`}>{toast.message}</div>}
    </div>
  )
}

export default ContactEditor
