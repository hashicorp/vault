import { productName, productSlug } from 'data/metadata'
import DocsPage from '@hashicorp/react-docs-page'
import Columns from 'components/columns'
import Tag from 'components/inline-tag'
// Imports below are used in server-side only
import { getStaticGenerationFunctions } from '@hashicorp/react-docs-page/server'

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

const { getStaticPaths, getStaticProps } = getStaticGenerationFunctions(
  process.env.ENABLE_VERSIONED_DOCS === 'true'
    ? {
        strategy: 'remote',
        basePath: basePath,
        fallback: 'blocking',
        revalidate: 360, // 1 hour
        product: productSlug,
      }
    : {
        strategy: 'fs',
        basePath: basePath,
        localContentDir: CONTENT_DIR,
        navDataFile: NAV_DATA_FILE,
        product: productSlug,
        revalidate: false,
      }
)

export { getStaticPaths, getStaticProps }
