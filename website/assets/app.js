const cssStandards = require('spike-css-standards')
const jsStandards = require('spike-js-standards')
const preactPreset = require('babel-preset-preact')
const extendRule = require('postcss-extend-rule')

module.exports = {
  ignore: ['yarn.lock', '**/_*'],
  entry: { 'js/main': './js/index.js' },
  postcss: cssStandards({
    appendPlugins: [extendRule()]
  }),
  babel: jsStandards({ appendPresets: [preactPreset] }),
  server: { open: false }
}
