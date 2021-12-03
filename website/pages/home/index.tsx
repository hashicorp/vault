import * as React from 'react'
import Head from 'next/head'
import rivetQuery from '@hashicorp/nextjs-scripts/dato/client'
import homepageQuery from './query.graphql'
import { renderMetaTags } from 'react-datocms'
import IoHomeHero from 'components/io-home-hero'
import IoHomeInPractice from 'components/io-home-in-practice'
import IoVideoCallout from 'components/io-video-callout'
import IoCardContainer from 'components/io-card-container'
import IoHomeFeature from 'components/io-home-feature'
import IoHomeCaseStudies from 'components/io-home-case-studies'
import IoHomeCallToAction from 'components/io-home-call-to-action'
import IoHomePreFooter from 'components/io-home-pre-footer'
import s from './style.module.css'

export default function Homepage({ data }): React.ReactElement {
  const {
    seo,
    hero,
    intro,
    inPractice,
    useCases,
    caseStudies,
    callToAction,
    preFooter,
  } = data

  return (
    <>
      <Head>{renderMetaTags(seo)}</Head>

      <IoHomeHero
        pattern="/img/home-hero-pattern.svg"
        brand="vault"
        {...hero}
      />

      <section className={s.intro}>
        <header className={s.introHeader}>
          <div className={s.container}>
            <div className={s.introHeaderInner}>
              <h2 className={s.introHeading}>{intro.heading}</h2>
              <p className={s.introDescription}>{intro.description}</p>
            </div>
          </div>
        </header>

        <ul className={s.features}>
          {intro.features.map((feature, index) => {
            return (
              // Index is stable
              // eslint-disable-next-line react/no-array-index-key
              <li key={index}>
                <div className={s.container}>
                  <IoHomeFeature {...feature} />
                </div>
              </li>
            )
          })}
        </ul>

        {intro.video ? (
          <div className={s.container}>
            <IoVideoCallout
              youtubeId={intro.video.youtubeId}
              thumbnail={intro.video.thumbnail.url}
              heading={intro.video.heading}
              description={intro.video.description}
              person={{
                name: intro.video.personName,
                description: intro.video.personDescription,
                avatar: intro.video.personAvatar?.url,
              }}
            />
          </div>
        ) : null}
      </section>

      <IoHomeInPractice
        brand="vault"
        pattern="/img/practice-pattern.svg"
        heading={inPractice.heading}
        description={inPractice.description}
        cards={inPractice.cards.map((card) => {
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
        cta={{
          heading: inPractice.cta.heading,
          description: inPractice.cta.description,
          link: inPractice.cta.link,
          image: inPractice.cta.image,
        }}
      />

      <section className={s.useCases}>
        <div className={s.container}>
          <IoCardContainer
            heading={useCases.heading}
            description={useCases.description}
            cardsPerRow={4}
            cards={useCases.cards.map((card) => {
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
      </section>

      <section className={s.caseStudies}>
        <div className={s.container}>
          <header className={s.caseStudiesHeader}>
            <h2 className={s.caseStudiesHeading}>{caseStudies.heading}</h2>
            <p className={s.caseStudiesDescription}>
              {caseStudies.description}
            </p>
          </header>

          <IoHomeCaseStudies
            primary={caseStudies.features}
            secondary={caseStudies.links}
          />
        </div>
      </section>

      <IoHomeCallToAction
        brand="vault"
        heading={callToAction.heading}
        content={callToAction.description}
        links={callToAction.links}
      />

      <IoHomePreFooter
        brand="vault"
        heading={preFooter.heading}
        description={preFooter.description}
        ctas={preFooter.ctas}
      />
    </>
  )
}

export async function getStaticProps() {
  const { vaultHomepage } = await rivetQuery({
    query: homepageQuery,
  })

  const {
    seo,
    heroHeading,
    heroDescription,
    heroCtas,
    heroCards,
    introHeading,
    introDescription,
    introFeatures,
    introVideo,
    inPracticeHeading,
    inPracticeDescription,
    inPracticeCards,
    inPracticeCtaHeading,
    inPracticeCtaDescription,
    inPracticeCtaLink,
    inPracticeCtaImage,
    useCasesHeading,
    useCasesDescription,
    useCasesCards,
    caseStudiesHeading,
    caseStudiesDescription,
    caseStudiesFeatured,
    caseStudiesLinks,
    callToActionHeading,
    callToActionDescription,
    callToActionCtas,
    preFooterHeading,
    preFooterDescription,
    preFooterCtas,
  } = vaultHomepage

  return {
    props: {
      data: {
        seo,
        hero: {
          heading: heroHeading,
          description: heroDescription,
          ctas: heroCtas,
          cards: heroCards.map((card) => {
            return {
              ...card,
              cta: card.cta[0],
            }
          }),
        },
        intro: {
          heading: introHeading,
          description: introDescription,
          features: introFeatures,
          video: introVideo[0],
        },
        inPractice: {
          heading: inPracticeHeading,
          description: inPracticeDescription,
          cards: inPracticeCards,
          cta: {
            heading: inPracticeCtaHeading,
            description: inPracticeCtaDescription,
            link: inPracticeCtaLink,
            image: inPracticeCtaImage,
          },
        },
        useCases: {
          heading: useCasesHeading,
          description: useCasesDescription,
          cards: useCasesCards,
        },
        caseStudies: {
          heading: caseStudiesHeading,
          description: caseStudiesDescription,
          features: caseStudiesFeatured,
          links: caseStudiesLinks,
        },
        callToAction: {
          heading: callToActionHeading,
          description: callToActionDescription,
          links: callToActionCtas,
        },
        preFooter: {
          heading: preFooterHeading,
          description: preFooterDescription,
          ctas: preFooterCtas,
        },
      },
    },
  }
}
