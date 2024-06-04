/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { settled } from '@ember/test-helpers';

export async function pollCluster(owner) {
  const store = owner.lookup('service:store');
  await store.peekAll('cluster')[0].reload();
  await settled();
}
