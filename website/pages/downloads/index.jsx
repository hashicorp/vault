import fetch from 'isomorphic-unfetch'
import { VERSION } from '../../data/version.js'
import ProductDownloader from '@hashicorp/react-product-downloader'
import Head from 'next/head'

export default function DownloadsPage({ downloadData }) {
  return (
    <div id="p-downloads" className="g-container">
      <Head>
        <title key="title">Downloads | Vault by HashiCorp</title>
      </Head>
      <ProductDownloader
        product="Vault"
        version={VERSION}
        downloads={downloadData}
        prerelease={{
          type: 'beta release',
          name: '1.4.0',
          version: '1.4.0-beta1'
        }}
      />
    </div>
  )
}

export async function unstable_getStaticProps() {
  return fetch(`https://releases.hashicorp.com/vault/${VERSION}/index.json`)
    .then(r => r.json())
    .then(r => {
      // TODO: restructure product-downloader to run this logic internally
      return r.builds.reduce((acc, build) => {
        if (!acc[build.os]) acc[build.os] = {}
        acc[build.os][build.arch] = build.url
        return acc
      }, {})
    })
    .then(r => ({ props: { downloadData: r } }))
    .catch(() => {
      throw new Error(
        `--------------------------------------------------------
        Unable to resolve version ${VERSION} on releases.hashicorp.com from link
        <https://releases.hashicorp.com/vault/${VERSION}/index.json>. Usually this
        means that the specified version has not yet been released. The downloads page
        version can only be updated after the new version has been released, to ensure
        that it works for all users.
        ----------------------------------------------------------`
      )
    })
}
