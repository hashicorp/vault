import Ember from 'ember';
const { Helper, inject, observer } = Ember;

export default Helper.extend({
  version: inject.service(),
  onFeaturesChange: observer('version.version', function() {
    this.recompute();
  }),
  compute([sku]) {
    if (sku !== 'OSS' && sku !== 'Enterprise') {
      Ember.assert(`${sku} is not one of the available values for Vault versions.`, false);
      return false;
    }
    return this.get(`version.is${sku}`);
  },
});
