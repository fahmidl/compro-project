import { useEffect, useState } from 'react'
import { getContent } from '../services/api'

function ContactPage() {
  const [content, setContent] = useState(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    getContent()
      .then((res) => setContent(res.data))
      .catch(console.error)
      .finally(() => setLoading(false))
  }, [])

  if (loading) return <div className="loading">Loading...</div>
  if (!content?.contact?.visible) return null

  return (
    <section className="contact-section">
      <div className="container">
        <h1>{content.contact.title}</h1>
        <div className="contact-info">
          <div className="contact-details">
            <div className="contact-item">
              <strong>Address:</strong>
              <p>{content.contact.address}</p>
            </div>
            <div className="contact-item">
              <strong>Phone:</strong>
              <p>{content.contact.phone}</p>
            </div>
            <div className="contact-item">
              <strong>Email:</strong>
              <p>
                <a href={`mailto:${content.contact.email}`}>{content.contact.email}</a>
              </p>
            </div>
          </div>
          {content.contact.mapEmbedUrl && (
            <div className="contact-map">
              <iframe
                src={content.contact.mapEmbedUrl}
                width="100%"
                height="300"
                style={{ border: 0 }}
                allowFullScreen=""
                loading="lazy"
                title="Location Map"
              ></iframe>
            </div>
          )}
        </div>
      </div>
    </section>
  )
}

export default ContactPage
