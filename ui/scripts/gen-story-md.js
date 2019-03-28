#!/usr/bin/env node
/* eslint-disable */
const fs = require('fs');
const jsdoc2md = require('jsdoc-to-markdown');
var args = process.argv.slice(2);
const name = args[0];
const inputFile = `app/components/${name}.js`;
const outputFile = `stories/${name}.md`;
const component = name
  .split('-')
  .map(word => word.charAt(0).toUpperCase() + word.slice(1))
  .join('');
const options = {
  files: inputFile,
  template: fs.readFileSync('./lib/story-md.hbs', 'utf8'),
};
let md = jsdoc2md.renderSync(options);

const pageBreakIndex = md.lastIndexOf('---'); //this is our last page break

const seeLinks = `**See**

- [Uses of ${component}](https://github.com/hashicorp/vault/search?l=Handlebars&q=${component})
- [${component} Source Code](https://github.com/hashicorp/vault/blob/master/ui/app/components/${name}.js) 

`;

md = md.slice(0, pageBreakIndex) + seeLinks + md.slice(pageBreakIndex);

fs.writeFileSync(outputFile, md);
