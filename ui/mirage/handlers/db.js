/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export default function (server) {
  server.get('/database/static-roles', function () {
    return {
      data: { keys: ['dev-static', 'prod-static'] },
    };
  });

  server.get('/database/static-roles/:rolename', function (db, req) {
    if (req.params.rolename.includes('tester')) {
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
}
