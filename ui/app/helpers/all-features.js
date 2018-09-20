import Ember from 'ember';

const ALL_FEATURES = [
  'UI',
  'HSM',
  'Performance Replication',
  'DR Replication',
  'MFA',
  'Sentinel',
  'AWS KMS Autounseal',
  'GCP CKMS Autounseal',
  'Seal Wrapping',
  'Control Groups',
  'Azure Key Vault Seal',
  'Performance Standby',
  'Namespaces',
];

export function allFeatures() {
  return ALL_FEATURES;
}

export default Ember.Helper.helper(allFeatures);
