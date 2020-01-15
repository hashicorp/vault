// This script replaces <br>, which is invalid in react, with <br /> for all markdown files

const glob = require('glob')
const path = require('path')
const fs = require('fs')

glob.sync(path.join(__dirname, '../pages/**/*.mdx')).map(fullPath => {
  let content = fs.readFileSync(fullPath, 'utf8')

  // fix unclosed br tag
  content = content.replace(/<br>/g, '<br />')
  // fix unclosed img tags
  content = content.replace(/(<img[^>]+)(?<!\/)>/g, (_, m1) => `${m1} />`)

  fs.writeFileSync(fullPath, content)
})
