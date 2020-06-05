import Button from '@hashicorp/react-button'

export default function UseCaseCtaSection() {
  return (
    <section className="g-section-block g-cta-section">
      <div>
        <h2>Ready to get started?</h2>
        <Button
          url="/downloads"
          title="Download"
          linkType="download"
          theme={{
            variant: 'primary',
            background: 'dark',
            brand: 'neutral'
          }}
        />
        <Button
          url="/docs"
          title="Explore Docs"
          theme={{ variant: 'secondary', background: 'dark' }}
        />
      </div>
    </section>
  )
}
