import * as React from 'react'
import Link from 'next/link'
import InlineSvg from '@hashicorp/react-inline-svg'
import classNames from 'classnames'
import { IconArrowRight24 } from '@hashicorp/flight-icons/svg-react/arrow-right-24'
import { IconExternalLink24 } from '@hashicorp/flight-icons/svg-react/external-link-24'
import { productLogos } from './product-logos'
import s from './style.module.css'

export interface IoCardProps {
  variant?: 'light' | 'gray' | 'dark'
  products?: Array<{
    name: keyof typeof productLogos
  }>
  link: {
    url: string
    type: 'inbound' | 'outbound'
  }
  inset?: 'none' | 'sm' | 'md'
  eyebrow?: string
  heading?: string
  description?: string
  children?: React.ReactNode
}

function IoCard({
  variant = 'light',
  products,
  link,
  inset = 'md',
  eyebrow,
  heading,
  description,
  children,
}: IoCardProps): React.ReactElement {
  const LinkWrapper = ({ className, children }) =>
    link.type === 'inbound' ? (
      <Link href={link.url}>
        <a className={className}>{children}</a>
      </Link>
    ) : (
      <a
        className={className}
        href={link.url}
        target="_blank"
        rel="noopener noreferrer"
      >
        {children}
      </a>
    )

  return (
    <article className={classNames(s.card)}>
      <LinkWrapper className={classNames(s[variant], s[inset])}>
        {children ? (
          children
        ) : (
          <>
            {eyebrow ? <Eyebrow>{eyebrow}</Eyebrow> : null}
            {heading ? <Heading>{heading}</Heading> : null}
            {description ? <Description>{description}</Description> : null}
          </>
        )}
        <footer className={s.footer}>
          {products && (
            <ul className={s.products}>
              {products.map(({ name }, index) => {
                const key = name.toLowerCase()
                const version = variant === 'dark' ? 'neutral' : 'color'
                return (
                  // eslint-disable-next-line react/no-array-index-key
                  <li key={index}>
                    <InlineSvg
                      className={s.logo}
                      src={productLogos[key][version]}
                    />
                  </li>
                )
              })}
            </ul>
          )}
          <span className={s.linkType}>
            {link.type === 'inbound' ? (
              <IconArrowRight24 />
            ) : (
              <IconExternalLink24 />
            )}
          </span>
        </footer>
      </LinkWrapper>
    </article>
  )
}

interface EyebrowProps {
  children: string
}

function Eyebrow({ children }: EyebrowProps) {
  return <p className={s.eyebrow}>{children}</p>
}

interface HeadingProps {
  as?: 'h2' | 'h3' | 'h4'
  children: React.ReactNode
}

function Heading({ as: Component = 'h2', children }: HeadingProps) {
  return <Component className={s.heading}>{children}</Component>
}

interface DescriptionProps {
  children: string
}

function Description({ children }: DescriptionProps) {
  return <p className={s.description}>{children}</p>
}

IoCard.Eyebrow = Eyebrow
IoCard.Heading = Heading
IoCard.Description = Description

export default IoCard
