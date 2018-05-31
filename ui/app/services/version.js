import Ember from 'ember';
import { task } from 'ember-concurrency';

const { Service, inject, computed } = Ember;

const hasFeature = featureKey => {
  return computed('features', 'features.[]', function() {
    const features = this.get('features');
    if (!features) {
      return false;
    }
    return features.includes(featureKey);
  });
};
export default Service.extend({
  _features: null,
  features: computed.readOnly('_features'),
  version: null,
  store: inject.service(),

  hasPerfReplication: hasFeature('Performance Replication'),

  hasDRReplication: hasFeature('DR Replication'),

  hasSentinel: hasFeature('Sentinel'),

  isEnterprise: computed.match('version', /\+.+$/),

  isOSS: computed.not('isEnterprise'),

  setVersion(resp) {
    this.set('version', resp.version);
  },

  setFeatures(resp) {
    if (!resp.features) {
      return;
    }
    this.set('_features', resp.features);
  },

  getVersion: task(function*() {
    if (this.get('version')) {
      return;
    }
    let response = yield this.get('store').adapterFor('cluster').health();
    this.setVersion(response);
    return;
  }),

  getFeatures: task(function*() {
    if (this.get('features.length') || this.get('isOSS')) {
      return;
    }
    try {
      let response = yield this.get('store').adapterFor('cluster').features();
      this.setFeatures(response);
      return;
    } catch (err) {
      // if we fail here, we're likely in DR Secondary mode and don't need to worry about it
    }
  }).keepLatest(),

  fetchVersion: function() {
    return this.get('getVersion').perform();
  },
  fetchFeatures: function() {
    return this.get('getFeatures').perform();
  },
});
