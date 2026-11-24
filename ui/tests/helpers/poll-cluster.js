/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { settled } from '@ember/test-helpers';

export async function pollCluster(owner) {
  await owner.lookup('service:current-cluster').cluster.reload();
  await settled();
}
