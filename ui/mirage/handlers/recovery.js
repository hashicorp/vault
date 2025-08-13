/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export default function (server) {
  server.get('/sys/storage/raft/configuration', () => {
    return server.create('configuration', 'withRaft');
  });

  server.get('/sys/storage/raft/snapshot-load', (schema) => {
    // Currently only one snapshot can be loaded at a time
    const { snapshot_id } = schema.db['snapshots'][0];
    return { data: { keys: [snapshot_id] } };
  });

  server.get('/sys/storage/raft/snapshot-load/:snapshot_id', (schema, req) => {
    const record = schema.db['snapshots'].findBy(req.params);
    if (record) {
      delete record.id; // "snapshot_id" is the id
      return { data: record };
    }
    return new Response(404, {}, { errors: [] });
  });
}
