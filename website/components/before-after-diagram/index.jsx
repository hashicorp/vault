import marked from 'marked'
import Image from '@hashicorp/react-image'
import alertIcon from './img/alert-icon.svg?include'
import checkIcon from './img/check-icon.svg?include'

function BeforeAfterDiagram(props) {
  const markedOptions = generateMarkedOptions()
  const {
    theme,
    beforeHeadline,
    beforeContent,
    beforeImage,
    afterHeadline,
    afterContent,
    afterImage
  } = props
  return (
    <div className={`g-before-after-diagrams ${theme}`}>
      <div className="before">
        <div className="image">
          <div>
            <Image {...beforeImage} />
          </div>
        </div>
        <div className="content">
          <span className="line">
            <span />
            <div
              dangerouslySetInnerHTML={{
                __html: alertIcon
              }}
            />
            <span />
          </span>
          <div>
            {beforeHeadline && (
              <h3
                className="g-type-display-3"
                dangerouslySetInnerHTML={{
                  __html: marked.inlineLexer(beforeHeadline, [])
                }}
              />
            )}
            {beforeContent && (
              <div
                className="g-type-body-small"
                dangerouslySetInnerHTML={{
                  __html: marked(beforeContent, markedOptions)
                }}
              />
            )}
          </div>
        </div>
      </div>
      <div className="after">
        <div className="image">
          <div>
            <Image {...afterImage} />
          </div>
        </div>
        <div className="content">
          <div className="line">
            <div
              dangerouslySetInnerHTML={{
                __html: checkIcon
              }}
            />
          </div>
          <div>
            {afterHeadline && (
              <h3
                className="g-type-display-3"
                dangerouslySetInnerHTML={{
                  __html: marked.inlineLexer(afterHeadline, [])
                }}
              />
            )}
            {afterContent && (
              <div
                dangerouslySetInnerHTML={{
                  __html: marked(afterContent, markedOptions)
                }}
              />
            )}
          </div>
        </div>
      </div>
    </div>
  )
}

export default BeforeAfterDiagram

function generateMarkedOptions() {
  const markedRenderer = new marked.Renderer()

  markedRenderer.heading = function(text, level) {
    return `<h${level} class="g-type-label">${text}</h${level}>`
  }
  markedRenderer.paragraph = function(text) {
    return `<p class="g-type-body-small">${text}</p>`
  }
  markedRenderer.list = function(text) {
    return `<ul class="g-type-body-small">${text}</ul>`
  }

  return { renderer: markedRenderer }
}
