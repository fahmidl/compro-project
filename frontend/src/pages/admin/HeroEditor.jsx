import { useEffect, useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { getContent, updateContent, uploadImage } from '../../services/api'

function HeroEditor() {
  const [content, setContent] = useState(null)
  const [form, setForm] = useState({ title: '', subtitle: '', backgroundImage: '', visible: true })
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
          title: res.data.hero.title || '',
          subtitle: res.data.hero.subtitle || '',
          backgroundImage: res.data.hero.backgroundImage || '',
          visible: res.data.hero.visible,
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
      const updated = { ...content, hero: { ...form } }
      await updateContent(updated)
      showToast('Hero section saved!')
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
      setForm((prev) => ({ ...prev, backgroundImage: res.data.url }))
      showToast('Image uploaded!')
    } catch {
      showToast('Image upload failed', 'error')
    }
  }

  if (loading) return <div className="loading">Loading...</div>

  return (
    <div className="admin-layout">
      <header className="admin-header">
        <h2>Edit Hero Section</h2>
        <nav>
          <Link to="/admin">Dashboard</Link>
          <Link to="/">View Site</Link>
        </nav>
      </header>
      <div className="admin-body">
        <form onSubmit={handleSubmit} className="admin-card">
          <h3>Hero Content</h3>
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
            <label>Subtitle</label>
            <input
              type="text"
              value={form.subtitle}
              onChange={(e) => setForm({ ...form, subtitle: e.target.value })}
            />
          </div>
          <div className="form-group">
            <label>Background Image</label>
            <input type="file" accept="image/*" onChange={handleImageUpload} />
            {form.backgroundImage && (
              <div className="image-preview">
                <img src={form.backgroundImage} alt="Hero background" />
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

export default HeroEditor
