'use strict';
const getPathOption = require('ember-cli-get-component-path-option');
const stringUtil = require('ember-cli-string-utils');
const path = require('path');

function findAddonByName(addonOrProject, name) {
  let addon = addonOrProject.addons.find(addon => addon.name === name);

  if (addon) {
    return addon;
  }

  return addonOrProject.addons.find(addon => findAddonByName(addon, name));
}

module.exports = {
  description: 'generates a story for storybook',

  fileMapTokens: function() {
    let { project } = this;
    return {
      __path__: function(options) {
        if (options.inRepoAddon) {
          let addon = findAddonByName(project, options.inRepoAddon);
          return path.relative(project.root, addon.root);
        }
        return path.relative(project.root, project.root);
      },
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

    let importMD = "import notes from './" + stringUtil.dasherize(options.entity.name) + ".md';\n";
    return {
      importMD: importMD,
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
