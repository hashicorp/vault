import * as React from 'react'
import s from './style.module.css'

interface IoUsecaseHeroProps {
  eyebrow: string
  heading: string
  description: string
}

export default function IoUsecaseHero({
  eyebrow,
  heading,
  description,
}: IoUsecaseHeroProps): React.ReactElement {
  return (
    <header className={s.hero}>
      <div className={s.container}>
        <div className={s.content}>
          <p className={s.eyebrow}>{eyebrow}</p>
          <h1 className={s.heading}>{heading}</h1>
          <p className={s.description}>{description}</p>
        </div>
      </div>
    </header>
  )
}
