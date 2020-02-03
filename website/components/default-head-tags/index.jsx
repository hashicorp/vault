import Head from 'next/head'

export default function DefaultHeadTags() {
  return (
    <Head>
      <title key="title">Vault by HashiCorp</title>
      <meta charSet="utf-8" />
      <meta httpEquiv="x-ua-compatible" content="ie=edge" />
      {/* ref: https://www.phpied.com/minimum-viable-sharing-meta-tags/ */}
      <meta property="og:locale" content="en_US" />
      <meta property="og:type" content="website" />
      <meta
        property="og:site_name"
        content="Vault by HashiCorp"
        key="og-name"
      />
      <meta name="twitter:site" content="@HashiCorp" />
      <meta name="twitter:card" content="summary_large_image" />
      <meta
        property="article:publisher"
        content="https://www.facebook.com/HashiCorp/"
      />
      <meta
        name="description"
        property="og:description"
        content="Vault secures, stores, and tightly controls access to tokens, passwords, certificates, API keys, and other secrets in modern computing. Vault handles leasing, key revocation, key rolling, auditing, and provides secrets as a service through a unified API."
        key="description"
      />
      <meta
        property="og:image"
        content="https://www.vaultproject.io/img/og-image.png"
        key="image"
      />
      <link
        sizes="16x16"
        type="image/png"
        rel="icon"
        href="https://www.datocms-assets.com/2885/1527033389-favicon.png?h=16&w=16"
      />
      <link
        sizes="32x32"
        type="image/png"
        rel="icon"
        href="https://www.datocms-assets.com/2885/1527033389-favicon.png?h=32&w=32"
      />
      <link
        sizes="96x96"
        type="image/png"
        rel="icon"
        href="https://www.datocms-assets.com/2885/1527033389-favicon.png?h=96&w=96"
      />
      <link
        sizes="192x192"
        type="image/png"
        rel="icon"
        href="https://www.datocms-assets.com/2885/1527033389-favicon.png?h=192&w=192"
      />
      <link rel="stylesheet" href="/css/nprogress.css"></link>
      <link
        href="https://fonts.googleapis.com/css?family=Open+Sans:300,400,600,700&display=swap"
        rel="stylesheet"
      />
    </Head>
  )
}
