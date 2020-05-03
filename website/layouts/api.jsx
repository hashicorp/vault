import DocsPage, { getInitialProps } from '../components/docs-page'
import orderData from '../data/api-navigation.js'
import { frontMatter } from '../pages/api-docs/**/*.mdx'

function ApiLayoutWrapper(pageMeta) {
  function ApiLayout(props) {
    return (
      <DocsPage
        {...props}
        orderData={orderData}
        frontMatter={frontMatter}
        category="api-docs"
        pageMeta={pageMeta}
      />
    )
  }

  ApiLayout.getInitialProps = getInitialProps

  return ApiLayout
}

export default ApiLayoutWrapper
