/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { Response } from 'miragejs';

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

  // NAMESPACE SEARCH SELECT
  server.get('/sys/internal/ui/namespaces', () => ({
    data: {
      keys: ['child-ns-1', 'child-ns-1/nested', 'child-ns-2'],
    },
  }));

  // MOUNT SEARCH SELECT
  server.get('/sys/internal/ui/mounts', () => ({
    data: {
      secret: {
        'cubbyhole/': {
          type: 'cubbyhole',
          description: 'per-token private secret storage',
          local: true,
          seal_wrap: false,
          external_entropy_access: false,
          config: {
            default_lease_ttl: 0,
            max_lease_ttl: 0,
            force_no_cache: false,
          },
        },
        'kv/': {
          type: 'kv',
          description: 'key/value secret storage',
          options: { version: '1' },
          local: false,
          seal_wrap: false,
          external_entropy_access: false,
          config: {
            default_lease_ttl: 0,
            max_lease_ttl: 0,
            force_no_cache: false,
          },
        },
      },
    },
  }));

  // READ SECRET HANDLERS
  server.get('/cubbyhole/:path', (schema, req) => {
    const path = req.params.path;

    // Mock data for different paths
    if (path === 'my-path') {
      return {
        data: {
          secret_key: 'secret_value',
          another_key: 'another_value',
        },
      };
    }

    if (path === 'nonexistent-secret') {
      return new Response(404, {}, { errors: ['An error occurred, please try again'] });
    }

    // Default mock data for any other path
    return {
      data: {
        key1: 'value1',
        key2: 'value2',
      },
    };
  });

  server.get('/kv/:path', (schema, req) => {
    const path = req.params.path;

    // Mock data for different paths
    if (path === 'my-path') {
      return {
        data: {
          username: 'admin',
          password: 'secret123',
        },
      };
    }

    // Default mock data for any other path
    return {
      data: {
        foo: 'bar',
        baz: 'qux',
      },
    };
  });

  // RECOVER SECRET HANDLERS
  server.post('/cubbyhole/:path', () => ({
    data: {},
  }));

  // server.post('/kv/:path', () => ({
  //   data: {},
  // }));
}
