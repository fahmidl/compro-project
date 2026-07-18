import { useEffect, useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { getNewsPost } from '../services/api'

function NewsDetailPage() {
  const { id } = useParams()
  const [post, setPost] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    getNewsPost(id)
      .then((res) => setPost(res.data))
      .catch(() => setError('Post not found'))
      .finally(() => setLoading(false))
  }, [id])

  if (loading) return <div className="loading">Loading...</div>
  if (error) {
    return (
      <section className="news-detail-section">
        <div className="container">
          <p style={{ textAlign: 'center', color: '#666' }}>{error}</p>
          <p style={{ textAlign: 'center' }}>
            <Link to="/news">&larr; Back to News</Link>
          </p>
        </div>
      </section>
    )
  }
  if (!post) return null

  return (
    <section className="news-detail-section">
      <div className="container">
        <Link to="/news" className="news-back-link">&larr; Back to News</Link>
        <article className="news-detail">
          {post.image && (
            <div className="news-detail-image">
              <img src={post.image} alt={post.title} />
            </div>
          )}
          <h1>{post.title}</h1>
          <p className="news-detail-meta">
            {new Date(post.publishedAt).toLocaleDateString('en-US', {
              year: 'numeric',
              month: 'long',
              day: 'numeric',
            })}{' '}
            &middot; {post.author}
          </p>
          <div className="news-detail-content">
            {post.content.split('\n').map((paragraph, i) => (
              <p key={i}>{paragraph}</p>
            ))}
          </div>
        </article>
      </div>
    </section>
  )
}

export default NewsDetailPage
