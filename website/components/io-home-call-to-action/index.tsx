import ReactCallToAction from '@hashicorp/react-call-to-action'
import { Products } from '@hashicorp/platform-product-meta'
import s from './style.module.css'

interface IoHomeCallToActionProps {
  brand: Products
  heading: string
  content: string
  links: Array<{
    text: string
    url: string
  }>
}

export default function IoHomeCallToAction({
  brand,
  heading,
  content,
  links,
}: IoHomeCallToActionProps) {
  return (
    <div className={s.callToAction}>
      <ReactCallToAction
        variant="compact"
        heading={heading}
        content={content}
        product={brand}
        theme="dark"
        links={links.map(({ text, url }, index) => {
          return {
            text,
            url,
            type: index === 1 ? 'inbound' : null,
          }
        })}
      />
    </div>
  )
}
