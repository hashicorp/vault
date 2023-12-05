/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export default function (server) {
  const types = ['aws-sm', 'azure-kv', 'gcp-sm', 'gh', 'vercel-project'];
  types.forEach((type) => {
    server.create('sync-destination', type);
    const destinationDetails = { type, name: `destination-${type.split('-')[0]}` };
    server.create('sync-association', destinationDetails);
    if (['azure-kv', 'vercel-project'].includes(type)) {
      server.create('sync-association', {
        ...destinationDetails,
        secret_name: 'my-path/nested-path/nested-secret-1',
        sync_status: 'UNSYNCED',
        updated_at: '2023-11-05T16:15:25.961861096-04:00',
      });
    }
  });
  // create destination with no associations
  server.create('sync-destination', 'aws-sm', { name: 'new-destination' });
}
