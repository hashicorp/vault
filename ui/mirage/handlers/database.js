/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { Response } from 'miragejs';

export default function (server) {
  const getRecord = (schema, req, dbKey) => {
    const { backend, name } = req.params;
    const record = schema.db[dbKey].findBy({ name, backend });
    if (record) {
      delete record.backend;
      delete record.id;
    }
    return record ? { data: record } : new Response(404, {}, { errors: [] });
  };
  const createOrUpdateRecord = (schema, req, key) => {
    const { backend, name } = req.params;
    const payload = JSON.parse(req.requestBody);
    const record = schema[key].findOrCreateBy({ name, backend });
    record.update(payload);
    return new Response(204);
  };
  const deleteRecord = (schema, req, dbKey) => {
    const { name } = req.params;
    const record = schema.db[dbKey].findBy({ name });
    if (record) {
      schema.db[dbKey].remove(record.id);
    }
    return new Response(204);
  };

  // Connection mgmt
  server.get('/:backend/config/:name', (schema, req) => {
    return getRecord(schema, req, 'database/connections');
  });
  server.get('/:backend/config', (schema) => {
    const keys = schema.db['databaseConnections'].map((record) => record.name);
    if (!keys.length) {
      return new Response(404, {}, { errors: [] });
    }
    return {
      data: {
        keys,
      },
    };
  });
  server.post('/:backend/config/:name', (schema, req) => {
    const { name } = req.params;
    const { username } = JSON.parse(req.requestBody);
    if (name === 'bad-connection') {
      return new Response(
        500,
        {},
        {
          errors: [
            `error creating database object: error verifying - ping: Error 1045 (28000): Access denied for user '${username}'@'192.168.65.1' (using password: YES)`,
          ],
        }
      );
    }
    return createOrUpdateRecord(schema, req, 'database/connections');
  });
  server.delete('/:backend/config/:name', (schema, req) => {
    return deleteRecord(schema, req, 'database-connection');
  });
  // Rotate root
  server.post('/:backend/rotate-root/:name', (schema, req) => {
    const { name } = req.params;
    if (name === 'fail-rotate') {
      return new Response(
        500,
        {},
        {
          errors: [
            "1 error occurred:\n\t* failed to update user: failed to change password: Error 1045 (28000): Access denied for user 'admin'@'%' (using password: YES)\n\n",
          ],
        }
      );
    }
    return new Response(204);
  });

  // Generate credentials
  server.get('/:backend/creds/:role', (schema, req) => {
    const { role } = req.params;
    if (role === 'static-role') {
      // static creds
      return {
        request_id: 'static-1234',
        lease_id: '',
        renewable: false,
        lease_duration: 0,
        data: {
          last_vault_rotation: '2024-01-18T10:45:47.227193-06:00',
          password: 'generated-password',
          rotation_period: 86400,
          ttl: 3600,
          username: 'static-username',
        },
        wrap_info: null,
        warnings: null,
        auth: null,
        mount_type: 'database',
      };
    }
    // dynamic creds
    return {
      request_id: 'dynamic-1234',
      lease_id: `database/creds/${role}/abcd`,
      renewable: true,
      lease_duration: 3600,
      data: {
        password: 'generated-password',
        username: 'generated-username',
      },
      wrap_info: null,
      warnings: null,
      auth: null,
      mount_type: 'database',
    };
  });
}
