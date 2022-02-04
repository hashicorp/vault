import * as React from 'react'
import Image from 'next/image'
import Button from '@hashicorp/react-button'
import { Products } from '@hashicorp/platform-product-meta'
import { IoCardProps } from 'components/io-card'
import IoCardContainer from 'components/io-card-container'
import s from './style.module.css'

interface IoHomeInPracticeProps {
  brand: Products
  pattern: string
  heading: string
  description: string
  cards: Array<IoCardProps>
  cta: {
    heading: string
    description: string
    link: string
    image: {
      url: string
      alt: string
      width: number
      height: number
    }
  }
}

export default function IoHomeInPractice({
  brand,
  pattern,
  heading,
  description,
  cards,
  cta,
}: IoHomeInPracticeProps) {
  return (
    <section
      className={s.inPractice}
      style={
        {
          '--pattern': `url(${pattern})`,
        } as React.CSSProperties
      }
    >
      <div className={s.container}>
        <IoCardContainer
          theme="dark"
          heading={heading}
          description={description}
          cardsPerRow={3}
          cards={cards}
        />

        {cta.heading ? (
          <div className={s.inPracticeCta}>
            <div className={s.inPracticeCtaContent}>
              <h3 className={s.inPracticeCtaHeading}>{cta.heading}</h3>
              {cta.description ? (
                <p className={s.inPracticeCtaDescription}>{cta.description}</p>
              ) : null}
              {cta.link ? (
                <Button
                  title="Learn more"
                  url={cta.link}
                  theme={{
                    brand: brand,
                  }}
                />
              ) : null}
            </div>
            {cta.image?.url ? (
              <div className={s.inPracticeCtaMedia}>
                <Image
                  src={cta.image.url}
                  width={cta.image.width}
                  height={cta.image.height}
                  alt={cta.image.alt}
                />
              </div>
            ) : null}
          </div>
        ) : null}
      </div>
    </section>
  )
}
