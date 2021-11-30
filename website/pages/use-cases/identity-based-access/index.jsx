import SectionHeader from '@hashicorp/react-section-header'
import Button from '@hashicorp/react-button'
import TextSplits from '@hashicorp/react-text-splits'
import BeforeAfterDiagram from 'components/before-after-diagram'
import UseCaseCtaSection from 'components/use-case-cta-section'
//  Imports below are used in getStaticProps
import RAW_CONTENT from './content.json'
import highlightData from '@hashicorp/platform-code-highlighting/highlight-data'
import processBeforeAfterDiagramProps from 'components/before-after-diagram/server'

export async function getStaticProps() {
  const content = await highlightData(RAW_CONTENT)
  content.beforeAfterDiagram = await processBeforeAfterDiagramProps(
    content.beforeAfterDiagram
  )
  return { props: { content } }
}

export default function DataEncryptionUseCase({ content }) {
  return (
    <main id="use-cases" className="g-section-block page-wrap">
      {/* Header / Buttons */}
      <section className="g-grid-container">
        <SectionHeader
          headline="Leverage Trusted Identities in Low Trust Networks"
          description="Authenticate and access different clouds, systems, and endpoints using trusted identities"
          useH1={true}
        />

        <div className="button-container">
          <Button
            title="Download"
            url="/downloads"
            theme={{ brand: 'vault' }}
          />
          <Button
            title="Get Started"
            label="Get started — external link to education platform"
            url="/intro"
            theme="dark-outline"
          />
        </div>
      </section>

      {/* Before/After Diagram */}
      <section>
        <div className="g-grid-container">
          <BeforeAfterDiagram {...content.beforeAfterDiagram} />
        </div>
      </section>

      {/* Features */}
      <section className="no-section-spacing">
        <div className="g-grid-container">
          <SectionHeader headline="Identity-based Access Features" />
        </div>
        <TextSplits textSplits={content.features} />
      </section>

      <UseCaseCtaSection />
    </main>
  )
}
