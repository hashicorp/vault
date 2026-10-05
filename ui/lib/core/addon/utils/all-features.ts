/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

// Used to display the active vs inactive features table on the license page (/ui/vault/license)
// and in test helpers to stub full feature sets.
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
  'Agentic IAM',
];

export function allFeatures() {
  return ALL_FEATURES;
}
