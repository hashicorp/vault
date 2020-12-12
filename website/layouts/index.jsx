import DocsPage from '@hashicorp/react-docs-page'
import Head from 'next/head'
import Link from 'next/link'

function DefaultLayoutWrapper(pageMeta) {
  function DefaultLayout(props) {
    return (
      <DocsPage
        {...props}
        product="vault"
        head={{
          is: Head,
          title: `${pageMeta.page_title} | Vault by HashiCorp`,
          description: pageMeta.description,
          siteName: 'Vault by HashiCorp',
        }}
        sidenav={{
          Link,
          category: 'docs',
          currentPage: props.path,
          data: [],
          order: [],
          disableFilter: true,
        }}
        resourceURL={`https://github.com/hashicorp/vault/blob/master/website/pages/${pageMeta.__resourcePath}`}
      />
    )
  }

  DefaultLayout.getInitialProps = ({ asPath }) => ({ path: asPath })

  return DefaultLayout
}

export default DefaultLayoutWrapper
