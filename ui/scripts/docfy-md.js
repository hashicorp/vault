#!/usr/bin/env node
/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/* eslint-disable */
// run this script the ui directory:
// yarn docfy-md some-component
//
// or if the docs are for a component in an in-repo-addon or an engine:
// yarn docfy-md some-component name-of-engine

const fs = require('fs');
const jsdoc2md = require('jsdoc-to-markdown');
const [nameOrFile, addonOrEngine] = process.argv.slice(2);

const name = nameOrFile.includes('.') ? nameOrFile?.split('.')[0] : nameOrFile; // can pass component-name or component-name.js
const path = nameOrFile.includes('.') ? nameOrFile : `${nameOrFile}.js`; // default to js

const inputFile = addonOrEngine ? `lib/${addonOrEngine}/addon/components/${path}` : `app/components/${path}`;
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

try {
  const md = jsdoc2md.renderSync(options);
  fs.writeFileSync(outputFile, md);
} catch (error) {
  console.log(name, error);
}
