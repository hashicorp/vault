import IoUsecaseHero from 'components/io-usecase-hero'
import IoUsecaseSection from 'components/io-usecase-section'
import IoUsecaseCustomer from 'components/io-usecase-customer'
// import IoCard from 'components/io-card'
import IoVideoCallout from 'components/io-video-callout'
import IoUsecaseCallToAction from 'components/io-usecase-call-to-action'
import s from './style.module.css'
import React from 'react'

export default function UseCases(): React.ReactElement {
  return (
    <>
      <IoUsecaseHero
        eyebrow="Common use case"
        heading="Credential rotation with Vault"
        description="Eliminate long standing shared credentials and reduce risk of breach and credential leakage with automated database credential rotation"
      />

      <IoUsecaseSection
        brand="vault"
        eyebrow="Challenge"
        heading="Each database your organization uses requires credentials for application, services, and users to access or use the data."
        description={
          <p>
            This creates potentially thousands of consumers that need access to
            one or more databases. Safeguarding and ensuring that one of these
            credentials isn’t leaked, or in the likelihood it is, that the
            organization can quickly revoke access and remediate, is a complex
            problem to solve.
          </p>
        }
        cta={{
          text: 'Learn more',
          link: '/',
        }}
      />

      <IoUsecaseSection
        brand="vault"
        eyebrow="Solution"
        heading="Create, rotate, and revoke database credentials through an automated workflow and API."
        media={{
          src: '/img/TEMP-customer-story.png',
          width: '592',
          height: '455',
          alt: '',
        }}
        description={
          <>
            <p>
              This allows each application, service, or user to dynamically get
              unique credentials to access the database(s) as well as lease and
              expiration times for the credentials. This means that the
              credentials will expire and reduce impact of breached from leaked
              credentials.
            </p>
            <p>
              In a scenario where credentials are lost or stolen, the window for
              those credentials to be valid can be reduced to almost nothing or
              instant-use only. If credentials are stolen or leaked, the same
              automated workflow for issuance and rotation can also
              automatically revoke access and seal Vault and lock down from
              outside access.
            </p>
          </>
        }
        cta={{
          text: 'Learn more',
          link: '/',
        }}
      />

      <IoUsecaseCustomer
        link="/"
        media={{
          src: '/img/TEMP-customer-story.png',
          width: '592',
          height: '455',
          alt: '',
        }}
        logo={{
          src: '/img/TEMP-logo.svg',
          width: '89',
          height: '25',
          alt: '',
        }}
        heading="Message received"
        description="Leading global advertising platform uses HashiCorp Consul to launch new services in <1 minute by eliminating all manual operations"
        stats={[
          {
            value: '150k+',
            key: 'Tokens and credentials generated per day',
          },
          {
            value: '10m+',
            key: 'Customers served per year',
          },
          {
            value: '25+',
            key: 'Manual work hours saved per week',
          },
        ]}
      />

      <div className={s.callToAction}>
        <IoUsecaseCallToAction
          theme="light"
          brand="vault"
          heading="Get started HCP Vault and secrets injection into Kubernetes"
          description="We’ve built a step-by-step guide on integrating HCP Vault and your Kubernetes cluster."
          links={[
            {
              text: 'Try HCP Vault for free',
              url: '/',
            },
          ]}
        />
      </div>

      <div className={s.videoCallout}>
        <IoVideoCallout
          youtubeId="Y7c_twmDxQ4"
          thumbnail="/img/TEMP-thumbnail.png"
          heading="How Vault works"
          description="Vault tightly controls access to secrets and encryption keys by authenticating against trusted sources of identity such as Active Directory, LDAP, Kubernetes, CloudFoundry, and cloud platforms."
          person={{
            name: 'Armon Dadgar',
            description: 'Co-founder & CTO',
            avatar: '/img/TEMP-customer-story.png',
          }}
        />
      </div>
    </>
  )
}
