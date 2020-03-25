import DocsSidenav from '@hashicorp/react-docs-sidenav'
import Content from '@hashicorp/react-content'
import InlineSvg from '@hashicorp/react-inline-svg'
import githubIcon from './img/github-icon.svg?include'
import Link from 'next/link'
import Head from 'next/head'
import HashiHead from '@hashicorp/react-head'

export default function DocsPage({
  children,
  path,
  orderData,
  frontMatter,
  category,
  pageMeta
}) {
  return (
    <div id="p-docs">
      <HashiHead
        is={Head}
        title={`${pageMeta.page_title} | Vault by Hashicorp`}
        description={pageMeta.description}
      >
        {pageMeta.deprecated && <meta name="robots" content="noindex" />}
      </HashiHead>
      <div className="content-wrap g-container">
        <div id="sidebar" role="complementary">
          <div className="nav docs-nav">
            <DocsSidenav
              currentPage={path}
              category={category}
              order={orderData}
              data={frontMatter}
              Link={Link}
            />
          </div>
        </div>

        <div id="inner" role="main">
          <Content product="vault" content={children} />
        </div>
      </div>
      <div id="edit-this-page" className="g-container">
        <a
          href={`https://github.com/hashicorp/vault/blob/master/website/pages/${pageMeta.__resourcePath}`}
        >
          <InlineSvg src={githubIcon} />
          <span>Edit this page</span>
        </a>
      </div>
    </div>
  )
}

export async function getInitialProps({ asPath }) {
  return { path: asPath }
}
