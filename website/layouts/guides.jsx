import DocsPage from '@hashicorp/react-docs-page'
import order from '../data/guides-navigation.js'
import { frontMatter as data } from '../pages/guides/**/*.mdx'
import { MDXProvider } from '@mdx-js/react'
import Head from 'next/head'
import Link from 'next/link'
import Tabs, { Tab } from '../components/tabs'
import EnterpriseAlert from '../components/enterprise-alert'

const DEFAULT_COMPONENTS = { Tabs, Tab, EnterpriseAlert }

export default function GuidesLayoutWrapper(pageMeta) {
  function GuidesLayout(props) {
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
            category: 'guides',
            currentPage: props.path,
            data,
            order,
          }}
          resourceURL={`https://github.com/hashicorp/vault/blob/master/website/pages/${pageMeta.__resourcePath}`}
        />
      </MDXProvider>
    )
  }

  GuidesLayout.getInitialProps = ({ asPath }) => ({ path: asPath })

  return GuidesLayout
}
