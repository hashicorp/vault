/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { Response } from 'miragejs';

export default function (server) {
  server.get('/sys/storage/raft/configuration', () => {
    return server.create('configuration', 'withRaft');
  });

  // LOAD SNAPSHOT FORM HANDLERS
  server.post('sys/storage/raft/snapshot-load', () => ({
    data: {},
  }));

  server.post('sys/storage/raft/snapshot-auto/snapshot-load/:config_name', () => ({
    data: {},
  }));

  // UNLOAD SNAPSHOT HANDLER
  server.delete('/sys/storage/raft/snapshot-load/:snapshot_id', (schema, req) => {
    const record = schema.db['snapshots'].findBy(req.params);
    if (record) {
      schema.db['snapshots'].remove(record);
      return new Response(204); // No content
    }
    return new Response(404, {}, { errors: [] });
  });

  // LOADED SNAPSHOT HANDLERS
  server.get('/sys/storage/raft/snapshot-load', (schema) => {
    // Currently only one snapshot can be loaded at a time
    const record = schema.db['snapshots'][0];

    if (record) {
      const { snapshot_id } = record;
      return { data: { keys: [snapshot_id] } };
    }
    return new Response(404, {}, { errors: [] });
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
          local: true,
          path: 'cubbyhole/',
        },
        'kv/': {
          type: 'kv',
          local: false,
          path: 'kv/',
        },
        'database/': {
          type: 'database',
          local: true,
          path: 'database/',
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

  server.get('/database/static-roles/:path', () => {
    return {
      data: {
        credential_type: 'password',
        db_name: 'test-db',
        rotation_period: 86400,
        rotation_statements: [],
        skip_import_rotation: true,
        username: 'super-user',
      },
    };
  });

  // RECOVER SECRET HANDLERS
  server.post('/cubbyhole/:path', () => ({
    data: {},
  }));

  server.post('/kv/:path', () => ({
    data: {},
  }));

  server.post('/database/static-roles/:path', () => ({
    data: {},
  }));
}
