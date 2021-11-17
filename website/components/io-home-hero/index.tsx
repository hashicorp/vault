import * as React from 'react'
import Button from '@hashicorp/react-button'
import classNames from 'classnames'
import s from './style.module.css'

interface IoHomeHeroProps {
  brand: 'vault' | 'consul'
  heading: string
  description: string
  ctas: Array<{
    title: string
    url: string
  }>
  cards: Array<IoHomeHeroCardProps>
}

export default function IoHomeHero({
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
    <header className={classNames(s.hero, loaded && s.loaded)}>
      <span className={s.pattern} />
      <div className={s.container}>
        <div className={s.content}>
          <h1 className={s.heading}>{heading}</h1>
          <p className={s.description}>{description}</p>
          {ctas && (
            <div className={s.ctas}>
              {ctas.map((cta) => {
                return (
                  <Button
                    title={cta.title}
                    url={cta.url}
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
                  index={index}
                  heading={card.heading}
                  description={card.description}
                  cta={{
                    brand: index === 0 ? 'neutral' : brand,
                    title: card.cta.title,
                    url: card.cta.url,
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
    url: string
    brand?: 'neutral' | 'vault' | 'consul'
  }
  subText: string
}

function IoHomeHeroCard({
  index,
  heading,
  description,
  cta,
  subText,
}: IoHomeHeroCardProps) {
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
        url={cta.url}
        theme={{
          variant: 'primary',
          brand: cta.brand,
        }}
      />
      <p className={s.cardSubText}>{subText}</p>
    </article>
  )
}
