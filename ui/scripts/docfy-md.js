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

see the readme ui/docs/how-to-docfy.md for more info
*/

const fs = require('fs');
const jsdoc2md = require('jsdoc-to-markdown');

// the fullFilepath arg will override the assumed inputFile which is built using the addonOrEngine arg
const [nameOrFile, addonOrEngine, fullFilepath] = process.argv.slice(2);

const isFile = nameOrFile.includes('.');
const name = isFile ? nameOrFile?.split('.')[0] : nameOrFile; // can pass component-name or component-name.js
const path = isFile ? nameOrFile : `${nameOrFile}.js`; // default to js

const inputFile = addonOrEngine ? `lib/${addonOrEngine}/addon/components/${path}` : `app/components/${path}`;
const outputFile = `docs/components/${name}.md`;

const outputFormat = `
{{#module}}
# {{name}}
{{>body}}
{{/module}}
`;

const options = {
  files: fullFilepath || inputFile,
  'example-lang': 'hbs preview-template',
  configure: './jsdoc2md.json',
  template: outputFormat,
};

try {
  const md = jsdoc2md.renderSync(options);
  // for some reason components without a jsdoc @module doesn't throw, so throw manually
  if (md.includes('ERROR')) throw `${md} (there is probably no jsdoc for this component)`;
  fs.writeFileSync(outputFile, md);
  console.log(`✅ ${name}`);
} catch (error) {
  console.log(`❌ ${name}`);
  console.log(error);
}
