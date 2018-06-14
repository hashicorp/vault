/* eslint-env node */
'use strict';

var path = require('path');
var Funnel = require('broccoli-funnel');
var mergeTrees = require('broccoli-merge-trees');

module.exports = {
  name: 'bulma',

  isDevelopingAddon() {
    return true;
  },

  included: function(app) {
    this._super.included.apply(this, arguments);

    // see: https://github.com/ember-cli/ember-cli/issues/3718
    while (typeof app.import !== 'function' && app.app) {
      app = app.app;
    }

    this.bulmaPath = path.dirname(require.resolve('bulma'));
    this.bulmaSwitchPath = path.dirname(require.resolve('bulma-switch/switch.sass'));
    this.bulmaCheckPath = path.dirname(require.resolve('cool-checkboxes-for-bulma.io'));
    this.sassSVGURIPath = path.dirname(require.resolve('sass-svg-uri'));
    return app;
  },

  treeForStyles: function() {
    var bulma = new Funnel(this.bulmaPath, {
      srcDir: '/',
      destDir: 'app/styles/bulma',
      annotation: 'Funnel (bulma)',
    });

    var bulmaSwitch = new Funnel(this.bulmaSwitchPath, {
      srcDir: '/',
      destDir: 'app/styles/bulma',
      annotation: 'Funnel (bulma-switch)',
    });
    var bulmaCheck = new Funnel(this.bulmaCheckPath, {
      srcDir: '/',
      destDir: 'app/styles/bulma',
      annotation: 'Funnel (bulma-check)',
    });
    var sassSVGURI = new Funnel(this.sassSVGURIPath, {
      srcDir: '/',
      destDir: 'app/styles/sass-svg-uri',
      annotation: 'Sass SVG URI',
    });

    return mergeTrees([bulmaCheck, bulmaSwitch, bulma, sassSVGURI], { overwrite: true });
  },
};
