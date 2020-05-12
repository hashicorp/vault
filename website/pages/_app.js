import './style.css'
import '@hashicorp/nextjs-scripts/lib/nprogress/style.css'

import '../lib/globalThis'
import Router from 'next/router'
import Head from 'next/head'
import { ErrorBoundary } from '@hashicorp/nextjs-scripts/lib/bugsnag'
import createConsentManager from '@hashicorp/nextjs-scripts/lib/consent-manager'
import NProgress from '@hashicorp/nextjs-scripts/lib/nprogress'
import useAnchorLinkAnalytics from '@hashicorp/nextjs-scripts/lib/anchor-link-analytics'
import HashiHead from '@hashicorp/react-head'
import ProductSubnav from '../components/subnav'
import MegaNav from '@hashicorp/react-mega-nav'
import Footer from '../components/footer'
import Error from './_error'

NProgress({ Router })
const { ConsentManager, openConsentManager } = createConsentManager({
  preset: 'oss',
})

function App({ Component, pageProps }) {
  useAnchorLinkAnalytics()

  return (
    <ErrorBoundary FallbackComponent={Error}>
      <HashiHead
        is={Head}
        title="Vault by HashiCorp"
        siteName="Vault by HashiCorp"
        description="Vault secures, stores, and tightly controls access to tokens, passwords, certificates, API keys, and other secrets in modern computing. Vault handles leasing, key revocation, key rolling, auditing, and provides secrets as a service through a unified API."
        image="https://www.vaultproject.io/img/og-image.png"
        stylesheet={[
          {
            href:
              'https://fonts.googleapis.com/css?family=Open+Sans:300,400,600,700&display=swap',
          },
        ]}
        icon={[
          {
            href:
              'https://www.datocms-assets.com/2885/1527033389-favicon.png?h=16&w=16',
            type: 'image/png',
            sizes: '16x16',
          },
          {
            href:
              'https://www.datocms-assets.com/2885/1527033389-favicon.png?h=32&w=32',
            type: 'image/png',
            sizes: '32x32',
          },
          {
            href:
              'https://www.datocms-assets.com/2885/1527033389-favicon.png?h=96&w=96',
            type: 'image/png',
            sizes: '96x96',
          },
          {
            href:
              'https://www.datocms-assets.com/2885/1527033389-favicon.png?h=192&w=192',
            type: 'image/png',
            sizes: '192x192',
          },
        ]}
      />
      <MegaNav product="Vault" />
      <ProductSubnav />
      <Component {...pageProps} />
      <Footer openConsentManager={openConsentManager} />
      <ConsentManager />
    </ErrorBoundary>
  )
}

App.getInitialProps = async ({ Component, ctx }) => {
  let pageProps = {}

  if (Component.getInitialProps) {
    pageProps = await Component.getInitialProps(ctx)
  } else if (Component.isMDXComponent) {
    // fix for https://github.com/mdx-js/mdx/issues/382
    const mdxLayoutComponent = Component({}).props.originalType
    if (mdxLayoutComponent.getInitialProps) {
      pageProps = await mdxLayoutComponent.getInitialProps(ctx)
    }
  }

  return { pageProps }
}

export default App
