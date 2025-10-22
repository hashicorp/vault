/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/* eslint-env node */
'use strict';

const path = require('path');

module.exports = {
  entry: {
    jsondiffpatch: require.resolve('jsondiffpatch'),
    htmlformatter: require.resolve('jsondiffpatch/formatters/html'),
  },
  output: {
    path: path.resolve(__dirname, 'vendor'),
    filename: '[name].umd.js',
    library: '[name]',
    libraryTarget: 'umd',
    globalObject: 'this',
  },
  mode: 'production',
};
