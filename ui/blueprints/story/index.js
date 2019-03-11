'use strict';
const getPathOption = require('ember-cli-get-component-path-option');

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
    };
  },
};
