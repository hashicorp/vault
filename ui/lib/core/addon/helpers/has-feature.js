/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/* eslint-disable ember/no-observers */
import { service } from '@ember/service';
import { assert } from '@ember/debug';
import Helper from '@ember/component/helper';
import { observer } from '@ember/object';

const POSSIBLE_FEATURES = [
  'HSM',
  'Performance Replication',
  'DR Replication',
  'MFA',
  'Sentinel',
  'Seal Wrapping',
  'Control Groups',
  'Namespaces',
  'KMIP',
  'Transform Secrets Engine',
  'Key Management Secrets Engine',
];

export function hasFeature(featureName, features) {
  if (!POSSIBLE_FEATURES.includes(featureName)) {
    assert(`${featureName} is not one of the available values for Vault Enterprise features.`, false);
    return false;
  }
  return features ? features.includes(featureName) : false;
}

export default Helper.extend({
  version: service(),
  onFeaturesChange: observer('version.features.[]', function () {
    this.recompute();
  }),
  compute([featureName]) {
    return hasFeature(featureName, this.version.features);
  },
});
