import * as React from 'react'
import rivetQuery from '@hashicorp/nextjs-scripts/dato/client'
import homepageQuery from './query.graphql'
import { isInternalLink } from 'lib/utils'
import IoHomeHero from 'components/io-home-hero'
import IoHomeVideoCallout from 'components/io-home-video-callout'
import IoHomeCard from 'components/io-home-card'
import IoHomeFeature from 'components/io-home-feature'
import IoHomeCaseStudies from 'components/io-home-case-studies'
import IoHomeCallToAction from 'components/io-home-call-to-action'
import IoHomePreFooter from 'components/io-home-pre-footer'
import s from './style.module.css'

export default function Homepage({ data }) {
  const {
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
              <li key={index}>
                <div className={s.container}>
                  <IoHomeFeature {...feature} />
                </div>
              </li>
            )
          })}
        </ul>

        <div className={s.container}>
          <IoHomeVideoCallout {...intro.video} />
        </div>
      </section>

      <section className={s.inPractice}>
        <div className={s.container}>
          <header className={s.inPracticeHeader}>
            <h2 className={s.inPracticeHeading}>{inPractice.heading}</h2>
            <p className={s.inPracticeDescription}>{inPractice.description}</p>
          </header>
          <ul className={s.inPracticeCards}>
            {inPractice.cards.map(
              ({ link, eyebrow, heading, description, products }, index) => {
                return (
                  <li key={index}>
                    <IoHomeCard
                      variant="dark"
                      link={{
                        url: link,
                        type: isInternalLink(link) ? 'inbound' : 'outbound',
                      }}
                      eyebrow={eyebrow}
                      heading={heading}
                      description={description}
                      products={products}
                    />
                  </li>
                )
              }
            )}
          </ul>
        </div>
      </section>

      <section className={s.useCases}>
        <div className={s.container}>
          <header className={s.useCasesHeader}>
            <h2 className={s.useCasesHeading}>{useCases.heading}</h2>
          </header>

          <ul className={s.useCasesCards}>
            {useCases.cards.map(
              ({ link, eyebrow, heading, description, products }, index) => {
                return (
                  <li key={index}>
                    <IoHomeCard
                      link={{
                        url: link,
                        type: isInternalLink(link) ? 'inbound' : 'outbound',
                      }}
                      inset="sm"
                      eyebrow={eyebrow}
                      heading={heading}
                      description={description}
                      products={products}
                    />
                  </li>
                )
              }
            )}
          </ul>
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
