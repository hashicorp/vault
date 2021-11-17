import * as React from 'react'
import IoHomeHero from 'components/io-home-hero'
import IoHomeVideoCallout from 'components/io-home-video-callout'
import IoHomeCard from 'components/io-home-card'
import IoHomeCaseStudies from 'components/io-home-case-studies'
import IoHomeCallToAction from 'components/io-home-call-to-action'
import IoHomePreFooter from 'components/io-home-pre-footer'
import s from './style.module.css'

export default function Homepage({ content }) {
  return (
    <>
      <IoHomeHero
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

      <div className={s.intro}>
        <div className={s.container}>
          <IoHomeVideoCallout
            heading="How Vault works"
            description="Vault tightly controls access to secrets and encryption keys by authenticating against trusted sources of identity such as Active Directory, LDAP, Kubernetes, CloudFoundry, and cloud platforms."
            thumbnail="/img/TEMP-thumbnail.png"
            person={{
              name: 'Armon Dadgar',
              description: 'Co-founder & CTO',
              thumbnail: '/img/TEMP-thumbnail.png',
            }}
          />
        </div>
      </div>

      <div className={s.inPractice}>
        <div className={s.container}>
          <header className={s.inPracticeHeader}>
            <h2 className={s.inPracticeHeading}>Vault in practice</h2>
            <p className={s.inPracticeDescription}>
              The best way to understand what Vault can enable for your projects
              is to see it in action
            </p>
          </header>
          <ul className={s.inPracticeCards}>
            <li>
              <IoHomeCard
                variant="dark"
                link={{
                  url: '/',
                  type: 'outbound',
                }}
                eyebrow="Documentation"
                heading="Secrets storage"
                description="Securely store and manage access to secrets and systems based on trusted sources of application and user identity."
              />
            </li>

            <li>
              <IoHomeCard
                variant="dark"
                link={{
                  url: '/',
                  type: 'outbound',
                }}
                eyebrow="Tutorial"
                heading="Dynamic secrets"
                description="Eliminate secret sprawl and reduce exposure risk with dynamically created and destroyed unique, on-demand credentials."
              />
            </li>

            <li>
              <IoHomeCard
                variant="dark"
                link={{
                  url: '/',
                  type: 'outbound',
                }}
                eyebrow="Tutorial"
                heading="Automate credential rotation"
                description="Reduce risk of secret exposure by automating how long secrets live and rotating secrets across your entire fleet."
              />
            </li>

            <li>
              <IoHomeCard
                variant="dark"
                link={{
                  url: '/',
                  type: 'outbound',
                }}
                eyebrow="Documentation"
                heading="Encryption as a service"
                description="Vault provides encryption as a service to simplify encrypting data, tokenizing sensitive values, signing and validating."
              />
            </li>

            <li>
              <IoHomeCard
                variant="dark"
                link={{
                  url: '/',
                  type: 'outbound',
                }}
                eyebrow="Tutorial"
                heading="API-driven encryption"
                description="Vault provides rich APIs to protect data, while using the state of the art in cryptography. Vault uses smart defaults but also."
              />
            </li>

            <li>
              <IoHomeCard
                variant="dark"
                link={{
                  url: '/',
                  type: 'outbound',
                }}
                eyebrow="Tutorial"
                heading="Encryption key rolling"
                description="Automatically update and rotate encryption keys without code changes, configuration updates, or re-deploys. Developers use a."
              />
            </li>
          </ul>
        </div>
      </div>

      <div className={s.caseStudies}>
        <div className={s.container}>
          <header className={s.caseStudiesHeader}>
            <h2 className={s.caseStudiesHeading}>Vault case studies</h2>
            <p className={s.caseStudiesDescription}>
              An inside look at powerful solutions from some of the world’s most
              innovative companies.
            </p>
          </header>

          <IoHomeCaseStudies
            primary={[
              {
                link: '',
                thumbnail: '/img/TEMP-thumbnail.png',
                alt: 'Sample alt text',
                heading: 'Accelerating the path to modern banking',
              },
              {
                link: '',
                thumbnail: '/img/TEMP-thumbnail.png',
                alt: 'Sample alt text',
                heading: 'A unified plan for secrets',
              },
            ]}
            secondary={[
              {
                link: '',
                heading:
                  'Uber Hadoop Cluster Process Secured with HashiCorp Vault',
              },
              {
                link: '',
                heading: 'Terraform for the lean engineering team at Compile ',
              },
              {
                link: '',
                heading: 'Seeding HashiCorp Vault with Terraform at Form3',
              },
              {
                link: '',
                heading: 'A G-Research story: 1 to 1000 Vault namespaces',
              },
            ]}
          />
        </div>
      </div>

      <IoHomeCallToAction
        brand="vault"
        heading="Get HashiCorp Certified"
        content="Level up your concepts, skills, and use cases associated with open source HashiCorp Vault."
        links={[
          { text: 'Prepare & get certified', url: '#TODO' },
          {
            text: 'Learn more about Vault Ops Pro',
            url: '#TODO',
            type: 'inbound',
          },
        ]}
      />

      <IoHomePreFooter
        brand="vault"
        heading="Next steps"
        description="HCP Vault simplifies cloud security automation on fully managed infrastructure. Get started for free, and pay only for what you use."
        ctas={[
          {
            link: '#TODO',
            heading: 'Open Source',
            description: 'Self-managed | always free',
            label: 'Download',
          },
          {
            link: '#TODO',
            heading: 'Cloud',
            description: 'Self-managed',
            label: 'Compare plans',
          },
          {
            link: '#TODO',
            heading: 'Enterprise',
            description: 'Self-Managed custom deployments',
            label: 'Learn more',
          },
        ]}
      />
    </>
  )
}
