import Head from 'next/head'
import Link from 'next/link'
import Button from '@hashicorp/react-button'
import ProductDownloader from '@hashicorp/react-product-downloader'
import HashiHead from '@hashicorp/react-head'
import { VERSION, CHANGELOG_URL, packageManagers } from 'data/version'
import { productName, productSlug } from 'data/metadata'
import s from './style.module.css'

function MerchandisingSlot() {
  return (
    <div className={s.merchandisingSlot}>
      <div className={s.centerWrapper}>
        <p>
          Want all of the power and security of Vault, without the complexity
          and overhead of managing it yourself?
        </p>
        <Button
          title="Sign up for HCP Vault"
          linkType="inbound"
          url="https://portal.cloud.hashicorp.com/sign-up?utm_source=vault_io&utm_content=download_cta"
          theme={{
            variant: 'tertiary',
            brand: 'vault',
          }}
        />
      </div>
    </div>
  )
}

export default function DownloadsPage({ releases }) {
  const changelogUrl = CHANGELOG_URL.length
    ? CHANGELOG_URL
    : `https://github.com/hashicorp/vault/blob/v${VERSION}/CHANGELOG.md`
  return (
    <div id="p-downloads" className="g-container">
      <HashiHead is={Head} title="Downloads | Vault by Hashicorp" />
      <ProductDownloader
        releases={releases}
        packageManagers={packageManagers}
        productName={productName}
        productId={productSlug}
        latestVersion={VERSION}
        changelog={changelogUrl}
        getStartedDescription="Follow step-by-step tutorials on the essentials of Nomad."
        getStartedLinks={[
          {
            label: 'Getting Started with the CLI',
            href:
              'http://learn.hashicorp.com/collections/vault/getting-started',
          },
          {
            label: 'Getting Started with Vault UI',
            href:
              'http://learn.hashicorp.com/collections/vault/getting-started-ui',
          },
          {
            label: 'Vault on HCP',
            href:
              'http://learn.hashicorp.com/collections/vault/getting-started-ui',
          },
          {
            label: 'View all Vault tutorials',
            href: 'https://learn.hashicorp.com/vault',
          },
        ]}
        logo={
          <img
            className={s.logo}
            alt="Nomad"
            src={require('./img/vault-logo.svg')}
          />
        }
        tutorialLink={{
          href: 'https://learn.hashicorp.com/vault',
          label: 'View Tutorials at HashiCorp Learn',
        }}
        merchandisingSlot={
          <>
            <MerchandisingSlot />
            <p className={s.releaseNote}>
              Release notes are available in our{' '}
              <Link href={`/docs/release-notes/${VERSION}`}>
                <a>documentation</a>
              </Link>
              .
            </p>
          </>
        }
      />
    </div>
  )
}

export async function getStaticProps() {
  return fetch(`https://releases.hashicorp.com/vault/index.json`, {
    headers: {
      'Cache-Control': 'no-cache',
    },
  })
    .then((res) => res.json())
    .then((result) => {
      return {
        props: {
          releases: result,
        },
      }
    })
    .catch(() => {
      throw new Error(
        `--------------------------------------------------------
        Unable to resolve version ${VERSION} on releases.hashicorp.com from link
        <https://releases.hashicorp.com/${productSlug}/${VERSION}/index.json>. Usually this
        means that the specified version has not yet been released. The downloads page
        version can only be updated after the new version has been released, to ensure
        that it works for all users.
        ----------------------------------------------------------`
      )
    })
}
