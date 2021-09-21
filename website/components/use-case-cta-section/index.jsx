import Button from '@hashicorp/react-button'

export default function UseCaseCtaSection() {
  return (
    <section className="g-section-block g-cta-section">
      <div>
        <h2 className='g-type-display-2'>Ready to get started?</h2>
        <Button
          url="/downloads"
          title="Download"
          label="Download CLI"
          linkType="download"
          theme={{
            variant: 'primary',
            brand: 'neutral',
          }}
        />
        <Button
          url="/docs"
          title="Explore Docs"
          theme={{ variant: 'secondary' }}
        />
      </div>
    </section>
  )
}
