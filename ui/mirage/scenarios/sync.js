/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export default function (server) {
  const types = ['aws-sm', 'azure-kv', 'gcp-sm', 'gh', 'vercel-project'];
  types.forEach((type) => {
    server.create('sync-destination', type);
    server.create('sync-association', { type, name: `destination-${type.split('-')[0]}` });
  });
}
