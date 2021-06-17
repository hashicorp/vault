const fs = require('fs')
const globby = require('globby')
const prettier = require('prettier')

/** Based on this blog by Lee Robinson: https://leerob.io/blog/nextjs-sitemap-robots */

;(async () => {
  const prettierConfig = await prettier.resolveConfig('./.prettierrc.js')

  // Ignore Next.js specific files (e.g., _app.js) and API routes.
  const pages = await globby([
    'pages/**/*{.js,.tsx,.jsx}',
    //  'content/**/*.mdx', // docs routes
    '!pages/**/[*', // ignore dynamic routes
    '!pages/home/*', // home is reexported at /index
    '!pages/_*.js*', // ignore nextjs specific pages
    '!pages/api', // ignore api routes
  ])

  const sitemap = `
          <?xml version="1.0" encoding="UTF-8"?>
          <urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
              ${pages
                .map((page) => {
                  const path = page
                    .replace('pages/', '')
                    .replace(/(index)?.(jsx?|tsx?|mdx)$/, '')

                  return `
                          <url>
                              <loc>${`https://www.vaultproject.io/${path}`}</loc>
                          </url>
                      `
                })
                .join('')}
          </urlset>
      `

  // Formats the xml
  const formatted = prettier.format(sitemap, {
    ...prettierConfig,
    parser: 'html',
  })

  fs.writeFileSync('public/sitemap.xml', formatted)
})()
