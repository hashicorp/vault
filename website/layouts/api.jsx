import DocsPage from '@hashicorp/react-docs-page'
import order from '../data/api-navigation.js'
import { frontMatter as data } from '../pages/api-docs/**/*.mdx'
import Head from 'next/head'
import Link from 'next/link'

function ApiLayoutWrapper(pageMeta) {
  function ApiLayout(props) {
    return (
      <DocsPage
        {...props}
        product="vault"
        head={{
          is: Head,
          title: `${pageMeta.page_title} | Vault by HashiCorp`,
          description: pageMeta.description,
          siteName: 'Vault by HashiCorp'
        }}
        sidenav={{
          Link,
          category: 'api-docs',
          currentPage: props.path,
          data,
          order
        }}
        resourceURL={`https://github.com/hashicorp/vault/blob/master/website/pages/${pageMeta.__resourcePath}`}
      />
    )
  }

  ApiLayout.getInitialProps = ({ asPath }) => ({ path: asPath })

  return ApiLayout
}

export default ApiLayoutWrapper
