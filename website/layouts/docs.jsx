import DocsPage from '@hashicorp/react-docs-page'
import order from '../data/docs-navigation.js'
import { frontMatter } from '../pages/docs/**/*.mdx'
import Head from 'next/head'
import Link from 'next/link'

function DocsLayoutWrapper(pageMeta) {
  function DocsLayout(props) {
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
          category: 'docs',
          currentPage: props.path,
          data: frontMatter,
          order
        }}
        resourceURL={`https://github.com/hashicorp/vault/blob/master/website/pages/${pageMeta.__resourcePath}`}
      />
    )
  }

  DocsLayout.getInitialProps = ({ asPath }) => ({ path: asPath })

  return DocsLayout
}

export default DocsLayoutWrapper
