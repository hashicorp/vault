/* eslint-env node */
'use strict';

const EmberApp = require('ember-cli/lib/broccoli/ember-app');
const config = require('./config/environment')();

const environment = EmberApp.env();
const isProd = environment === 'production';
const isTest = environment === 'test';
const isCI = !!process.env.CI;

module.exports = function(defaults) {
  var app = new EmberApp(defaults, {
    svgJar: {
      //optimize: false,
      //paths: [],
      optimizer: {},
      sourceDirs: ['node_modules/@hashicorp/structure-icons/dist', 'public'],
      rootURL: '/ui/',
    },
    assetLoader: {
      generateURI: function(filePath) {
        return `${config.rootURL.replace(/\/$/, '')}${filePath}`;
      },
    },

    codemirror: {
      modes: ['javascript', 'ruby'],
      keyMaps: ['sublime'],
    },
    babel: {
      plugins: ['transform-object-rest-spread'],
    },
    'ember-cli-babel': {
      includePolyfill: isTest || isProd || isCI,
    },
    hinting: isTest,
    tests: isTest,
    sourcemaps: {
      enabled: !isProd,
    },
    sassOptions: {
      sourceMap: false,
      onlyIncluded: true,
      implementation: require('node-sass'),
    },
    autoprefixer: {
      enabled: isTest || isProd,
      grid: true,
      browsers: ['defaults', 'ie 11'],
    },
    autoImport: {
      webpack: {
        // this makes `unsafe-eval` CSP unnecessary
        // see https://github.com/ef4/ember-auto-import/issues/50
        // and https://github.com/webpack/webpack/issues/5627
        devtool: 'inline-source-map',
      },
    },
    'ember-test-selectors': {
      strip: isProd,
    },
    // https://github.com/ember-cli/ember-fetch/issues/204
    'ember-fetch': {
      preferNative: true,
    },
  });

  app.import('vendor/string-includes.js');
  app.import('node_modules/string.prototype.endswith/endswith.js');
  app.import('node_modules/string.prototype.startswith/startswith.js');

  app.import('node_modules/jsonlint/lib/jsonlint.js');
  app.import('node_modules/codemirror/addon/lint/lint.css');
  app.import('node_modules/codemirror/addon/lint/lint.js');
  app.import('node_modules/codemirror/addon/lint/json-lint.js');
  app.import('node_modules/text-encoder-lite/index.js');

  app.import('app/styles/bulma/bulma-radio-checkbox.css');

  app.import('node_modules/@hashicorp/structure-icons/dist/loading.css');
  app.import('node_modules/@hashicorp/structure-icons/dist/run.css');
  // Use `app.import` to add additional libraries to the generated
  // output files.
  //
  // If you need to use different assets in different
  // environments, specify an object as the first parameter. That
  // object's keys should be the environment name and the values
  // should be the asset to use in that environment.
  //
  // If the library that you are including contains AMD or ES6
  // modules that you would like to import into your application
  // please specify an object with the list of modules as keys
  // along with the exports of each module as its value.

  return app.toTree();
};
