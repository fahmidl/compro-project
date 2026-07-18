import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { getNews } from '../services/api'

function NewsPage() {
  const [posts, setPosts] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    getNews()
      .then((res) => setPosts(res.data))
      .catch(console.error)
      .finally(() => setLoading(false))
  }, [])

  if (loading) return <div className="loading">Loading...</div>

  return (
    <section className="news-section">
      <div className="container">
        <h1>Latest News</h1>
        {posts.length === 0 ? (
          <p className="news-empty">No news posts yet.</p>
        ) : (
          <div className="news-grid">
            {posts.map((post) => (
              <Link to={`/news/${post.slug || post.id}`} key={post.id} className="news-card">
                {post.image && (
                  <div className="news-card-image">
                    <img src={post.image} alt={post.title} />
                  </div>
                )}
                <div className="news-card-body">
                  <h2>{post.title}</h2>
                  <p className="news-card-meta">
                    {new Date(post.publishedAt).toLocaleDateString('en-US', {
                      year: 'numeric',
                      month: 'long',
                      day: 'numeric',
                    })}{' '}
                    &middot; {post.author}
                  </p>
                  <p className="news-card-summary">{post.summary || post.content.substring(0, 150) + '...'}</p>
                  <span className="news-card-link">Read more &rarr;</span>
                </div>
              </Link>
            ))}
          </div>
        )}
      </div>
    </section>
  )
}

export default NewsPage
