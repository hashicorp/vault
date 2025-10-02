/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import isDeleted from 'kv/helpers/is-deleted';

// helps in long logic statements for state of a currentVersion
export default function currentSecret(metadata) {
  if (metadata?.versions && metadata?.current_version) {
    const data = metadata.versions[metadata.current_version];
    const state = data.destroyed ? 'destroyed' : isDeleted(data.deletion_time) ? 'deleted' : 'created';
    return {
      state,
      isDeactivated: state !== 'created',
      deletionTime: data.deletion_time,
    };
  }
  return false;
}
