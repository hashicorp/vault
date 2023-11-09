#!/usr/bin/env node
/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/* eslint-disable */
// run this script via yarn in the ui directory:
// yarn docfy-md some-component
//
// or if the story is for a component in an in-repo-addon or an engine:
// yarn docfy-md some-component name-of-engine

const fs = require('fs');
const jsdoc2md = require('jsdoc-to-markdown');
const [name, addonOrEngine] = process.argv.slice(2);

const inputFile = addonOrEngine
  ? `lib/${addonOrEngine}/addon/components/${name}.js`
  : `app/components/${name}.js`;
const outputFile = `docs/components/${name}.md`;

const outputFormat = `
{{#module}}
# {{name}}
{{>body}}
{{/module}}
`;

const options = {
  files: inputFile,
  'example-lang': 'hbs preview-template',
  template: outputFormat,
};
const md = jsdoc2md.renderSync(options);
fs.writeFileSync(outputFile, md);
