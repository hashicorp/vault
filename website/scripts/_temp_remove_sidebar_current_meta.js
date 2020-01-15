// This script removed the "sidebar_current" metadata in frontmatter,
// because its not used at all anymore

const glob = require('glob')
const path = require('path')
const fs = require('fs')
const matter = require('gray-matter')

glob.sync(path.join(__dirname, '../pages/**/*.mdx')).map(fullPath => {
  const { content, data } = matter.read(fullPath)
  delete data.sidebar_current
  fs.writeFileSync(fullPath, matter.stringify(content, data))
})
