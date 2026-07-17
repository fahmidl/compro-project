import { useEffect, useState } from 'react'
import { getContent } from '../services/api'

function HomePage() {
  const [content, setContent] = useState(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    getContent()
      .then((res) => setContent(res.data))
      .catch(console.error)
      .finally(() => setLoading(false))
  }, [])

  if (loading) return <div className="loading">Loading...</div>
  if (!content?.hero?.visible) return null

  return (
    <section
      className="hero-section"
      style={
        content.hero.backgroundImage
          ? { backgroundImage: `url(${content.hero.backgroundImage})` }
          : {}
      }
    >
      <div className="hero-overlay">
        <h1>{content.hero.title}</h1>
        <p>{content.hero.subtitle}</p>
      </div>
    </section>
  )
}

export default HomePage
