import * as React from 'react'
import Button from '@hashicorp/react-button'
import { IoCardProps } from 'components/io-card'
import s from './style.module.css'

interface IoCardContaianerProps {
  heading?: string
  description?: string
  label?: string
  cta?: {
    url: string
    text: string
  }
  cardsPerRow: 3 | 4
  children:
    | Array<React.ReactElement<IoCardProps>>
    | React.ReactElement<IoCardProps>
}

export default function IoCardContaianer({
  heading,
  description,
  label,
  cta,
  cardsPerRow = 3,
  children,
}: IoCardContaianerProps): React.ReactElement {
  return (
    <div className={s.cardContainer}>
      {heading || description ? (
        <header className={s.header}>
          {heading ? <h2 className={s.heading}>{heading}</h2> : null}
          {description ? <p className={s.description}>{description}</p> : null}
        </header>
      ) : null}
      {label || cta ? (
        <header className={s.subHeader}>
          {label ? <h3 className={s.label}>{label}</h3> : null}
          {cta ? (
            <Button
              title={cta.text}
              href={cta.url}
              linkType="inbound"
              theme={{
                brand: 'neutral',
                variant: 'tertiary',
              }}
            />
          ) : null}
        </header>
      ) : null}
      <ul
        className={s.cardList}
        style={
          {
            '--per-row': cardsPerRow,
            '--length': React.Children.count(children),
          } as React.CSSProperties
        }
      >
        {React.Children.map(children, (child, index) => {
          // Index is stable
          // eslint-disable-next-line react/no-array-index-key
          return <li key={index}>{React.cloneElement(child)}</li>
        })}
      </ul>
    </div>
  )
}
