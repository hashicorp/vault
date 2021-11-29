import Image from 'next/image'
import * as React from 'react'
import classNames from 'classnames'
import Button from '@hashicorp/react-button'
import s from './style.module.css'

interface IoUsecaseCallToActionProps {
  brand: string
  theme?: 'light' | 'dark'
  heading: string
  description: string
  links: Array<{
    text: string
    url: string
  }>
  // TODO document intended usage
  pattern: string
}

export default function IoUsecaseCallToAction({
  brand,
  theme,
  heading,
  description,
  links,
  pattern,
}: IoUsecaseCallToActionProps): React.ReactElement {
  return (
    <div
      className={classNames(s.callToAction, s[theme])}
      style={
        {
          '--background-color': `var(--${brand})`,
        } as React.CSSProperties
      }
    >
      <h2 className={s.heading}>{heading}</h2>
      <div className={s.content}>
        <p className={s.description}>{description}</p>
        <div className={s.links}>
          {links.map((link, index) => {
            return (
              <Button
                // Index is stable
                // eslint-disable-next-line react/no-array-index-key
                key={index}
                title={link.text}
                url={link.url}
                theme={{
                  brand: 'neutral',
                  variant: index === 0 ? 'primary' : 'secondary',
                  background: theme,
                }}
              />
            )
          })}
        </div>
      </div>
      <div className={s.pattern}>
        <Image
          src={pattern}
          layout="fill"
          objectFit="cover"
          objectPosition="center left"
          alt=""
        />
      </div>
    </div>
  )
}
