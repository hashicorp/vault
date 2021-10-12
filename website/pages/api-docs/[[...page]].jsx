import { productName, productSlug } from 'data/metadata'
import DocsPage from '@hashicorp/react-docs-page'
// Imports below are used in server-side only
import {
  generateStaticPaths,
  generateStaticProps,
} from '@hashicorp/react-docs-page/server'

const NAV_DATA_FILE_HIDDEN = 'data/api-docs-nav-data-hidden.json'
const NAV_DATA_FILE = 'data/api-docs-nav-data.json'
const CONTENT_DIR = 'content/api-docs'
const basePath = 'api-docs'

export default function DocsLayout(props) {
  return (
    <DocsPage
      product={{ name: productName, slug: productSlug }}
      baseRoute={basePath}
      staticProps={props}
    />
  )
}

export async function getStaticPaths() {
  return {
    fallback: false,
    paths: await generateStaticPaths({
      navDataFile: NAV_DATA_FILE,
      navDataFileHidden: NAV_DATA_FILE_HIDDEN,
      localContentDir: CONTENT_DIR,
    }),
  }
}

export async function getStaticProps({ params }) {
  return {
    props: await generateStaticProps({
      navDataFile: NAV_DATA_FILE,
      navDataFileHidden: NAV_DATA_FILE_HIDDEN,
      localContentDir: CONTENT_DIR,
      product: { name: productName, slug: productSlug },
      params,
    }),
  }
}
