/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export default function (server) {
  server.create('snapshot');
  // Other snapshot states, only one snapshot should be "loaded" at a time
  // server.create('snapshot', { status: 'loading' });
  // server.create('snapshot', { status: 'error' });
}
