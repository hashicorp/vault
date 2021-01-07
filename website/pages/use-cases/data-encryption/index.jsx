import SectionHeader from '@hashicorp/react-section-header'
import Button from '@hashicorp/react-button'
import TextSplits from '@hashicorp/react-text-splits'
import BeforeAfterDiagram from 'components/before-after-diagram'
import UseCaseCtaSection from 'components/use-case-cta-section'
//  Imports below are used in getStaticProps
import RAW_CONTENT from './content.json'
import highlightData from '@hashicorp/nextjs-scripts/prism/highlight-data'

export async function getStaticProps() {
  const content = await highlightData(RAW_CONTENT)
  return { props: { content } }
}

export default function DataEncryptionUseCase({ content }) {
  return (
    <div id="use-cases" className="g-section-block page-wrap">
      {/* Header / Buttons */}
      <section className="g-container">
        <SectionHeader
          headline="Encrypt Application Data in Low Trust Networks"
          description="Keep application data secure with one centralized workflow to encrypt data in flight and at rest"
          useH1={true}
        />

        <div className="button-container">
          <Button title="Download" url="/downloads" />
          <Button title="Get Started" url="/intro" theme="dark-outline" />
        </div>
      </section>

      {/* Before/After Diagram */}
      <section>
        <div className="g-container">
          <BeforeAfterDiagram
            beforeImage={{
              url:
                'https://www.datocms-assets.com/2885/1539885046-data-protectionchallenge.svg',
              format: 'svg',
            }}
            beforeHeadline="The Challenge"
            beforeContent="All application data should be encrypted, but deploying a cryptography and key management infrastructure is expensive, hard to develop against, and not cloud or multi-datacenter friendly"
            afterImage={{
              url:
                'https://www.datocms-assets.com/2885/1539885039-data-protectionsolution.svg',
              format: 'svg',
            }}
            afterHeadline="The Solution"
            afterContent="Vault provides encryption as a service with centralized key management to simplify encrypting data in transit and at rest across clouds and data centers"
          />
        </div>
      </section>

      {/* Features */}
      <section className="no-section-spacing">
        <div className="g-grid-container">
          <SectionHeader headline=" Encryption Features" />
        </div>
        <TextSplits textSplits={content.features} />
      </section>

      <UseCaseCtaSection />
    </div>
  )
}
