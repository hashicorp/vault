const cssStandards = require('spike-css-standards')
const jsStandards = require('spike-js-standards')
const preactPreset = require('babel-preset-preact')

module.exports = {
  ignore: ['yarn.lock', '**/_*'],
  entry: { 'js/main': './js/index.js' },
  postcss: cssStandards(),
  babel: jsStandards({ appendPresets: [preactPreset] }),
  server: { open: false }
}
