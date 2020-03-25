import DocsPage, { getInitialProps } from '../components/docs-page'
import orderData from '../data/guides-navigation.js'
import { frontMatter } from '../pages/guides/**/*.mdx'

function GuidesLayoutWrapper(pageMeta) {
  function GuidesLayout(props) {
    return (
      <DocsPage
        {...props}
        orderData={orderData}
        frontMatter={frontMatter}
        category="guides"
        pageMeta={pageMeta}
      />
    )
  }

  GuidesLayout.getInitialProps = getInitialProps

  return GuidesLayout
}

export default GuidesLayoutWrapper
