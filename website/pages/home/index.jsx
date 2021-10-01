import HomepageHero from 'components/homepage-hero'
import SectionHeader from '@hashicorp/react-section-header'
import UseCases from '@hashicorp/react-use-cases'
import TextSplits from '@hashicorp/react-text-splits'
import Button from '@hashicorp/react-button'
import BeforeAfterDiagram from '../../components/before-after-diagram'
import HcpCalloutSection from 'components/hcp-callout-section'
//  Imports below are used in getStaticProps only
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

export default function Homepage({ content }) {
  return (
    <main id="page-home">
      <div className="g-section-block page-wrap">
        <HomepageHero
          uiVideo="https://www.datocms-assets.com/2885/1543956852-vault-v1-0-ui-opt.mp4"
          cliVideo="https://www.datocms-assets.com/2885/1543956847-vault-v1-0-cli-opt.mp4"
          title="Manage Secrets and Protect Sensitive Data"
          description="Secure, store and tightly control access to tokens, passwords, certificates, encryption keys for protecting secrets and other sensitive data using a UI, CLI, or HTTP API."
          buttons={[
            {
              external: false,
              title: 'Try Cloud',
              url:
                'https://portal.cloud.hashicorp.com/sign-up?utm_source=vault_io&utm_content=hero',
            },
            {
              external: false,
              title: 'Download CLI',
              url: 'https://www.vaultproject.io/downloads',
            },
            {
              type: 'inbound',
              title: 'Get Started with Vault',
              url: 'https://www.vaultproject.io/intro/getting-started',
              theme: { variant: 'tertiary' },
            },
          ]}
        />

        {/* Text Section */}

        <section className="g-grid-container remove-bottom-padding">
          <SectionHeader
            headline="Secure dynamic infrastructure across clouds and environments"
            description="The shift from static, on-premise infrastructure to dynamic, multi-provider infrastructure changes the approach to security. Security in static infrastructure relies on dedicated servers, static IP addresses, and a clear network perimeter. Security in dynamic infrastructure is defined by ephemeral applications and servers, trusted sources of user and application identity, and software-based encryption."
          />
        </section>

        {/* Before-After Diagram */}

        <section className="g-grid-container before-after">
          <BeforeAfterDiagram
            {...content.beforeAfterDiagram}
            beforeImage={{
              alt: 'Static database graphic',
              format: 'png',
              url: require('./img/vault_static_isometric@2x.png'),
            }}
            afterImage={{
              alt: 'Dynamic VM and database graphic',
              format: 'png',
              url: require('./img/vault_dynamic_isometric@2x.png'),
            }}
          />
        </section>

        {/* Use cases */}

        <section>
          <div className="g-grid-container">
            <UseCases
              product="vault"
              items={[
                {
                  title: 'Secrets Management',
                  description:
                    'Centrally store, access, and deploy secrets across applications, systems, and infrastructure',
                  image: {
                    alt: 'Key icon',
                    format: 'png',
                    url: require('./img/use-cases/secrets-management.svg?url'),
                  },
                  link: {
                    external: false,
                    title: 'Learn more',
                    url: '/use-cases/secrets-management',
                  },
                },
                {
                  title: 'Data Encryption',
                  description:
                    'Keep secrets and application data secure with one centralized workflow to encrypt data in flight and at rest',
                  image: {
                    alt: 'Lock icon',
                    format: 'png',
                    url: require('./img/use-cases/data_encryption.svg?url'),
                  },
                  link: {
                    external: false,
                    title: 'Learn more',
                    url: '/use-cases/data-encryption',
                  },
                },
                {
                  title: 'Identity-based Access',
                  description:
                    'Authenticate and access different clouds, systems, and endpoints using trusted identities',
                  image: {
                    alt: 'Access badge icon',
                    format: 'png',
                    url: require('./img/use-cases/identity-based-access.svg?url'),
                  },
                  link: {
                    external: false,
                    title: 'Learn more',
                    url: '/use-cases/identity-based-access',
                  },
                },
              ]}
            />
          </div>
        </section>

        <HcpCalloutSection
          id="cloud-offerings"
          title="HCP Vault"
          chin="Available on AWS"
          description="HCP Vault provides all of the power and security of Vault, without the complexity and overhead of managing it yourself. Access Vault’s best-in-class secrets management and encryption capabilities instantly and onboard applications and teams easily."
          image={require('./img/hcp_vault.svg?url')}
          links={[
            {
              text: 'Learn More',
              url:
                'https://cloud.hashicorp.com/?utm_source=vault_io&utm_content=hcp_vault_detail',
            },
          ]}
        />

        {/* Principles / Text & Content Blocks */}
        <section className="no-section-spacing">
          <div className="g-grid-container">
            <SectionHeader headline="Vault Principles" />
          </div>
          <TextSplits textSplits={content.principles} />
        </section>

        <section className="g-grid-container">
          <SectionHeader
            headline="Open Source and Enterprise"
            description="Vault Open Source addresses the technical complexity of managing secrets by leveraging trusted identities across distributed infrastructure and clouds. Vault Enterprise addresses the organizational complexity of large user bases and compliance requirements with collaboration and governance features."
          />
          <div className="button-container">
            <Button
              title="Learn More"
              label="Learn more — Vault pricing tiers"
              url="https://www.hashicorp.com/products/vault/enterprise"
              theme={{ brand: 'vault' }}
            />
          </div>
        </section>
      </div>
    </main>
  )
}
