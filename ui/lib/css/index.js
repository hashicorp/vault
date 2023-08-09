/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

/* eslint-env node */
/* eslint-disable n/no-extraneous-require */
'use strict';

var path = require('path');
var Funnel = require('broccoli-funnel');
var mergeTrees = require('broccoli-merge-trees');

module.exports = {
  name: 'sassSvgUri',

  isDevelopingAddon() {
    return true;
  },

  included: function (app) {
    this._super.included.apply(this, arguments);

    // see: https://github.com/ember-cli/ember-cli/issues/3718
    while (typeof app.import !== 'function' && app.app) {
      app = app.app;
    }

    this.sassSVGURIPath = path.dirname(require.resolve('sass-svg-uri'));
    return app;
  },

  treeForStyles: function () {
    var sassSVGURI = new Funnel(this.sassSVGURIPath, {
      srcDir: '/',
      destDir: 'app/styles/sass-svg-uri',
      annotation: 'Sass SVG URI',
    });

    return mergeTrees([sassSVGURI], { overwrite: true });
  },
};
