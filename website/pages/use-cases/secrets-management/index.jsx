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

export default function SecretsManagmentUseCase({ content }) {
  return (
    <div id="use-cases" className="g-section-block page-wrap">
      <section className="g-container">
        {/* Header / Buttons */}

        <SectionHeader
          headline="Secrets Management in Low Trust Networks"
          description="Centrally store, access, and deploy secrets across applications, systems, and infrastructure"
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
                'https://www.datocms-assets.com/2885/1539885048-secrets-managementchallenge.svg',
              format: 'svg',
            }}
            beforeHeadline="The Challenge"
            beforeContent="Secrets for applications and systems need to be centralized and static IP-based solutions don't scale in dynamic environments with frequently changing applications and machines"
            afterImage={{
              url:
                'https://www.datocms-assets.com/2885/1539885054-secrets-managementsolution.svg',
              format: 'svg',
            }}
            afterHeadline="The Solution"
            afterContent="Vault centrally manages and enforces access to secrets and systems based on trusted sources of application and user identity"
          />
        </div>
      </section>

      {/* Case study slider */}

      <section className="g-section-block theme-black-background-white-text">
        <div className="g-container">
          <div className="g-case-study-slider">
            <div className="case-study-container">
              <div className="slider-container">
                <div className="slider-frame single">
                  <div className="case-study">
                    <div className="feature-image">
                      <a href="https://www.hashicorp.com/resources/adobe-100-trillion-transactions-hashicorp-vault">
                        <picture>
                          <source
                            type="image/webp"
                            srcSet="https://www.datocms-assets.com/2885/1520367019-dan_mcteer_adobe_hashiconf2017.jpg?fit=crop&amp;fm=webp&amp;h=156.25&amp;q=80&amp;w=250 250w,https://www.datocms-assets.com/2885/1520367019-dan_mcteer_adobe_hashiconf2017.jpg?fit=crop&amp;fm=webp&amp;h=312.5&amp;q=80&amp;w=500 500w,https://www.datocms-assets.com/2885/1520367019-dan_mcteer_adobe_hashiconf2017.jpg?fit=crop&amp;fm=webp&amp;h=468.75&amp;q=80&amp;w=750 750w,https://www.datocms-assets.com/2885/1520367019-dan_mcteer_adobe_hashiconf2017.jpg?fit=crop&amp;fm=webp&amp;h=625&amp;q=80&amp;w=1000 1000w,https://www.datocms-assets.com/2885/1520367019-dan_mcteer_adobe_hashiconf2017.jpg?fit=crop&amp;fm=webp&amp;h=937.5&amp;q=80&amp;w=1500 1500w,https://www.datocms-assets.com/2885/1520367019-dan_mcteer_adobe_hashiconf2017.jpg?fit=crop&amp;fm=webp&amp;h=1250&amp;q=80&amp;w=2000 2000w,https://www.datocms-assets.com/2885/1520367019-dan_mcteer_adobe_hashiconf2017.jpg?fit=crop&amp;fm=webp&amp;h=1562.5&amp;q=80&amp;w=2500 2500w"
                            sizes="100vw"
                          />
                          <img
                            src="https://www.datocms-assets.com/2885/1520367019-dan_mcteer_adobe_hashiconf2017.jpg?fit=crop&amp;fm=jpg&amp;h=312.5&amp;q=80&amp;w=500"
                            srcSet="https://www.datocms-assets.com/2885/1520367019-dan_mcteer_adobe_hashiconf2017.jpg?fit=crop&amp;fm=jpg&amp;h=156.25&amp;q=80&amp;w=250 250w,https://www.datocms-assets.com/2885/1520367019-dan_mcteer_adobe_hashiconf2017.jpg?fit=crop&amp;fm=jpg&amp;h=312.5&amp;q=80&amp;w=500 500w,https://www.datocms-assets.com/2885/1520367019-dan_mcteer_adobe_hashiconf2017.jpg?fit=crop&amp;fm=jpg&amp;h=468.75&amp;q=80&amp;w=750 750w,https://www.datocms-assets.com/2885/1520367019-dan_mcteer_adobe_hashiconf2017.jpg?fit=crop&amp;fm=jpg&amp;h=625&amp;q=80&amp;w=1000 1000w,https://www.datocms-assets.com/2885/1520367019-dan_mcteer_adobe_hashiconf2017.jpg?fit=crop&amp;fm=jpg&amp;h=937.5&amp;q=80&amp;w=1500 1500w,https://www.datocms-assets.com/2885/1520367019-dan_mcteer_adobe_hashiconf2017.jpg?fit=crop&amp;fm=jpg&amp;h=1250&amp;q=80&amp;w=2000 2000w,https://www.datocms-assets.com/2885/1520367019-dan_mcteer_adobe_hashiconf2017.jpg?fit=crop&amp;fm=jpg&amp;h=1562.5&amp;q=80&amp;w=2500 2500w"
                            sizes="100vw"
                            alt="Dan McTeer at HashiConf 2017"
                          />
                        </picture>
                      </a>
                    </div>
                    <div className="feature-content">
                      <div className="single-logo">
                        <img
                          src="https://www.datocms-assets.com/2885/1539889072-1524097013-adobe-white-1.svg"
                          alt="Adobe logo"
                        />
                      </div>
                      <h3>
                        Using Vault to Protect Adobe&apos;s Secrets and User
                        Data Across Clouds and Datacenters
                      </h3>
                      <p>
                        Securing secrets and application data is a complex task
                        for globally distributed organizations. For Adobe,
                        managing secrets for over 20 products across 100,000
                        hosts, four regions, and trillions of transactions
                        annually requires a different approach altogether.
                      </p>
                      <a
                        className="g-btn primary-hashicorp-light"
                        href="https://www.hashicorp.com/resources/adobe-100-trillion-transactions-hashicorp-vault"
                      >
                        Read Case Study
                      </a>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Features */}
      <section className="no-section-spacing">
        <div className="g-grid-container">
          <SectionHeader headline="Secret Management Features" />
        </div>
        <TextSplits textSplits={content.features} />
      </section>

      <UseCaseCtaSection />
    </div>
  )
}
