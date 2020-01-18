import SectionHeader from '@hashicorp/react-section-header'
import Button from '@hashicorp/react-button'
import TextAndContent from '@hashicorp/react-text-and-content'
import BeforeAfterDiagram from '../../../components/before-after-diagram'
import UseCaseCtaSection from '../../../components/use-case-cta-section'

export default function DataEncryptionUseCase() {
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
          <Button title="Download" url="/downloads.html" />
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
              format: 'svg'
            }}
            beforeHeadline="The Challenge"
            beforeContent="All application data should be encrypted, but deploying a cryptography and key management infrastructure is expensive, hard to develop against, and not cloud or multi-datacenter friendly"
            afterImage={{
              url:
                'https://www.datocms-assets.com/2885/1539885039-data-protectionsolution.svg',
              format: 'svg'
            }}
            afterHeadline="The Solution"
            afterContent="Vault provides encryption as a service with centralized key management to simplify encrypting data in transit and at rest across clouds and data centers"
          />
        </div>
      </section>

      {/* Features / Text and content */}
      <section className="g-container">
        <SectionHeader headline=" Encryption Features" />

        <TextAndContent
          data={{
            text: `### API-driven Encryption

Encrypt application data during transit and rest with AES 256-bit CBC data encryption and TLS in transit.`,
            content: {
              __typename: 'SbcImageRecord',
              image: {
                url: 'https://www.datocms-assets.com/2885/1539314348-eaas.png',
                format: 'png'
              }
            }
          }}
        />

        <TextAndContent
          data={{
            reverseDirection: true,
            text: `### Encryption Key Rolling

Update and roll new keys throughout distributed infrastructure while retaining the ability to decrypt encrypted data`,
            content: {
              __typename: 'SbcImageRecord',
              image: {
                url:
                  'https://www.datocms-assets.com/2885/1539314609-encryption-key-rolling.png',
                format: 'png'
              }
            }
          }}
        />
      </section>

      <UseCaseCtaSection />
    </div>
  )
}
