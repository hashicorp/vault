const { h, Component } = require('preact')
const { decode } = require('reshape-preact-components')
const marked = require('8fold-marked')
const Image = require('@hashicorp/hashi-image').default
const AlertIcon = require('./alert-icon')
const CheckIcon = require('./check-icon')

module.exports = class BeforeAfterDiagram extends Component {
  render() {
    const data = decode(this.props._data)

    return (
      <div class={`g-before-after-diagrams ${data.theme}`}>
        <div class="before">
          <div class="image">
            <div>
              <Image src={data.before_image.url} svg="true" />
            </div>
          </div>
          <div class="content">
            <span class="line">
              <span />
              <AlertIcon />
              <span />
            </span>
            <div>
              {data.before_headline && (
                <h3
                  dangerouslySetInnerHTML={{
                    __html: marked.inlineLexer(data.before_headline, [])
                  }}
                />
              )}
              {data.before_content && (
                <div
                  dangerouslySetInnerHTML={{
                    __html: marked(data.before_content)
                  }}
                />
              )}
            </div>
          </div>
        </div>
        <div class="after">
          <div class="image">
            <div>
              <Image src={data.after_image.url} svg="true" />
            </div>
          </div>
          <div class="content">
            <div class="line">
              <CheckIcon />
            </div>
            <div>
              {data.after_headline && (
                <h3
                  dangerouslySetInnerHTML={{
                    __html: marked.inlineLexer(data.after_headline, [])
                  }}
                />
              )}
              {data.after_content && (
                <div
                  dangerouslySetInnerHTML={{
                    __html: marked(data.after_content)
                  }}
                />
              )}
            </div>
          </div>
        </div>
      </div>
    )
  }
}
