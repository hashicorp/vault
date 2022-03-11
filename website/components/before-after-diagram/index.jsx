import Image from '@hashicorp/react-image'
import InlineSvg from '@hashicorp/react-inline-svg'
import alertIcon from 'public/img/icons/alert.svg?include'
import checkIcon from 'public/img/icons/check.svg?include'
import fragment from './fragment.graphql'
import s from './style.module.css'
function BeforeAfterDiagram(props) {
  const {
    theme,
    beforeHeadline,
    beforeContent,
    beforeImage,
    afterHeadline,
    afterContent,
    afterImage,
  } = props
  return (
    <div className={s.beforeAfterDiagram} data-theme={theme}>
      <div className={s.beforeSide}>
        <div className={s.image}>
          <div>
            <Image {...beforeImage} />
          </div>
        </div>
        <div className={s.contentContainer}>
          <span className={s.iconLineContainer}>
            <InlineSvg className={s.beforeIcon} src={alertIcon} />
            <span className={s.lineSegment} />
          </span>
          <div>
            {beforeHeadline && (
              <h2
                className={s.contentHeadline}
                dangerouslySetInnerHTML={{
                  __html: beforeHeadline,
                }}
              />
            )}
            {beforeContent && (
              <div
                className={s.beforeContent}
                dangerouslySetInnerHTML={{
                  __html: beforeContent,
                }}
              />
            )}
          </div>
        </div>
      </div>
      <div className={s.afterSide}>
        <div className={s.image}>
          <div>
            <Image {...afterImage} />
          </div>
        </div>
        <div className={s.contentContainer}>
          <span className={s.iconLineContainer}>
            <InlineSvg className={s.afterIcon} src={checkIcon} />
          </span>
          <div>
            {afterHeadline && (
              <h2
                className={s.contentHeadline}
                dangerouslySetInnerHTML={{
                  __html: afterHeadline,
                }}
              />
            )}
            {afterContent && (
              <div
                className={s.afterContent}
                data-theme={theme}
                dangerouslySetInnerHTML={{
                  __html: afterContent,
                }}
              />
            )}
          </div>
        </div>
      </div>
    </div>
  )
}

BeforeAfterDiagram.fragmentSpec = { fragment, dependencies: [Image] }

export default BeforeAfterDiagram
