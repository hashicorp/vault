import * as React from 'react'
import classNames from 'classnames'
import Button from '@hashicorp/react-button'
import IoCard, { IoCardProps } from 'components/io-card'
import s from './style.module.css'

interface IoCardContaianerProps {
  theme?: 'light' | 'dark'
  heading?: string
  description?: string
  label?: string
  cta?: {
    url: string
    text: string
  }
  cardsPerRow: 3 | 4
  cards: Array<IoCardProps>
}

export default function IoCardContaianer({
  theme = 'light',
  heading,
  description,
  label,
  cta,
  cardsPerRow = 3,
  cards,
}: IoCardContaianerProps): React.ReactElement {
  return (
    <div className={classNames(s.cardContainer, s[theme])}>
      {heading || description ? (
        <header className={s.header}>
          {heading ? <h2 className={s.heading}>{heading}</h2> : null}
          {description ? <p className={s.description}>{description}</p> : null}
        </header>
      ) : null}
      {cards.length ? (
        <>
          {label || cta ? (
            <header className={s.subHeader}>
              {label ? <h3 className={s.label}>{label}</h3> : null}
              {cta ? (
                <Button
                  title={cta.text}
                  url={cta.url}
                  linkType="inbound"
                  theme={{
                    brand: 'neutral',
                    variant: 'tertiary',
                    background: theme,
                  }}
                />
              ) : null}
            </header>
          ) : null}
          <ul
            className={classNames(
              s.cardList,
              cardsPerRow === 3 && s.threeUp,
              cardsPerRow === 4 && s.fourUp
            )}
            style={
              {
                '--length': cards.length,
              } as React.CSSProperties
            }
          >
            {cards.map((card, index) => {
              return (
                // Index is stable
                // eslint-disable-next-line react/no-array-index-key
                <li key={index}>
                  <IoCard variant={theme} {...card} />
                </li>
              )
            })}
          </ul>
        </>
      ) : null}
    </div>
  )
}
