import config from '../config/environment';

export default {
  name: 'ember-inspect-disable',
  initialize: function() {
    if (config.environment === 'production') {
      // disables ember inspector
      window.NO_EMBER_DEBUG = true;
    }
  },
};
