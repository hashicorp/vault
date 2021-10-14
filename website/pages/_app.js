import './style.css'
import '@hashicorp/platform-util/nprogress/style.css'

import Router from 'next/router'
import Head from 'next/head'
import { ErrorBoundary } from '@hashicorp/platform-runtime-error-monitoring'
import createConsentManager from '@hashicorp/react-consent-manager/loader'
import NProgress from '@hashicorp/platform-util/nprogress'
import useAnchorLinkAnalytics from '@hashicorp/platform-util/anchor-link-analytics'
import HashiHead from '@hashicorp/react-head'
import ProductSubnav from 'components/subnav'
import HashiStackMenu from '@hashicorp/react-hashi-stack-menu'
import Footer from 'components/footer'
import Error from './_error'
import AlertBanner from '@hashicorp/react-alert-banner'
import alertBannerData, { ALERT_BANNER_ACTIVE } from '../data/alert-banner'

NProgress({ Router })
const { ConsentManager, openConsentManager } = createConsentManager({
  preset: 'oss',
})

export default function App({ Component, pageProps }) {
  useAnchorLinkAnalytics()

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
            href:
              'https://www.datocms-assets.com/2885/1597163356-vault-favicon.png?h=16&w=16',
            type: 'image/png',
            sizes: '16x16',
          },
          {
            href:
              'https://www.datocms-assets.com/2885/1597163356-vault-favicon.png?h=32&w=32',
            type: 'image/png',
            sizes: '32x32',
          },
          {
            href:
              'https://www.datocms-assets.com/2885/1597163356-vault-favicon.png?h=96&w=96',
            type: 'image/png',
            sizes: '96x96',
          },
          {
            href:
              'https://www.datocms-assets.com/2885/1597163356-vault-favicon.png?h=192&w=192',
            type: 'image/png',
            sizes: '192x192',
          },
        ]}
      />
      {ALERT_BANNER_ACTIVE && (
        <AlertBanner {...alertBannerData} product="vault" hideOnMobile />
      )}
      <HashiStackMenu />
      <ProductSubnav />
      <Component {...pageProps} />
      <Footer openConsentManager={openConsentManager} />
      <ConsentManager className="g-consent-manager" />
    </ErrorBoundary>
  )
}
