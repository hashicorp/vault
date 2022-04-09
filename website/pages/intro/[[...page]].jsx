import { productName, productSlug } from 'data/metadata'
import DocsPage from '@hashicorp/react-docs-page'
import { getStaticGenerationFunctions } from '@hashicorp/react-docs-page/server'

const NAV_DATA_FILE = 'data/intro-nav-data.json'
const CONTENT_DIR = 'content/intro'
const basePath = 'intro'

export default function DocsLayout(props) {
  return (
    <DocsPage
      product={{ name: productName, slug: productSlug }}
      baseRoute={basePath}
      staticProps={props}
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
        localContentDir: CONTENT_DIR,
        navDataFile: NAV_DATA_FILE,
        product: productSlug,
      }
)

export { getStaticPaths, getStaticProps }
