import * as React from 'react'
import Head from 'next/head'
import rivetQuery from '@hashicorp/nextjs-scripts/dato/client'
import homepageQuery from './query.graphql'
import { renderMetaTags } from 'react-datocms'
import IoHomeHero from 'components/io-home-hero'
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

      <IoHomeHero brand="vault" {...hero} />

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

        <div className={s.container}>
          <IoVideoCallout {...intro.video} />
        </div>
      </section>

      <section className={s.inPractice}>
        <div className={s.container}>
          <IoCardContainer
            theme="dark"
            heading="Vault in practice"
            description="The best way to understand what Vault can enable for your projects is to see it in action"
            cardsPerRow={3}
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
          />
        </div>
      </section>

      <section className={s.useCases}>
        <div className={s.container}>
          <IoCardContainer
            heading={useCases.heading}
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
        links={callToAction.links.map(({ text, url }, index) => {
          return {
            text,
            url,
            type: index === 1 ? 'inbound' : null,
          }
        })}
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
    inPracticeHeading,
    inPracticeDescription,
    inPracticeCards,
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
          video: {
            youtubeId: vaultHomepage.introVideo[0].youtubeId,
            heading: vaultHomepage.introVideo[0].heading,
            description: vaultHomepage.introVideo[0].description,
            thumbnail: vaultHomepage.introVideo[0].thumbnail.url,
            person: {
              name: vaultHomepage.introVideo[0].personName,
              description: vaultHomepage.introVideo[0].personDescription,
              avatar: vaultHomepage.introVideo[0].personAvatar.url,
            },
          },
        },
        inPractice: {
          heading: inPracticeHeading,
          description: inPracticeDescription,
          cards: inPracticeCards,
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
