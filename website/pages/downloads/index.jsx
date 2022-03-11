import ProductDownloadsPage from '@hashicorp/react-product-downloads-page'
import { generateStaticProps } from '@hashicorp/react-product-downloads-page/server'
import { VERSION } from 'data/version'
import { productSlug } from 'data/metadata'
import baseProps from 'components/downloads-props'

export default function DownloadsPage(staticProps) {
  return <ProductDownloadsPage {...baseProps()} {...staticProps} />
}

export function getStaticProps() {
  return generateStaticProps({
    product: productSlug,
    latestVersion: VERSION,
  })
}
