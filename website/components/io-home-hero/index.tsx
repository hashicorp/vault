import * as React from 'react'
import { Products } from '@hashicorp/platform-product-meta'
import Button from '@hashicorp/react-button'
import classNames from 'classnames'
import s from './style.module.css'

interface IoHomeHeroProps {
  pattern: string
  brand: Products | 'neutral'
  heading: string
  description: string
  ctas: Array<{
    title: string
    link: string
  }>
  cards: Array<IoHomeHeroCardProps>
}

export default function IoHomeHero({
  pattern,
  brand,
  heading,
  description,
  ctas,
  cards,
}: IoHomeHeroProps) {
  const [loaded, setLoaded] = React.useState(false)

  React.useEffect(() => {
    setTimeout(() => {
      setLoaded(true)
    }, 250)
  }, [])

  return (
    <header
      className={classNames(s.hero, loaded && s.loaded)}
      style={
        {
          '--pattern': `url(${pattern})`,
        } as React.CSSProperties
      }
    >
      <span className={s.pattern} />
      <div className={s.container}>
        <div className={s.content}>
          <h1 className={s.heading}>{heading}</h1>
          <p className={s.description}>{description}</p>
          {ctas && (
            <div className={s.ctas}>
              {ctas.map((cta, index) => {
                return (
                  <Button
                    key={index}
                    title={cta.title}
                    url={cta.link}
                    linkType="inbound"
                    theme={{
                      brand: 'neutral',
                      variant: 'tertiary',
                      background: 'light',
                    }}
                  />
                )
              })}
            </div>
          )}
        </div>
        {cards && (
          <div className={s.cards}>
            {cards.map((card, index) => {
              return (
                <IoHomeHeroCard
                  key={index}
                  index={index}
                  heading={card.heading}
                  description={card.description}
                  cta={{
                    brand: index === 0 ? 'neutral' : brand,
                    title: card.cta.title,
                    link: card.cta.link,
                  }}
                  subText={card.subText}
                />
              )
            })}
          </div>
        )}
      </div>
    </header>
  )
}

interface IoHomeHeroCardProps {
  index?: number
  heading: string
  description: string
  cta: {
    title: string
    link: string
    brand?: 'neutral' | Products
  }
  subText: string
}

function IoHomeHeroCard({
  index,
  heading,
  description,
  cta,
  subText,
}: IoHomeHeroCardProps): React.ReactElement {
  return (
    <article
      className={s.card}
      style={
        {
          '--index': index,
        } as React.CSSProperties
      }
    >
      <h2 className={s.cardHeading}>{heading}</h2>
      <p className={s.cardDescription}>{description}</p>
      <Button
        title={cta.title}
        url={cta.link}
        theme={{
          variant: 'primary',
          brand: cta.brand,
        }}
      />
      <p className={s.cardSubText}>{subText}</p>
    </article>
  )
}
