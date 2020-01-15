// This script fixes any erb-based markdown includes to use our markdown plugin

const glob = require('glob')
const path = require('path')
const fs = require('fs')
const matter = require('gray-matter')

glob.sync(path.join(__dirname, '../pages/**/*.mdx')).map(fullPath => {
  let { content, data } = matter.read(fullPath)
  content = content.replace(
    /<%=\s*partial[(\s]["'](.*)["'][)\s]\s*%>/gm,
    (_, partialPath) => `@include '${partialPath}.mdx'`
  )
  fs.writeFileSync(fullPath, matter.stringify(content, data))
})
