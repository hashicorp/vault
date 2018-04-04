import Ember from 'ember';
const { Helper, inject, observer } = Ember;

const FEATURES = [
  'HSM',
  'Performance Replication',
  'DR Replication',
  'MFA',
  'Sentinel',
  'AWS KMS Autounseal',
  'GCP CKMS Autounseal',
  'Seal Wrapping',
  'Control Groups',
];

export function hasFeature(featureName, features) {
  if (!FEATURES.includes(featureName)) {
    Ember.assert(`${featureName} is not one of the available values for Vault Enterprise features.`, false);
    return false;
  }
  return features ? features.includes(featureName) : false;
}

export default Helper.extend({
  version: inject.service(),
  onFeaturesChange: observer('version.features.[]', function() {
    this.recompute();
  }),
  compute([featureName]) {
    return hasFeature(featureName, this.get('version.features'));
  },
});
