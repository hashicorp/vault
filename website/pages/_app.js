import './style.css'
import '../lib/globalThis'
import App from 'next/app'
import NProgress from 'nprogress'
import Router from 'next/router'
import Head from 'next/head'
import HashiHead from '@hashicorp/react-head'
import ProductSubnav from '../components/subnav'
import MegaNav from '@hashicorp/react-mega-nav'
import Footer from '@hashicorp/react-footer'
import { ConsentManager, open } from '@hashicorp/react-consent-manager'
import consentManagerConfig from '../lib/consent-manager-config'
import Bugsnag from '../lib/bugsnag'
import Error from './_error'

Router.events.on('routeChangeStart', NProgress.start)
Router.events.on('routeChangeError', NProgress.done)
Router.events.on('routeChangeComplete', (url) => {
  setTimeout(() => window.analytics.page(url), 0)
  NProgress.done()
})

// Bugsnag
const ErrorBoundary = Bugsnag.getPlugin('react')

class NextApp extends App {
  static async getInitialProps({ Component, ctx }) {
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

  render() {
    const { Component, pageProps } = this.props

    return (
      <ErrorBoundary FallbackComponent={Error}>
        <HashiHead
          is={Head}
          title="Vault by HashiCorp"
          siteName="Vault by HashiCorp"
          description="Vault secures, stores, and tightly controls access to tokens, passwords, certificates, API keys, and other secrets in modern computing. Vault handles leasing, key revocation, key rolling, auditing, and provides secrets as a service through a unified API."
          image="https://www.vaultproject.io/img/og-image.png"
          stylesheet={[
            { href: '/css/nprogress.css' },
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
        <Footer openConsentManager={open} />
        <ConsentManager {...consentManagerConfig} />
      </ErrorBoundary>
    )
  }
}

export default NextApp
