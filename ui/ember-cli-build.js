/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/* eslint-env node */
'use strict';

const EmberApp = require('ember-cli/lib/broccoli/ember-app');
const config = require('./config/environment')();

const environment = EmberApp.env();
const isProd = environment === 'production';
const isTest = environment === 'test';
// const isCI = !!process.env.CI;

const appConfig = {
  'ember-service-worker': {
    serviceWorkerScope: config.serviceWorkerScope,
    skipWaitingOnMessage: true,
  },
  babel: {
    plugins: [require.resolve('ember-concurrency/async-arrow-task-transform')],
  },
  svgJar: {
    optimizer: {},
    sourceDirs: ['public'],
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
  'ember-cli-babel': {
    enableTypeScriptTransform: true,
    throwUnlessParallelizable: true,
  },
  hinting: isTest,
  tests: isTest,
  sourcemaps: {
    enabled: !isProd,
  },
  sassOptions: {
    sourceMap: false,
    onlyIncluded: true,
    precision: 4,
    includePaths: [
      './node_modules/@hashicorp/design-system-components/dist/styles',
      './node_modules/@hashicorp/ember-flight-icons/dist/styles',
      './node_modules/@hashicorp/design-system-tokens/dist/products/css',
    ],
  },
  minifyCSS: {
    options: {
      advanced: false,
    },
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
  const app = new EmberApp(defaults, appConfig);

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

  return app.toTree();
};
