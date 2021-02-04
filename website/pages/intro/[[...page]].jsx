import { productName, productSlug } from 'data/metadata'
import order from 'data/intro-navigation.js'
import DocsPage from '@hashicorp/react-docs-page'
import {
  generateStaticPaths,
  generateStaticProps,
} from '@hashicorp/react-docs-page/server'

const subpath = 'intro'

export default function GuidesLayout(props) {
  return (
    <DocsPage
      product={{ name: productName, slug: productSlug }}
      subpath={subpath}
      order={order}
      mainBranch="master"
      staticProps={props}
    />
  )
}

export async function getStaticPaths() {
  return generateStaticPaths(subpath)
}

export async function getStaticProps({ params }) {
  return generateStaticProps({
    subpath,
    productName,
    params,
  })
}
