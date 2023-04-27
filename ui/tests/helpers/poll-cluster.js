/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { settled } from '@ember/test-helpers';

export async function pollCluster(owner) {
  const store = owner.lookup('service:store');
  await store.peekAll('cluster').firstObject.reload();
  await settled();
}
