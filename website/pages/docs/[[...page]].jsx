import { productName, productSlug } from 'data/metadata'
import DocsPage from '@hashicorp/react-docs-page'
import Columns from 'components/columns'
import Tag from 'components/inline-tag'
// Imports below are used in server-side only
/**
 * DEBT: short term patch for "hidden" docs-sidenav items.
 * See components/_temp-enable-hidden-pages for details.
 * Revert to importing from @hashicorp/react-docs-page/server
 * once https://app.asana.com/0/1100423001970639/1200197752405255/f
 * is complete.
 **/
import {
  generateStaticPaths,
  generateStaticProps,
} from 'components/_temp-enable-hidden-pages'

const NAV_DATA_FILE_HIDDEN = 'data/docs-nav-data-hidden.json'
const NAV_DATA_FILE = 'data/docs-nav-data.json'
const CONTENT_DIR = 'content/docs'
const basePath = 'docs'
const additionalComponents = { Columns, Tag }

export default function DocsLayout(props) {
  return (
    <DocsPage
      product={{ name: productName, slug: productSlug }}
      baseRoute={basePath}
      staticProps={props}
      additionalComponents={additionalComponents}
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
      additionalComponents,
    }),
  }
}
