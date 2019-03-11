'use strict';
const getPathOption = require('ember-cli-get-component-path-option');
const stringUtil = require('ember-cli-string-utils');

module.exports = {
  description: 'generates a story for storybook',

  fileMapTokens: function() {
    return {
      __markdownname__: function(options) {
        return options.dasherizedModuleName;
      },
      __name__: function(options) {
        return options.dasherizedModuleName;
      },
    };
  },

  locals: function(options) {
    let contents = '';
    return {
      contents: contents,
      path: getPathOption(options),
      header: stringUtil
        .dasherize(options.entity.name)
        .split('-')
        .map(word => stringUtil.capitalize(word))
        .join(' '),
    };
  },
};
