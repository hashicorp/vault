/* eslint-env node */
/* eslint-disable ember/avoid-leaking-state-in-ember-objects */
'use strict';

const EngineAddon = require('ember-engines/lib/engine-addon');

module.exports = EngineAddon.extend({
  name: 'open-api-explorer',

  included() {
    this._super.included && this._super.included.apply(this, arguments);
    // we want to lazy load these deps, importing them here will result in them being added to the
    // engine-vendor files that will be lazy loaded with the engine
    this.import('node_modules/swagger-ui-dist/swagger-ui-bundle.js');
    this.import('node_modules/swagger-ui-dist/swagger-ui.css');
  },

  lazyLoading: {
    enabled: true,
  },

  isDevelopingAddon() {
    return true;
  },
});
