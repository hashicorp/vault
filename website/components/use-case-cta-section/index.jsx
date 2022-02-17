import Button from '@hashicorp/react-button'

export default function UseCaseCtaSection() {
  return (
    <section className="g-section-block g-cta-section">
      <div>
        <h2 className="g-type-display-2">Ready to get started?</h2>
        <Button
          url="/downloads"
          title="Download"
          label="Download CLI"
          linkType="download"
          className="g-btn"
          theme={{
            variant: 'primary',
            brand: 'neutral',
          }}
        />
        <Button
          url="/docs"
          title="Explore Docs"
          className="g-btn"
          theme={{ variant: 'secondary' }}
        />
      </div>
    </section>
  )
}
