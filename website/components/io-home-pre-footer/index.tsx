import * as React from 'react'
import classNames from 'classnames'
import { Products } from '@hashicorp/platform-product-meta'
import { IconArrowRight16 } from '@hashicorp/flight-icons/svg-react/arrow-right-16'
import s from './style.module.css'

interface IoHomePreFooterProps {
  brand: Products
  heading: string
  description: string
  ctas: [IoHomePreFooterCard, IoHomePreFooterCard, IoHomePreFooterCard]
}

export default function IoHomePreFooter({
  brand,
  heading,
  description,
  ctas,
}: IoHomePreFooterProps) {
  return (
    <div className={classNames(s.preFooter, s[brand])}>
      <div className={s.container}>
        <div className={s.content}>
          <h2 className={s.heading}>{heading}</h2>
          <p className={s.description}>{description}</p>
        </div>
        <div className={s.cards}>
          {ctas.map((cta, index) => {
            return (
              <IoHomePreFooterCard
                key={index}
                brand={brand}
                link={cta.link}
                heading={cta.heading}
                description={cta.description}
                cta={cta.cta}
              />
            )
          })}
        </div>
      </div>
    </div>
  )
}

interface IoHomePreFooterCard {
  brand?: string
  link: string
  heading: string
  description: string
  cta: string
}

function IoHomePreFooterCard({
  brand,
  link,
  heading,
  description,
  cta,
}: IoHomePreFooterCard): React.ReactElement {
  return (
    <a
      href={link}
      className={s.card}
      style={
        {
          '--primary': `var(--${brand})`,
          '--secondary': `var(--${brand}-secondary)`,
        } as React.CSSProperties
      }
    >
      <h3 className={s.cardHeading}>{heading}</h3>
      <p className={s.cardDescription}>{description}</p>
      <span className={s.cardCta}>
        {cta} <IconArrowRight16 />
      </span>
    </a>
  )
}
