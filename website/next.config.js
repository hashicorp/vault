const withHashicorp = require('@hashicorp/nextjs-scripts')
const path = require('path')

module.exports = withHashicorp({
  defaultLayout: true,
  transpileModules: ['is-absolute-url', '@hashicorp/react-mega-nav'],
  mdx: { resolveIncludes: path.join(__dirname, 'pages') }
})({
  experimental: {
    css: true,
    modern: true,
    polyfillsOptimization: true,
    rewrites: () => [
      {
        source: '/api/:path*',
        destination: '/api-docs/:path*'
      }
    ],
    redirects: () => [
      {
        source: '/intro',
        destination: '/intro/getting-started',
        permanent: false
      }
    ]
  },
  exportTrailingSlash: true,
  env: {
    HASHI_ENV: process.env.HASHI_ENV
  }
})
