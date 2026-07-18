import { useEffect, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { getNews, createNews, updateNews, deleteNews, uploadImage } from '../../services/api'

function NewsEditor() {
  const [posts, setPosts] = useState([])
  const [loading, setLoading] = useState(true)
  const [showForm, setShowForm] = useState(false)
  const [editingPost, setEditingPost] = useState(null)
  const [title, setTitle] = useState('')
  const [summary, setSummary] = useState('')
  const [content, setContent] = useState('')
  const [image, setImage] = useState('')
  const [saving, setSaving] = useState(false)
  const [toast, setToast] = useState(null)
  const navigate = useNavigate()

  useEffect(() => {
    const token = localStorage.getItem('token')
    if (!token) {
      navigate('/admin/login')
      return
    }
    fetchPosts()
  }, [navigate])

  const fetchPosts = () => {
    getNews()
      .then((res) => setPosts(res.data))
      .catch(console.error)
      .finally(() => setLoading(false))
  }

  const showToast = (message, type = 'success') => {
    setToast({ message, type })
    setTimeout(() => setToast(null), 3000)
  }

  const resetForm = () => {
    setTitle('')
    setSummary('')
    setContent('')
    setImage('')
    setEditingPost(null)
    setShowForm(false)
  }

  const handleEdit = (post) => {
    setEditingPost(post)
    setTitle(post.title)
    setSummary(post.summary || '')
    setContent(post.content)
    setImage(post.image || '')
    setShowForm(true)
  }

  const handleImageUpload = async (e) => {
    const file = e.target.files[0]
    if (!file) return
    try {
      const res = await uploadImage(file)
      setImage(res.data.url)
      showToast('Image uploaded successfully')
    } catch {
      showToast('Failed to upload image', 'error')
    }
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    setSaving(true)
    try {
      if (editingPost) {
        await updateNews(editingPost.id, { title, summary, content, image })
        showToast('Post updated successfully')
      } else {
        await createNews({ title, summary, content, image })
        showToast('Post created successfully')
      }
      resetForm()
      fetchPosts()
    } catch (err) {
      showToast(err.response?.data?.error || 'Failed to save post', 'error')
    } finally {
      setSaving(false)
    }
  }

  const handleDelete = async (id) => {
    if (!window.confirm('Are you sure you want to delete this post?')) return
    try {
      await deleteNews(id)
      showToast('Post deleted successfully')
      fetchPosts()
    } catch {
      showToast('Failed to delete post', 'error')
    }
  }

  if (loading) return <div className="loading">Loading...</div>

  return (
    <div className="admin-layout">
      <header className="admin-header">
        <h2>News Editor</h2>
        <nav>
          <Link to="/admin">Dashboard</Link>
          <Link to="/">View Site</Link>
        </nav>
      </header>
      <div className="admin-body">
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 20 }}>
          <h2>News Posts</h2>
          {!showForm && (
            <button className="btn btn-primary" onClick={() => setShowForm(true)}>
              + Create Post
            </button>
          )}
        </div>

        {showForm && (
          <div className="admin-card">
            <h3>{editingPost ? 'Edit Post' : 'New Post'}</h3>
            <form onSubmit={handleSubmit}>
              <div className="form-group">
                <label>Title</label>
                <input type="text" value={title} onChange={(e) => setTitle(e.target.value)} required />
              </div>
              <div className="form-group">
                <label>Summary</label>
                <textarea value={summary} onChange={(e) => setSummary(e.target.value)} placeholder="Brief description shown in the news list" style={{ minHeight: 60 }} />
              </div>
              <div className="form-group">
                <label>Content</label>
                <textarea value={content} onChange={(e) => setContent(e.target.value)} required placeholder="Full post content..." style={{ minHeight: 200 }} />
              </div>
              <div className="form-group">
                <label>Image</label>
                <input type="file" accept="image/*" onChange={handleImageUpload} />
                {image && (
                  <div className="image-preview">
                    <img src={image} alt="Preview" />
                  </div>
                )}
              </div>
              <div style={{ display: 'flex', gap: 10 }}>
                <button type="submit" className="btn btn-primary" disabled={saving}>
                  {saving ? 'Saving...' : editingPost ? 'Update Post' : 'Create Post'}
                </button>
                <button type="button" className="btn btn-secondary" onClick={resetForm}>
                  Cancel
                </button>
              </div>
            </form>
          </div>
        )}

        <div className="news-table-wrapper">
          {posts.length === 0 ? (
            <p style={{ textAlign: 'center', color: '#666', padding: 40 }}>No posts yet. Create your first post!</p>
          ) : (
            <table className="news-table">
              <thead>
                <tr>
                  <th>Title</th>
                  <th>Author</th>
                  <th>Date</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {posts.map((post) => (
                  <tr key={post.id}>
                    <td>{post.title}</td>
                    <td>{post.author}</td>
                    <td>{new Date(post.publishedAt).toLocaleDateString()}</td>
                    <td>
                      <button className="btn btn-primary btn-sm" onClick={() => handleEdit(post)}>
                        Edit
                      </button>
                      <button className="btn btn-danger btn-sm" onClick={() => handleDelete(post.id)}>
                        Delete
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      </div>

      {toast && (
        <div className={`toast ${toast.type === 'error' ? 'toast-error' : 'toast-success'}`}>
          {toast.message}
        </div>
      )}
    </div>
  )
}

export default NewsEditor
