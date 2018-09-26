#!/usr/bin/env node
/* eslint-disable */

// We need an array in this format for all of the files
//https://github.com/ember-template-lint/ember-cli-template-lint/blob/1bc03444ecf367473108cb28208cb3123199f950/.template-lintrc.js#L9

var walkSync = require('walk-sync');
var templates = walkSync('app', { globs: ['**/*.hbs'] });

templates = templates.map(path => {
  // we want the relative path w/o the extension:
  // 'app/templates/path/to/file/filename'
  return `app/${path.replace(/\.hbs$/, '')}`;
});

// stringify because if we don't console won't output the full list lol
console.log(JSON.stringify(templates, null, 2));
