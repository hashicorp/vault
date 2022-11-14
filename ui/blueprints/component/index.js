'use strict';
/* eslint-disable node/no-extraneous-require */
/* eslint-disable ember/no-string-prototype-extensions */
const path = require('path');
const stringUtil = require('ember-cli-string-utils');
const pathUtil = require('ember-cli-path-utils');
const getPathOption = require('ember-cli-get-component-path-option');
const normalizeEntityName = require('ember-cli-normalize-entity-name');

module.exports = {
  description: 'Generates a component.',

  availableOptions: [
    {
      name: 'path',
      type: String,
      default: 'components',
      aliases: [{ 'no-path': '' }],
    },
  ],

  filesPath: function () {
    const filesDirectory = 'files';

    return path.join(this.path, filesDirectory);
  },

  fileMapTokens: function () {
    return {
      __path__: function (options) {
        if (options.pod) {
          return path.join(options.podPath, options.locals.path, options.dasherizedModuleName);
        } else {
          return 'components';
        }
      },
      __templatepath__: function (options) {
        if (options.pod) {
          return path.join(options.podPath, options.locals.path, options.dasherizedModuleName);
        }
        return 'templates/components';
      },
      __templatename__: function (options) {
        if (options.pod) {
          return 'template';
        }
        return options.dasherizedModuleName;
      },
    };
  },

  normalizeEntityName: function (entityName) {
    return normalizeEntityName(entityName);
  },

  locals: function (options) {
    let exportDefault = 'export default ';
    let exportAddOn = '';
    let importTemplate = '';
    let setComponentTemplate = '';
    let templatePath = '';

    // if we're in an addon, build import statement and set layout
    if (options.project.isEmberCLIAddon() || (options.inRepoAddon && !options.inDummy) || !!options.in) {
      if (options.pod) {
        templatePath = './template';
      } else {
        templatePath =
          pathUtil.getRelativeParentPath(options.entity.name) +
          'templates/components/' +
          stringUtil.dasherize(options.entity.name);
      }
      exportDefault = '';
      exportAddOn = `export default setComponentTemplate(layout, ${stringUtil.classify(
        options.entity.name
      )})`;
      importTemplate = "import layout from '" + templatePath + "';";
      setComponentTemplate = "import { setComponentTemplate } from '@ember/component'; \n";
    }

    return {
      exportDefault: exportDefault,
      exportAddOn: exportAddOn,
      importTemplate: importTemplate,
      setComponentTemplate: setComponentTemplate,
      path: getPathOption(options),
    };
  },
};
