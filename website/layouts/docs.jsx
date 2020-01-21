import DocsPage, { getInitialProps } from '../components/docs-page'
import orderData from '../data/docs-navigation.js'
import { frontMatter } from '../pages/docs/**/*.mdx'

function DocsLayoutWrapper(pageMeta) {
  function DocsLayout(props) {
    return (
      <DocsPage
        {...props}
        orderData={orderData}
        frontMatter={frontMatter}
        category="docs"
        pageMeta={pageMeta}
      />
    )
  }

  DocsLayout.getInitialProps = getInitialProps

  return DocsLayout
}

export default DocsLayoutWrapper
