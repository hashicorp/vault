import Hero from '@hashicorp/react-hero'
import Button from '@hashicorp/react-button'
import styles from './HomepageHero.module.css'
import classNames from 'classnames'

/* A simple Facade wrapper around the Hero component */
export default function HomepageHero({
  title,
  description,
  buttons,
  uiVideo,
  cliVideo,
}) {
  return (
    <div className={styles.homepageHero}>
      <Hero
        data={{
          backgroundTheme: 'light',
          buttons: buttons.slice(0, 2),
          centered: false,
          description: description,
          product: 'vault',
          title: title,
          videos: [
            {
              name: 'UI',
              playbackRate: 2,
              src: [
                {
                  srcType: 'mp4',
                  url: uiVideo,
                },
              ],
            },
            {
              name: 'CLI',
              playbackRate: 2,
              src: [
                {
                  srcType: 'mp4',
                  url: cliVideo,
                },
              ],
            },
          ],
        }}
      />
      {/* A hack to inject a third link styled in tertiary style
           this is very much a non-ideal way to handle this. */}
      <div className={classNames('g-grid-container', styles.thirdLinkWrapper)}>
        {buttons[2] && (
          <div className="third-link">
            <Button
              // eslint-disable-next-line react/no-array-index-key
              linkType={buttons[2].type}
              theme={{
                variant: 'tertiary-neutral',
                brand: 'vault',
                background: 'light',
              }}
              title={buttons[2].title}
              url={buttons[2].url}
            />
          </div>
        )}
      </div>
    </div>
  )
}
