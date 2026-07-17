import { useEffect, useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { getContent, updateContent, uploadImage } from '../../services/api'

function AboutEditor() {
  const [content, setContent] = useState(null)
  const [form, setForm] = useState({ title: '', description: '', image: '', visible: true })
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
          title: res.data.about.title || '',
          description: res.data.about.description || '',
          image: res.data.about.image || '',
          visible: res.data.about.visible,
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
      const updated = { ...content, about: { ...form } }
      await updateContent(updated)
      showToast('About section saved!')
    } catch {
      showToast('Failed to save', 'error')
    } finally {
      setSaving(false)
    }
  }

  const handleImageUpload = async (e) => {
    const file = e.target.files[0]
    if (!file) return
    try {
      const res = await uploadImage(file)
      setForm((prev) => ({ ...prev, image: res.data.url }))
      showToast('Image uploaded!')
    } catch {
      showToast('Image upload failed', 'error')
    }
  }

  if (loading) return <div className="loading">Loading...</div>

  return (
    <div className="admin-layout">
      <header className="admin-header">
        <h2>Edit About Section</h2>
        <nav>
          <Link to="/admin">Dashboard</Link>
          <Link to="/">View Site</Link>
        </nav>
      </header>
      <div className="admin-body">
        <form onSubmit={handleSubmit} className="admin-card">
          <h3>About Content</h3>
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
            <label>Description</label>
            <textarea
              value={form.description}
              onChange={(e) => setForm({ ...form, description: e.target.value })}
            />
          </div>
          <div className="form-group">
            <label>Image</label>
            <input type="file" accept="image/*" onChange={handleImageUpload} />
            {form.image && (
              <div className="image-preview">
                <img src={form.image} alt="About" />
              </div>
            )}
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

export default AboutEditor
