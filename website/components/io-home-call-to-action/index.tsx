import ReactCallToAction from '@hashicorp/react-call-to-action'
import s from './style.module.css'

interface IoHomeCallToActionProps {
  brand: string
  heading: string
  content: string
  links: Array<{
    text: string
    url: string
    type?: 'inbound'
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
        links={links}
      />
    </div>
  )
}
