/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { Response } from 'miragejs';

export default function (server) {
  const getRecord = (schema, req, dbKey) => {
    const { backend, name } = req.params;
    const record = schema.db[dbKey].findBy((r) => r.name === name && r.backend === backend);
    if (record) {
      delete record.backend;
      delete record.id;
    }
    return record ? { data: record } : new Response(404, {}, { errors: [] });
  };
  const createRecord = (req, key) => {
    const data = JSON.parse(req.requestBody);
    server.create(key, data);
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
    return getRecord(schema, req, 'databaseConnections');
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

    return createRecord(req, 'database-connection');
  });
  server.delete('/:backend/config/:name', (schema, req) => {
    return deleteRecord(schema, req, 'database-connection');
  });
  // Rotate root
  server.post('/:backend/rotate-root/:name', () => {
    new Response(204);
  });

  // Generate credentials
  server.post('/:path/creds/:role', (schema, req) => {
    const { role } = req.params;
    const record = schema.db.databaseRoles.findBy({ name: role });

    let errors;
    if (!record) {
      errors = [`role '${role}' does not exist`];
    }
    // creds cannot be fetched after creation so we don't need to store them
    return errors
      ? new Response(400, {}, { errors })
      : {
          request_id: 'iiiiiiii',
          lease_id: 'database/creds/dynamic-role/ijijijijjji',
          renewable: true,
          lease_duration: 3600,
          data: {
            password: 'some-generated-password',
            username: 'dynamic-username-abcdefg',
          },
          wrap_info: null,
          warnings: null,
          auth: null,
          mount_type: 'database',
        };
  });

  server.get('/sys/internal/ui/mounts/database', () => ({
    data: {
      accessor: 'database_9f846a87',
      path: 'database/',
      type: 'database',
    },
  }));
}

/* Connection failed due to bad verification:
POST v1/:backend/config/:name
{
  "backend": "database",
  "name": "awesome-db",
  "plugin_name": "mysql-database-plugin",
  "verify_connection": true,
  "connection_url": "{{username}}:{{password}}@tcp(127.0.0.1:33060)/",
  "username": "sudo2",
  "password": "my-oiwejfowijef",
  "max_open_connections": 4,
  "max_idle_connections": 0,
  "max_connection_lifetime": "0s",
  "root_rotation_statements": [
    "SELECT user from mysql.user",
    "GRANT ALL PRIVILEGES ON *.* to 'sudo'@'%'"
  ]
}
{ errors: [
  "error creating database object: error verifying - ping: Error 1045 (28000): Access denied for user 'sudo2'@'192.168.65.1' (using password: YES)"
]}

Connection succeeded: 204 no response
*/

/* Rotate root (no body) POST http://localhost:8200/v1/database/rotate-root/awesome-db
failed due to something:

{ errors: [
  "1 error occurred:\n\t* failed to update user: failed to change password: Error 1045 (28000): Access denied for user 'sudo2'@'%' (using password: YES)\n\n"
]}

Success: 204 no response
*/

/* Create role

First gets then updates the connection with allowed_roles update
Then POST http://localhost:8200/v1/database/roles/awesome-role with:
{
  "backend": "database",
  "name": "awesome-role",
  "type": "dynamic",
  "default_ttl": "1h",
  "max_ttl": "24h",
  "rotation_period": "24h",
  "path": "roles",
  "db_name": "awesome-db"
}

Success: 204 no response
*/

/* CREDS
dynamic GET http://localhost:8200/v1/database/creds/awesome-role
{
    "request_id": "d6248a7d-85db-c989-53d6-37a52ddf98cc",
    "lease_id": "database/creds/awesome-role/hbmpLDbXXJAH9Q23PdTYqIJX",
    "renewable": true,
    "lease_duration": 3600,
    "data": {
        "password": "abcd",
        "username": "v-token-awesome-ro-YYHIPH2BdpE5h"
    },
    "wrap_info": null,
    "warnings": null,
    "auth": null,
    "mount_type": "database"
}

export default function (server) {
  server.get('/database/roles', function () {
    return {
      data: { keys: ['my-role'] },
    };
  });
  server.get('/database/static-roles', function () {
    return {
      data: { keys: ['dev-static', 'prod-static'] },
    };
  });

  server.get('/database/static-roles/:rolename', function (db, req) {
    if (req.params.rolename.includes('no-exist')) {
      return new Response(400);
    }
    return {
      data: {
        rotation_statements: [
          '{ "db": "admin", "roles": [{ "role": "readWrite" }, {"role": "read", "db": "foo"}] }',
        ],
        db_name: 'connection',
        username: 'alice',
        rotation_period: '1h',
      },
    };
  });

  server.post('/database/rotate-role/:rolename', function () {
    return new Response(204);
  });

  server.get('/database/roles/my-role', function () {
    return {
      data: {
        creation_statements: [
          "CREATE ROLE \"{{name}}\" WITH LOGIN PASSWORD '{{password}}' VALID UNTIL '{{expiration}}';",
          'GRANT SELECT ON ALL TABLES IN SCHEMA public TO "{{name}}";',
        ],
        credential_type: 'password',
        db_name: 'mysql',
        default_ttl: 3600,
        max_ttl: 86400,
        renew_statements: [],
        revocation_statements: [],
        rollback_statements: [],
      },
    };
  });
}

/*
Failure on create due to bad password:
{
    "errors": [
        "error creating database object: error verifying - ping: Error 1045 (28000): Access denied for user 'root'@'192.168.65.1' (using password: YES)"
    ]
}

Success on create & rotate:
204 No Content
*/
