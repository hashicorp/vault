import * as React from 'react'
import Link from 'next/link'
import InlineSvg from '@hashicorp/react-inline-svg'
import camelCase from 'camelcase'
import classNames from 'classnames'
import { IconArrowRight16 } from '@hashicorp/flight-icons/svg-react/arrow-right-16'
import { IconExternalLink16 } from '@hashicorp/flight-icons/svg-react/external-link-16'
import { productLogos } from './product-logos'
import s from './style.module.css'

interface CardProps {
  variant?: 'light' | 'gray' | 'dark'
  products?: Array<keyof typeof productLogos>
  link: {
    url: string
    type: 'inbound' | 'outbound'
  }
  inset?:
    | 'none'
    | 'sm'
    | 'md'
    | {
        horizontal: 'none' | 'sm' | 'md'
        vertical: 'none' | 'sm' | 'md'
      }
  children: React.ReactNode
}

function Card({
  variant = 'light',
  products,
  link,
  inset = 'md',
  children,
}: CardProps) {
  const space =
    typeof inset === 'string' ? camelCase(['space', inset.toString()]) : null
  const spaceHorizontal =
    typeof inset === 'object' && inset.horizontal
      ? camelCase(['space', 'horizontal', inset.horizontal])
      : null
  const spaceVertical =
    typeof inset === 'object' && inset.vertical
      ? camelCase(['space', 'vertical', inset.vertical])
      : null

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
      <LinkWrapper
        className={classNames(
          s[variant],
          s[space],
          s[spaceHorizontal],
          s[spaceVertical]
        )}
      >
        {children}
        <footer className={s.footer}>
          {products && (
            <ul className={s.products}>
              {products.map((product) => {
                return (
                  <li>
                    <InlineSvg className={s.logo} src={productLogos[product]} />
                  </li>
                )
              })}
            </ul>
          )}
          <span className={s.linkType}>
            {link.type === 'inbound' ? (
              <IconArrowRight16 />
            ) : (
              <IconExternalLink16 />
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

Card.Eyebrow = Eyebrow
Card.Heading = Heading
Card.Description = Description

export default Card
