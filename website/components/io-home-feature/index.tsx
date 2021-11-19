import * as React from 'react'
import Image from 'next/image'
import { IconArrowRight16 } from '@hashicorp/flight-icons/svg-react/arrow-right-16'
import s from './style.module.css'

interface IoHomeFeatureProps {
  link: string
  image: {
    url: string
    alt: string
  }
  heading: string
  description: string
}

export default function IoHomeFeature({
  link,
  image,
  heading,
  description,
}: IoHomeFeatureProps): React.ReactElement {
  return (
    <a href={link} className={s.feature}>
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
        <span className={s.featureCta} aria-hidden={true}>
          Learn more{' '}
          <span>
            <IconArrowRight16 />
          </span>
        </span>
      </div>
    </a>
  )
}
