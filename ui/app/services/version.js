import { readOnly, match, not } from '@ember/object/computed';
import Service, { inject as service } from '@ember/service';
import { computed } from '@ember/object';
import { task } from 'ember-concurrency';

const hasFeatureMethod = (context, featureKey) => {
  const features = context.get('features');
  if (!features) {
    return false;
  }
  return features.includes(featureKey);
};
const hasFeature = featureKey => {
  return computed('features', 'features.[]', function() {
    return hasFeatureMethod(this, featureKey);
  });
};
export default Service.extend({
  _features: null,
  features: readOnly('_features'),
  version: null,
  store: service(),

  hasPerfReplication: hasFeature('Performance Replication'),

  hasDRReplication: hasFeature('DR Replication'),

  hasSentinel: hasFeature('Sentinel'),
  hasNamespaces: hasFeature('Namespaces'),

  isEnterprise: match('version', /\+.+$/),

  isOSS: not('isEnterprise'),

  setVersion(resp) {
    this.set('version', resp.version);
  },

  hasFeature(feature) {
    return hasFeatureMethod(this, feature);
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
    let response = yield this.get('store')
      .adapterFor('cluster')
      .health();
    this.setVersion(response);
    return;
  }),

  getFeatures: task(function*() {
    if (this.get('features.length') || this.get('isOSS')) {
      return;
    }
    try {
      let response = yield this.get('store')
        .adapterFor('cluster')
        .features();
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
