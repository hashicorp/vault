import * as React from 'react'
import Image from 'next/image'
import Link from 'next/link'
import { IconArrowRight16 } from '@hashicorp/flight-icons/svg-react/arrow-right-16'
import s from './style.module.css'

export interface IoHomeFeatureProps {
  isInternalLink: (link: string) => boolean
  link?: string
  image: {
    url: string
    alt: string
  }
  heading: string
  description: string
}

export default function IoHomeFeature({
  isInternalLink,
  link,
  image,
  heading,
  description,
}: IoHomeFeatureProps): React.ReactElement {
  return (
    <IoHomeFeatureWrap isInternalLink={isInternalLink} href={link}>
      <div className={s.featureMedia}>
        <Image
          src={image.url}
          width={400}
          height={200}
          layout="responsive"
          alt={image.alt}
        />
      </div>
      <div className={s.featureContent}>
        <h3 className={s.featureHeading}>{heading}</h3>
        <p className={s.featureDescription}>{description}</p>
        {link ? (
          <span className={s.featureCta} aria-hidden={true}>
            Learn more{' '}
            <span>
              <IconArrowRight16 />
            </span>
          </span>
        ) : null}
      </div>
    </IoHomeFeatureWrap>
  )
}

interface IoHomeFeatureWrapProps {
  isInternalLink: (link: string) => boolean
  href: string
  children: React.ReactNode
}

function IoHomeFeatureWrap({
  isInternalLink,
  href,
  children,
}: IoHomeFeatureWrapProps) {
  if (!href) {
    return <div className={s.feature}>{children}</div>
  }

  if (isInternalLink(href)) {
    return (
      <Link href={href}>
        <a className={s.feature}>{children}</a>
      </Link>
    )
  }

  return (
    <a className={s.feature} href={href}>
      {children}
    </a>
  )
}
