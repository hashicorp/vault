import * as React from 'react'
import Link from 'next/link'
import camelCase from 'camelcase'
import classNames from 'classnames'
import { IconArrowRight16 } from '@hashicorp/flight-icons/svg-react/arrow-right-16'
import { IconExternalLink16 } from '@hashicorp/flight-icons/svg-react/external-link-16'
import s from './style.module.css'

function Card({
  variant = 'light',
  link,
  inset = 'md',
  children,
}: {
  variant: 'light' | 'gray' | 'dark'
  link: {
    url: string
    type: 'inbound' | 'outbound'
  }
  inset:
    | 'none'
    | 'sm'
    | 'md'
    | {
        horizontal: 'none' | 'sm' | 'md'
        vertical: 'none' | 'sm' | 'md'
      }
  children: React.ReactNode
}) {
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
          {link.type === 'inbound' ? (
            <IconArrowRight16 />
          ) : (
            <IconExternalLink16 />
          )}
        </footer>
      </LinkWrapper>
    </article>
  )
}

function Eyebrow({ children }: { children: string }) {
  return <p className={s.eyebrow}>{children}</p>
}

function Heading({
  as: Component = 'h2',
  children,
}: {
  as: 'h2' | 'h3' | 'h4'
  children: React.ReactNode
}) {
  return <Component className={s.heading}>{children}</Component>
}

function Description({ children }: { children: string }) {
  return <p className={s.description}>{children}</p>
}

Card.Eyebrow = Eyebrow
Card.Heading = Heading
Card.Description = Description

export default Card
