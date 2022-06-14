import * as React from 'react'
import Image from 'next/image'
import classNames from 'classnames'
import { Products } from '@hashicorp/platform-product-meta'
import Button from '@hashicorp/react-button'
import IoVideoCallout, {
  IoHomeVideoCalloutProps,
} from 'components/io-video-callout'
import IoHomeFeature, { IoHomeFeatureProps } from 'components/io-home-feature'
import s from './style.module.css'

interface IoHomeIntroProps {
  brand: Products
  heading: string
  description: string
  features?: Array<IoHomeFeatureProps>
  offerings?: {
    image: {
      src: string
      width: number
      height: number
      alt: string
    }
    list: Array<{
      heading: string
      description: string
    }>
    cta?: {
      title: string
      link: string
    }
  }
  video?: IoHomeVideoCalloutProps
}

export default function IoHomeIntro({
  brand,
  heading,
  description,
  features,
  offerings,
  video,
}: IoHomeIntroProps) {
  return (
    <section
      className={classNames(
        s.root,
        s[brand],
        features && s.withFeatures,
        offerings && s.withOfferings
      )}
      style={
        {
          '--brand': `var(--${brand})`,
        } as React.CSSProperties
      }
    >
      <header className={s.header}>
        <div className={s.container}>
          <div className={s.headerInner}>
            <h2 className={s.heading}>{heading}</h2>
            <p className={s.description}>{description}</p>
          </div>
        </div>
      </header>

      {features ? (
        <ul className={s.features}>
          {features.map((feature, index) => {
            return (
              // Index is stable
              // eslint-disable-next-line react/no-array-index-key
              <li key={index}>
                <div className={s.container}>
                  <IoHomeFeature
                    image={{
                      url: feature.image.url,
                      alt: feature.image.alt,
                    }}
                    heading={feature.heading}
                    description={feature.description}
                    link={feature.link}
                  />
                </div>
              </li>
            )
          })}
        </ul>
      ) : null}

      {offerings ? (
        <div className={s.offerings}>
          {offerings.image ? (
            <div className={s.offeringsMedia}>
              <Image
                src={offerings.image.src}
                width={offerings.image.width}
                height={offerings.image.height}
                alt={offerings.image.alt}
              />
            </div>
          ) : null}
          <div className={s.offeringsContent}>
            <ul className={s.offeringsList}>
              {offerings.list.map((offering, index) => {
                return (
                  // Index is stable
                  // eslint-disable-next-line react/no-array-index-key
                  <li key={index}>
                    <h3 className={s.offeringsListHeading}>
                      {offering.heading}
                    </h3>
                    <p className={s.offeringsListDescription}>
                      {offering.description}
                    </p>
                  </li>
                )
              })}
            </ul>
            {offerings.cta ? (
              <div className={s.offeringsCta}>
                <Button
                  title={offerings.cta.title}
                  url={offerings.cta.link}
                  theme={{
                    brand: 'neutral',
                  }}
                />
              </div>
            ) : null}
          </div>
        </div>
      ) : null}

      {video ? (
        <div className={s.video}>
          <IoVideoCallout
            youtubeId={video.youtubeId}
            thumbnail={video.thumbnail}
            heading={video.heading}
            description={video.description}
            person={{
              name: video.person.name,
              description: video.person.description,
              avatar: video.person.avatar,
            }}
          />
        </div>
      ) : null}
    </section>
  )
}
