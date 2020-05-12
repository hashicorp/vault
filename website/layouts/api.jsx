import DocsPage from '@hashicorp/react-docs-page'
import order from '../data/api-navigation.js'
import { frontMatter as data } from '../pages/api-docs/**/*.mdx'
import { MDXProvider } from '@mdx-js/react'
import Head from 'next/head'
import Link from 'next/link'
import Tabs, { Tab } from '../components/tabs'
import EnterpriseAlert from '@hashicorp/react-enterprise-alert'

const DEFAULT_COMPONENTS = { Tabs, Tab, EnterpriseAlert }

export default function ApiLayoutWrapper(pageMeta) {
  function ApiLayout(props) {
    return (
      <MDXProvider components={DEFAULT_COMPONENTS}>
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
            category: 'api-docs',
            currentPage: props.path,
            data,
            order,
          }}
          resourceURL={`https://github.com/hashicorp/vault/blob/master/website/pages/${pageMeta.__resourcePath}`}
        />
      </MDXProvider>
    )
  }

  ApiLayout.getInitialProps = ({ asPath }) => ({ path: asPath })

  return ApiLayout
}
