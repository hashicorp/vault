import { productName, productSlug } from 'data/metadata'
import DocsPage from '@hashicorp/react-docs-page'
// Imports below are used in server-side only
import { getStaticGenerationFunctions } from '@hashicorp/react-docs-page/server'

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
      showVersionSelect={process.env.ENABLE_VERSIONED_DOCS === 'true'}
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
