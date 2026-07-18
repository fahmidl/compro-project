import { useEffect, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { listUsers, createUser, deleteUser } from '../../services/api'

function UserManagement() {
  const [users, setUsers] = useState([])
  const [loading, setLoading] = useState(true)
  const [showForm, setShowForm] = useState(false)
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [role, setRole] = useState('editor')
  const [saving, setSaving] = useState(false)
  const [toast, setToast] = useState(null)
  const navigate = useNavigate()

  useEffect(() => {
    const token = localStorage.getItem('token')
    const userRole = localStorage.getItem('role')
    if (!token || userRole !== 'admin') {
      navigate('/admin')
      return
    }
    fetchUsers()
  }, [navigate])

  const fetchUsers = () => {
    listUsers()
      .then((res) => setUsers(res.data))
      .catch(console.error)
      .finally(() => setLoading(false))
  }

  const showToast = (message, type = 'success') => {
    setToast({ message, type })
    setTimeout(() => setToast(null), 3000)
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    setSaving(true)
    try {
      await createUser({ username, password, role })
      showToast('User created successfully')
      setUsername('')
      setPassword('')
      setRole('editor')
      setShowForm(false)
      fetchUsers()
    } catch (err) {
      showToast(err.response?.data?.error || 'Failed to create user', 'error')
    } finally {
      setSaving(false)
    }
  }

  const handleDelete = async (id, name) => {
    if (!window.confirm(`Are you sure you want to delete user "${name}"?`)) return
    try {
      await deleteUser(id)
      showToast('User deleted successfully')
      fetchUsers()
    } catch (err) {
      showToast(err.response?.data?.error || 'Failed to delete user', 'error')
    }
  }

  if (loading) return <div className="loading">Loading...</div>

  return (
    <div className="admin-layout">
      <header className="admin-header">
        <h2>User Management</h2>
        <nav>
          <Link to="/admin">Dashboard</Link>
          <Link to="/">View Site</Link>
        </nav>
      </header>
      <div className="admin-body">
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 20 }}>
          <h2>Users</h2>
          {!showForm && (
            <button className="btn btn-primary" onClick={() => setShowForm(true)}>
              + Add User
            </button>
          )}
        </div>

        {showForm && (
          <div className="admin-card">
            <h3>Create New User</h3>
            <form onSubmit={handleSubmit}>
              <div className="form-group">
                <label>Username</label>
                <input type="text" value={username} onChange={(e) => setUsername(e.target.value)} required />
              </div>
              <div className="form-group">
                <label>Password</label>
                <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} required minLength={6} />
              </div>
              <div className="form-group">
                <label>Role</label>
                <select value={role} onChange={(e) => setRole(e.target.value)} className="form-select">
                  <option value="editor">Editor</option>
                  <option value="admin">Admin</option>
                </select>
              </div>
              <div style={{ display: 'flex', gap: 10 }}>
                <button type="submit" className="btn btn-primary" disabled={saving}>
                  {saving ? 'Creating...' : 'Create User'}
                </button>
                <button type="button" className="btn btn-secondary" onClick={() => setShowForm(false)}>
                  Cancel
                </button>
              </div>
            </form>
          </div>
        )}

        <div className="news-table-wrapper">
          <table className="news-table">
            <thead>
              <tr>
                <th>Username</th>
                <th>Role</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {users.map((user) => (
                <tr key={user.id}>
                  <td>{user.username}</td>
                  <td>
                    <span className={`role-badge role-${user.role || 'editor'}`}>
                      {user.role || 'editor'}
                    </span>
                  </td>
                  <td>
                    <button className="btn btn-danger btn-sm" onClick={() => handleDelete(user.id, user.username)}>
                      Delete
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
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

export default UserManagement
