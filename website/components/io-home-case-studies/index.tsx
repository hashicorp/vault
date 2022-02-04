import * as React from 'react'
import Image from 'next/image'
import { isInternalLink } from 'lib/utils'
import { IconExternalLink16 } from '@hashicorp/flight-icons/svg-react/external-link-16'
import { IconArrowRight16 } from '@hashicorp/flight-icons/svg-react/arrow-right-16'
import s from './style.module.css'

interface IoHomeCaseStudiesProps {
  heading: string
  description: string
  primary: Array<{
    thumbnail: {
      url: string
      alt: string
    }
    link: string
    heading: string
  }>
  secondary: Array<{
    link: string
    heading: string
  }>
}

export default function IoHomeCaseStudies({
  heading,
  description,
  primary,
  secondary,
}: IoHomeCaseStudiesProps): React.ReactElement {
  return (
    <section className={s.root}>
      <div className={s.container}>
        <header className={s.header}>
          <h2 className={s.heading}>{heading}</h2>
          <p className={s.description}>{description}</p>
        </header>
        <div className={s.caseStudies}>
          <ul className={s.primary}>
            {primary.map((item, index) => {
              return (
                <li key={index} className={s.primaryItem}>
                  <a className={s.card} href={item.link}>
                    <h3 className={s.cardHeading}>{item.heading}</h3>
                    <Image
                      className={s.cardThumbnail}
                      src={item.thumbnail.url}
                      layout="fill"
                      objectFit="cover"
                      alt={item.thumbnail.alt}
                    />
                  </a>
                </li>
              )
            })}
          </ul>

          <ul className={s.secondary}>
            {secondary.map((item, index) => {
              return (
                <li key={index} className={s.secondaryItem}>
                  <a className={s.link} href={item.link}>
                    <span className={s.linkInner}>
                      <h3 className={s.linkHeading}>{item.heading}</h3>
                      {isInternalLink(item.link) ? (
                        <IconArrowRight16 />
                      ) : (
                        <IconExternalLink16 />
                      )}
                    </span>
                  </a>
                </li>
              )
            })}
          </ul>
        </div>
      </div>
    </section>
  )
}
