import * as React from 'react'
import rivetQuery from '@hashicorp/nextjs-scripts/dato/client'
import useCasesQuery from './query.graphql'
import IoUsecaseHero from 'components/io-usecase-hero'
import IoUsecaseSection from 'components/io-usecase-section'
import IoUsecaseCustomer from 'components/io-usecase-customer'
import IoCardContainer from 'components/io-card-container'
import IoVideoCallout from 'components/io-video-callout'
import IoUsecaseCallToAction from 'components/io-usecase-call-to-action'
import s from './style.module.css'

export default function UseCasePage({ data }) {
  const {
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
    caseStudyImage,
    caseStudyLogo,
    caseStudyHeading,
    caseStudyDescription,
    caseStudyLink,
    caseStudyStats,
    callToActionHeading,
    callToActionDescription,
    callToActionLinks,
    videoCallout,
  } = data
  const _videoCallout = videoCallout[0]
  return (
    <>
      <IoUsecaseHero
        eyebrow="Common use case"
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

      <IoUsecaseCustomer
        link={caseStudyLink}
        media={{
          src: caseStudyImage.url,
          width: caseStudyImage.width,
          height: caseStudyImage.height,
          alt: caseStudyImage.alt,
        }}
        logo={{
          src: caseStudyLogo.url,
          width: caseStudyLogo.width,
          height: caseStudyLogo.height,
          alt: caseStudyLogo.alt,
        }}
        heading={caseStudyHeading}
        description={caseStudyDescription}
        stats={caseStudyStats.map((stat) => {
          return {
            value: stat.value,
            key: stat.label,
          }
        })}
      />

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
              avatar: _videoCallout.personAvatar.url,
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
    fallback: false,
  }
}

export async function getStaticProps({ params }) {
  const { slug } = params

  const { allVaultUseCases } = await rivetQuery({
    query: useCasesQuery,
  })

  const page = allVaultUseCases.find((page) => page.slug === slug)

  return {
    props: {
      data: page,
    },
  }
}
