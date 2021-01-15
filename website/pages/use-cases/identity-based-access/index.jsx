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
          headline="Leverage Trusted Identities in Low Trust Networks"
          description="Authenticate and access different clouds, systems, and endpoints using trusted identities"
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
              url: require('./img/challenge.png'),
              format: 'png',
            }}
            beforeHeadline="The Challenge"
            beforeContent="With the proliferation of different clouds, services, and systems all with their own identity providers, organizations need a way to manage identity sprawl"
            afterImage={{
              url: require('./img/solution.png'),
              format: 'png',
            }}
            afterHeadline="The Solution"
            afterContent="Vault merges identities across providers and uses a unified ACL system to broker access to systems and secrets"
          />
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
    </div>
  )
}
