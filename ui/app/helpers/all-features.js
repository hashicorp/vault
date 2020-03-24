import { helper as buildHelper } from '@ember/component/helper';

const ALL_FEATURES = [
  'Control Groups',
  'DR Replication',
  'Entropy Augmentation',
  'HSM',
  'KMIP',
  'MFA',
  'Namespaces',
  'Performance Replication',
  'Performance Standby',
  'Sentinel',
  'Seal Wrapping',
  'Transform Secrets Engine',
];

export function allFeatures() {
  return ALL_FEATURES;
}

export default buildHelper(allFeatures);
