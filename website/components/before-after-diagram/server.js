import md from '@hashicorp/platform-markdown-utils/markdown-to-html'
import mdInline from '@hashicorp/platform-markdown-utils/markdown-to-inline-html'

export default async function processBeforeAfterDiagramProps(props) {
  const { beforeHeadline, beforeContent, afterHeadline, afterContent } = props
  //  Transform headline markdown to HTML, for inline bold / italic support
  props.beforeHeadline = await mdInline(beforeHeadline)
  props.afterHeadline = await mdInline(afterHeadline)
  //  Transform content markdown to HTML, using custom type classes
  const contentOptions = {
    contentPlugins: {
      pluginOptions: {
        typography: {
          map: {
            h1: 'g-type-label',
            h2: 'g-type-label',
            h3: 'g-type-label',
            h4: 'g-type-label',
            h5: 'g-type-label',
            h6: 'g-type-label',
            p: 'g-type-body-small',
            li: 'g-type-body-small',
          },
        },
      },
    },
  }
  props.beforeContent = await md(beforeContent, contentOptions)
  props.afterContent = await md(afterContent, contentOptions)
  return props
}
