import Hero from 'components/hero'
import RAW_CONTENT from './content.json'

export async function getStaticProps() {
  return { props: {} }
}

export default function Homepage({ content }) {
  return (
    <>
      <Hero
        brand="vault"
        heading="Manage Secrets &amp; Protect Sensitive Data"
        description="Secure, store and tightly control access to tokens, passwords, certificates, encryption keys for protecting secrets and other sensitive data using a UI, CLI, or HTTP API."
        ctas={[
          {
            title: 'View tutorials',
            url: '#TODO',
          },
          {
            title: 'View documentation',
            url: '#TODO',
          },
        ]}
        cards={[
          {
            heading: 'Open Source',
            description: 'Self-managed | always free',
            cta: {
              title: 'Download',
              url: '#TODO',
            },
            subText:
              'Download the open source Vault binary and run locally or within your environments.',
          },
          {
            heading: 'Cloud',
            description: 'Managed Vault',
            cta: {
              title: 'Get started for free',
              url: '#TODO',
            },
            subText:
              'Get up and running in minutes with a fully managed Vault cluster on HCP (HashiCorp Cloud Platform)',
          },
        ]}
      />
    </>
  )
}
