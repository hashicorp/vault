import SectionHeader from '@hashicorp/react-section-header'
import Button from '@hashicorp/react-button'
import TextSplits from '@hashicorp/react-text-splits'
import BeforeAfterDiagram from 'components/before-after-diagram'
import UseCaseCtaSection from 'components/use-case-cta-section'
//  Imports below are used in getStaticProps
import RAW_CONTENT from './content.json'
import highlightData from '@hashicorp/platform-code-highlighting/highlight-data'
import processBeforeAfterDiagramProps from 'components/before-after-diagram/server'
import FeaturedSlider from '@hashicorp/react-featured-slider'

export async function getStaticProps() {
  const content = await highlightData(RAW_CONTENT)
  content.beforeAfterDiagram = await processBeforeAfterDiagramProps(
    content.beforeAfterDiagram
  )
  return { props: { content } }
}

export default function SecretsManagmentUseCase({ content }) {
  return (
    <main id="use-cases" className="g-section-block page-wrap">
      <section className="g-grid-container">
        {/* Header / Buttons */}

        <SectionHeader
          headline="Secrets Management in Low Trust Networks"
          description="Centrally store, access, and deploy secrets across applications, systems, and infrastructure"
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

      {/* Case study slider */}

      <FeaturedSlider
        theme="dark"
        features={[
          {
            logo: {
              url:
                'https://www.datocms-assets.com/2885/1539889072-1524097013-adobe-white-1.svg',
              alt: 'Adobe Logo',
            },
            image: {
              url:
                'https://www.datocms-assets.com/2885/1520367019-dan_mcteer_adobe_hashiconf2017.jpg?fit=crop&amp;fm=jpg&amp;h=312.5&amp;q=80&amp;w=500',
              alt: 'Dan McTeer at HashiConf 2017',
            },
            heading:
              "Using Vault to Protect Adobe's Secrets and User Data Across Clouds and Datacenters",
            content:
              'Securing secrets and application data is a complex task for globally distributed organizations. For Adobe, managing secrets for over 20 products across 100,000 hosts, four regions, and trillions of transactions annually requires a different approach altogether.',
            link: {
              text: 'Read Case Study',
              url:
                'https://www.hashicorp.com/resources/adobe-100-trillion-transactions-hashicorp-vault',
              type: 'outbound',
            },
          },
        ]}
      />

      {/* Features */}
      <section className="no-section-spacing">
        <div className="g-grid-container">
          <SectionHeader headline="Secret Management Features" />
        </div>
        <TextSplits textSplits={content.features} />
      </section>

      <UseCaseCtaSection />
    </main>
  )
}
