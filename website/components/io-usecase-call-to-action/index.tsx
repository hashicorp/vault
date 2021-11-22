import * as React from 'react'
import ReactCallToAction from '@hashicorp/react-call-to-action'
import s from './style.module.css'

interface IoUsecaseCallToActionProps {
  brand: string
  theme?: 'light' | 'dark'
  heading: string
  content: string
  links: Array<{
    text: string
    url: string
  }>
}

export default function IoUsecaseCallToAction({
  brand,
  theme,
  heading,
  content,
  links,
}: IoUsecaseCallToActionProps): React.ReactElement {
  return (
    <div
      className={s.callToAction}
      style={
        {
          '--background-color': `var(--${brand})`,
        } as React.CSSProperties
      }
    >
      <ReactCallToAction
        variant="compact"
        heading={heading}
        content={content}
        product="neutral"
        theme={theme}
        links={links}
      />
    </div>
  )
}
