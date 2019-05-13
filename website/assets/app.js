const cssStandards = require('spike-css-standards')
const jsStandards = require('spike-js-standards')
const preactPreset = require('babel-preset-preact')
const extendRule = require('postcss-extend-rule')
const webpack = require('webpack')

/* eslint-disable-next-line */
console.log(`Building assets for environment *${process.env.NODE_ENV}*`)

const isProd =
  process.env.NODE_ENV === 'production' ||
  process.env.NODE_ENV === 'tmp-production'

let utilServerUrl
if (isProd) {
  utilServerUrl = 'https://util.hashicorp.com'
} else {
  utilServerUrl = 'https://hashicorp-web-util-staging.herokuapp.com'
}

if (process.env.UTIL_SERVER) {
  utilServerUrl = process.env.UTIL_SERVER

  // remove trailing slash
  utilServerUrl = utilServerUrl.replace(/\/$/, '')

  /* eslint-disable-next-line */
  console.log(`utilServerUrl=${utilServerUrl}`)
}

let segmentWriteKey
if (isProd) {
  segmentWriteKey = 'OdSFDq9PfujQpmkZf03dFpcUlywme4sC'
} else {
  segmentWriteKey = '0EXTgkNx0Ydje2PGXVbRhpKKoe5wtzcE'
}

module.exports = {
  ignore: ['yarn.lock', '**/_*'],
  entry: {
    'js/main': './js/index.js',
    'js/analytics.js': './js/analytics.js',
    'js/consent-manager': './js/consent-manager.js'
  },
  postcss: cssStandards({
    appendPlugins: [extendRule()]
  }),
  plugins: [
    new webpack.DefinePlugin({
      'process.env.NODE_ENV': JSON.stringify(process.env.NODE_ENV),
      utilityServerRoot: JSON.stringify(utilServerUrl),
      segmentWriteKey: JSON.stringify(segmentWriteKey)
    })
  ],
  babel: jsStandards({ appendPresets: [preactPreset] }),
  server: { open: false }
}
