import { useEffect, useState } from 'react'
import { getContent } from '../services/api'

function AboutPage() {
  const [content, setContent] = useState(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    getContent()
      .then((res) => setContent(res.data))
      .catch(console.error)
      .finally(() => setLoading(false))
  }, [])

  if (loading) return <div className="loading">Loading...</div>
  if (!content?.about?.visible) return null

  return (
    <section className="about-section">
      <div className="container">
        <h1>{content.about.title}</h1>
        <div className="about-content">
          {content.about.image && (
            <img src={content.about.image} alt="About us" className="about-image" />
          )}
          <p className="about-description">{content.about.description}</p>
        </div>
      </div>
    </section>
  )
}

export default AboutPage
