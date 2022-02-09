import './style.css'
import '@hashicorp/platform-util/nprogress/style.css'

import Router from 'next/router'
import Head from 'next/head'
import rivetQuery from '@hashicorp/nextjs-scripts/dato/client'
import { ErrorBoundary } from '@hashicorp/platform-runtime-error-monitoring'
import createConsentManager from '@hashicorp/react-consent-manager/loader'
import localConsentManagerServices from 'lib/consent-manager-services'
import NProgress from '@hashicorp/platform-util/nprogress'
import useFathomAnalytics from '@hashicorp/platform-analytics'
import useAnchorLinkAnalytics from '@hashicorp/platform-util/anchor-link-analytics'
import HashiHead from '@hashicorp/react-head'
import Error from './_error'
import AlertBanner from '@hashicorp/react-alert-banner'
import alertBannerData, { ALERT_BANNER_ACTIVE } from '../data/alert-banner'
import StandardLayout from 'layouts/standard'

NProgress({ Router })
const { ConsentManager } = createConsentManager({
  preset: 'oss',
  otherServices: [...localConsentManagerServices],
})

export default function App({ Component, pageProps, layoutData }) {
  useFathomAnalytics()
  useAnchorLinkAnalytics()

  const Layout = Component.layout ?? StandardLayout

  return (
    <ErrorBoundary FallbackComponent={Error}>
      <HashiHead
        is={Head}
        title="Vault by HashiCorp"
        siteName="Vault by HashiCorp"
        description="Vault secures, stores, and tightly controls access to tokens, passwords, certificates, API keys, and other secrets in modern computing. Vault handles leasing, key revocation, key rolling, auditing, and provides secrets as a service through a unified API."
        image="https://www.vaultproject.io/img/og-image.png"
        icon={[
          {
            href: 'https://www.datocms-assets.com/2885/1597163356-vault-favicon.png?h=16&w=16',
            type: 'image/png',
            sizes: '16x16',
          },
          {
            href: 'https://www.datocms-assets.com/2885/1597163356-vault-favicon.png?h=32&w=32',
            type: 'image/png',
            sizes: '32x32',
          },
          {
            href: 'https://www.datocms-assets.com/2885/1597163356-vault-favicon.png?h=96&w=96',
            type: 'image/png',
            sizes: '96x96',
          },
          {
            href: 'https://www.datocms-assets.com/2885/1597163356-vault-favicon.png?h=192&w=192',
            type: 'image/png',
            sizes: '192x192',
          },
        ]}
      />
      {ALERT_BANNER_ACTIVE && (
        <AlertBanner {...alertBannerData} product="vault" hideOnMobile />
      )}
      <Layout {...(layoutData && { data: layoutData })}>
        <Component {...pageProps} />
      </Layout>
      <ConsentManager className="g-consent-manager" />
    </ErrorBoundary>
  )
}

App.getInitialProps = async ({ Component, ctx }) => {
  const layoutQuery = Component.layout
    ? Component.layout?.rivetParams ?? null
    : StandardLayout.rivetParams

  const layoutData = layoutQuery ? await rivetQuery(layoutQuery) : null

  let pageProps = {}

  if (Component.getInitialProps) {
    pageProps = await Component.getInitialProps(ctx)
  }
  return { pageProps, layoutData }
}
