/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import isDeleted from 'kv/helpers/is-deleted';

export default function sortedVersions(versions) {
  const array = [];
  for (const key in versions) {
    const version = versions[key];
    const isSecretDeleted = isDeleted(version.deletion_time);
    array.push({ version: key, isSecretDeleted, ...version });
  }
  // version keys are in order created with 1 being the oldest, we want newest first
  return array.reverse();
}
