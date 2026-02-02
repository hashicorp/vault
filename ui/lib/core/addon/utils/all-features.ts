/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

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
  'KMIP',
  'Entropy Augmentation',
  'Transform Secrets Engine',
  'Secrets Sync',
  'PKI-only Secrets',
];

export function allFeatures() {
  return ALL_FEATURES;
}
