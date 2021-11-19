const withHashicorp = require('@hashicorp/platform-nextjs-plugin')
const redirects = require('./redirects.next')

module.exports = withHashicorp({
  dato: {
    // This token is safe to be in this public repository, it only has access to content that is publicly viewable on the website
    token: '88b4984480dad56295a8aadae6caad',
  },
  nextOptimizedImages: true,
  transpileModules: ['@hashicorp/flight-icons'],
})({
  svgo: { plugins: [{ removeViewBox: false }] },
  rewrites: () => [
    {
      source: '/api/:path*',
      destination: '/api-docs/:path*',
    },
  ],
  redirects: () => redirects,
  env: {
    HASHI_ENV: process.env.HASHI_ENV || 'development',
    SEGMENT_WRITE_KEY: 'OdSFDq9PfujQpmkZf03dFpcUlywme4sC',
    BUGSNAG_CLIENT_KEY: '07ff2d76ce27aded8833bf4804b73350',
    BUGSNAG_SERVER_KEY: 'fb2dc40bb48b17140628754eac6c1b11',
  },
  images: {
    domains: ['www.datocms-assets.com'],
    disableStaticImages: true,
  },
})
