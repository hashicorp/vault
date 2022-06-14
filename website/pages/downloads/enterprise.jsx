import { VERSION } from 'data/version'
import { productSlug } from 'data/metadata'
import ProductDownloadsPage from '@hashicorp/react-product-downloads-page'
import { generateStaticProps } from '@hashicorp/react-product-downloads-page/server'
import baseProps from 'components/downloads-props'
import s from './style.module.css'

export default function DownloadsPage(staticProps) {
  return (
    <>
      <ProductDownloadsPage
        enterpriseMode={true}
        {...baseProps(
          <p className={s.legalNotice}>
            <em>
              The following shall apply unless your organization has a
              separately signed Enterprise License Agreement or Evaluation
              Agreement governing your use of the package: Enterprise packages
              in this repository are subject to the license terms located in the
              package. Please read the license terms prior to using the package.
              Your installation and use of the package constitutes your
              acceptance of these terms. If you do not accept the terms, do not
              use the package.
            </em>
          </p>
        )}
        {...staticProps}
      />
    </>
  )
}

export async function getStaticProps() {
  return generateStaticProps({
    product: productSlug,
    latestVersion: VERSION,
  })
}
