import * as React from 'react'
import Head from 'next/head'
import rivetQuery from '@hashicorp/nextjs-scripts/dato/client'
import useCasesQuery from './query.graphql'
import { renderMetaTags } from 'react-datocms'
import IoUsecaseHero from 'components/io-usecase-hero'
import IoUsecaseSection from 'components/io-usecase-section'
import IoUsecaseCustomer from 'components/io-usecase-customer'
import IoCardContainer from 'components/io-card-container'
import IoVideoCallout from 'components/io-video-callout'
import IoUsecaseCallToAction from 'components/io-usecase-call-to-action'
import s from './style.module.css'

export default function UseCasePage({ data }) {
  const {
    seo,
    heroHeading,
    heroDescription,
    challengeHeading,
    challengeDescription,
    challengeImage,
    challengeLink,
    solutionHeading,
    solutionDescription,
    solutionImage,
    solutionLink,
    customerCaseStudy,
    cardsHeading,
    cardsDescription,
    tutorialsLink,
    tutorialCards,
    documentationLink,
    documentationCards,
    callToActionHeading,
    callToActionDescription,
    callToActionLinks,
    videoCallout,
  } = data
  const _customerCaseStudy = customerCaseStudy[0]
  const _videoCallout = videoCallout[0]

  return (
    <>
      <Head>{renderMetaTags(seo)}</Head>

      <IoUsecaseHero
        eyebrow="Use case"
        heading={heroHeading}
        description={heroDescription}
        pattern="/img/usecase-hero-pattern.svg"
      />

      <IoUsecaseSection
        brand="vault"
        eyebrow="Challenge"
        heading={challengeHeading}
        description={challengeDescription}
        media={{
          src: challengeImage?.url,
          width: challengeImage?.width,
          height: challengeImage?.height,
          alt: challengeImage?.alt,
        }}
        cta={{
          text: 'Learn more',
          link: challengeLink,
        }}
      />

      <IoUsecaseSection
        brand="vault"
        bottomIsFlush={_customerCaseStudy}
        eyebrow="Solution"
        heading={solutionHeading}
        description={solutionDescription}
        media={{
          src: solutionImage?.url,
          width: solutionImage?.width,
          height: solutionImage?.height,
          alt: solutionImage?.alt,
        }}
        cta={{
          text: 'Learn more',
          link: solutionLink,
        }}
      />

      {_customerCaseStudy ? (
        <IoUsecaseCustomer
          link={_customerCaseStudy.link}
          media={{
            src: _customerCaseStudy.image.url,
            width: _customerCaseStudy.image.width,
            height: _customerCaseStudy.image.height,
            alt: _customerCaseStudy.image.alt,
          }}
          logo={{
            src: _customerCaseStudy.logo.url,
            width: _customerCaseStudy.logo.width,
            height: _customerCaseStudy.logo.height,
            alt: _customerCaseStudy.logo.alt,
          }}
          heading={_customerCaseStudy.heading}
          description={_customerCaseStudy.description}
          stats={_customerCaseStudy.stats.map((stat) => {
            return {
              value: stat.value,
              key: stat.label,
            }
          })}
        />
      ) : null}

      <div className={s.cards}>
        <IoCardContainer
          heading={cardsHeading}
          description={cardsDescription}
          label="Docs"
          cta={{
            url: documentationLink ? documentationLink : '/docs',
            text: 'Explore all',
          }}
          cardsPerRow={3}
          cards={documentationCards.map((card) => {
            return {
              eyebrow: card.eyebrow,
              link: {
                url: card.link,
                type: 'inbound',
              },
              heading: card.heading,
              description: card.description,
              products: card.products,
            }
          })}
        />

        <IoCardContainer
          label="Tutorials"
          cta={{
            url: tutorialsLink
              ? tutorialsLink
              : 'https://learn.hashicorp.com/vault',
            text: 'Explore all',
          }}
          cardsPerRow={3}
          cards={tutorialCards.map((card) => {
            return {
              eyebrow: card.eyebrow,
              link: {
                url: card.link,
                type: 'inbound',
              },
              heading: card.heading,
              description: card.description,
              products: card.products,
            }
          })}
        />
      </div>

      <div className={s.callToAction}>
        <IoUsecaseCallToAction
          theme="light"
          brand="vault"
          heading={callToActionHeading}
          description={callToActionDescription}
          links={callToActionLinks.map((link) => {
            return {
              text: link.title,
              url: link.link,
            }
          })}
          pattern="/img/usecase-callout-pattern.svg"
        />
      </div>

      {_videoCallout ? (
        <div className={s.videoCallout}>
          <IoVideoCallout
            youtubeId={_videoCallout.youtubeId}
            thumbnail={_videoCallout.thumbnail.url}
            heading={_videoCallout.heading}
            description={_videoCallout.description}
            person={{
              avatar: _videoCallout.personAvatar?.url,
              name: _videoCallout.personName,
              description: _videoCallout.personDescription,
            }}
          />
        </div>
      ) : null}
    </>
  )
}

export async function getStaticPaths() {
  const { allVaultUseCases } = await rivetQuery({
    query: useCasesQuery,
  })

  return {
    paths: allVaultUseCases.map((page) => {
      return {
        params: {
          slug: page.slug,
        },
      }
    }),
    fallback: 'blocking',
  }
}

export async function getStaticProps({ params }) {
  const { slug } = params

  const { allVaultUseCases } = await rivetQuery({
    query: useCasesQuery,
  })

  const page = allVaultUseCases.find((page) => page.slug === slug)

  if (!page) {
    return { notFound: true }
  }

  return {
    props: {
      data: page,
    },
    revalidate:
      process.env.HASHI_ENV === 'production'
        ? process.env.GLOBAL_REVALIDATE
        : 10,
  }
}
