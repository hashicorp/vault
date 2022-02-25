import * as React from 'react'
import { Products } from '@hashicorp/platform-product-meta'
import classNames from 'classnames'
import Image from 'next/image'
import Button from '@hashicorp/react-button'
import s from './style.module.css'

interface IoUsecaseSectionProps {
  brand?: Products | 'neutral'
  bottomIsFlush?: boolean
  eyebrow: string
  heading: string
  description: string
  media?: {
    src: string
    width: string
    height: string
    alt: string
  }
  cta?: {
    text: string
    link: string
  }
}

export default function IoUsecaseSection({
  brand = 'neutral',
  bottomIsFlush = false,
  eyebrow,
  heading,
  description,
  media,
  cta,
}: IoUsecaseSectionProps): React.ReactElement {
  return (
    <section
      className={classNames(s.section, s[brand], bottomIsFlush && s.isFlush)}
    >
      <div className={s.container}>
        <p className={s.eyebrow}>{eyebrow}</p>
        <div className={s.columns}>
          <div className={s.column}>
            <h2 className={s.heading}>{heading}</h2>
            {media?.src ? (
              <div
                className={s.description}
                dangerouslySetInnerHTML={{
                  __html: description,
                }}
              />
            ) : null}
            {cta?.link && cta?.text ? (
              <div className={s.cta}>
                <Button
                  title={cta.text}
                  url={cta.link}
                  theme={{
                    brand: brand,
                  }}
                />
              </div>
            ) : null}
          </div>
          <div className={s.column}>
            {media?.src ? (
              // eslint-disable-next-line jsx-a11y/alt-text
              <Image {...media} />
            ) : (
              <div
                className={s.description}
                dangerouslySetInnerHTML={{
                  __html: description,
                }}
              />
            )}
          </div>
        </div>
      </div>
    </section>
  )
}
