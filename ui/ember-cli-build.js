/* eslint-env node */
'use strict';

const EmberApp = require('ember-cli/lib/broccoli/ember-app');
const config = require('./config/environment')();

const environment = EmberApp.env();
const isProd = environment === 'production';
const isTest = environment === 'test';
const isCI = !!process.env.CI;

const appConfig = {
  'ember-service-worker': {
    serviceWorkerScope: config.serviceWorkerScope,
    skipWaitingOnMessage: true,
  },
  svgJar: {
    //optimize: false,
    //paths: [],
    optimizer: {},
    sourceDirs: ['node_modules/@hashicorp/structure-icons/dist', 'public'],
    rootURL: '/ui/',
  },
  fingerprint: {
    exclude: ['images/'],
  },
  assetLoader: {
    generateURI: function (filePath) {
      return `${config.rootURL.replace(/\/$/, '')}${filePath}`;
    },
  },
  babel: {
    plugins: ['@babel/plugin-proposal-object-rest-spread', ['inline-json-import', {}]],
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
  },
  autoprefixer: {
    enabled: isTest || isProd,
    grid: true,
    browsers: ['defaults'],
  },
  autoImport: {
    forbidEval: true,
  },
  'ember-test-selectors': {
    strip: isProd,
  },
  'ember-composable-helpers': {
    except: ['array'],
  },
  'ember-cli-deprecation-workflow': {
    enabled: true,
  },
};

module.exports = function (defaults) {
  let app = new EmberApp(defaults, appConfig);

  app.import('vendor/string-includes.js');
  app.import('node_modules/string.prototype.endswith/endswith.js');
  app.import('node_modules/string.prototype.startswith/startswith.js');

  app.import('node_modules/jsonlint/lib/jsonlint.js');
  app.import('node_modules/codemirror/addon/lint/lint.css');
  app.import('node_modules/codemirror/lib/codemirror.css');
  app.import('node_modules/text-encoder-lite/text-encoder-lite.js');
  app.import('node_modules/jsondiffpatch/dist/jsondiffpatch.umd.js');
  app.import('node_modules/jsondiffpatch/dist/formatters-styles/html.css');

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
