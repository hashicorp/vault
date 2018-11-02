import { helper as buildHelper } from '@ember/component/helper';

const ALL_FEATURES = [
  'HSM',
  'Performance Replication',
  'DR Replication',
  'MFA',
  'Sentinel',
  'Seal Wrapping',
  'Control Groups',
  'Performance Standby',
  'Namespaces',
];

export function allFeatures() {
  return ALL_FEATURES;
}

export default buildHelper(allFeatures);
