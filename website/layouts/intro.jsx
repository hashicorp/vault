import DocsPage, { getInitialProps } from '../components/docs-page'
import orderData from '../data/intro-navigation.js'
import { frontMatter } from '../pages/intro/**/*.mdx'

function IntroLayoutWrapper(pageMeta) {
  function IntroLayout(props) {
    return (
      <DocsPage
        {...props}
        orderData={orderData}
        frontMatter={frontMatter}
        category="intro"
        pageMeta={pageMeta}
      />
    )
  }

  IntroLayout.getInitialProps = getInitialProps

  return IntroLayout
}

export default IntroLayoutWrapper
