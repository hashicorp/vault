import Image from 'next/image'
import { IconArrowRight16 } from '@hashicorp/flight-icons/svg-react/arrow-right-16'
import s from './style.module.css'

interface IoHomeCaseStudiesProps {
  primary: Array<{
    thumbnail: string
    alt: string
    link: string
    heading: string
  }>
  secondary: Array<{
    link: string
    heading: string
  }>
}

export default function IoHomeCaseStudies({
  primary,
  secondary,
}: IoHomeCaseStudiesProps) {
  return (
    <div className={s.caseStudies}>
      <ul className={s.primary}>
        {primary.map((item) => {
          return (
            <li className={s.primaryItem}>
              <a className={s.card} href={item.link}>
                <h3 className={s.cardHeading}>{item.heading}</h3>
                <Image
                  src={item.thumbnail}
                  layout="fill"
                  objectFit="cover"
                  alt={item.alt}
                />
              </a>
            </li>
          )
        })}
      </ul>

      <ul className={s.secondary}>
        {secondary.map((item) => {
          return (
            <li className={s.secondaryItem}>
              <a className={s.link} href={item.link}>
                <span className={s.linkInner}>
                  <h3 className={s.linkHeading}>{item.heading}</h3>
                  <IconArrowRight16 />
                </span>
              </a>
            </li>
          )
        })}
      </ul>
    </div>
  )
}
