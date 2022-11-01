/* eslint-disable node/no-extraneous-require */
const { buildEngine } = require('ember-engines/lib/engine-addon');

module.exports = buildEngine({
  name: 'pki',
  lazyLoading: {
    enabled: false,
  },
  isDevelopingAddon() {
    return true;
  },
});
