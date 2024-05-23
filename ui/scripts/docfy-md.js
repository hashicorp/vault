#!/usr/bin/env node

/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/* eslint-disable */

/*
run from the ui directory:
yarn docfy-md some-component

or if the docs are for a component in an in-repo-addon or an engine:
yarn docfy-md some-component name-of-engine

this script is currently limited to only .js files
for some reason when the input file is a typescript file
the following error is returned "JSDOC_ERROR: There are no input files to process."
*/

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
  if (md.includes('ERROR')) throw `${md} (there is probably no jsdoc for this component)`;
  fs.writeFileSync(outputFile, md);
  console.log(`✅ ${name}`);
} catch (error) {
  console.log(`❌ ${name}`);
  console.log(error);
}
