/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { Response } from 'miragejs';

export default function (server) {
  const getRecord = (schema, req, dbKey) => {
    const { path, name } = req.params;
    const findBy = dbKey === 'kubernetesConfigs' ? { path } : { name };
    const record = schema.db[dbKey].findBy(findBy);
    if (record) {
      delete record.path;
      delete record.id;
    }
    return record ? { data: record } : new Response(404, {}, { errors: [] });
  };
  const createRecord = (req, key) => {
    const data = JSON.parse(req.requestBody);
    if (key === 'kubernetes-config') {
      data.path = req.params.path;
    }
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

  server.get('/:path/config', (schema, req) => {
    return getRecord(schema, req, 'kubernetesConfigs');
  });
  server.post('/:path/config', (schema, req) => {
    return createRecord(req, 'kubernetes-config');
  });
  server.delete('/:path/config', (schema, req) => {
    return deleteRecord(schema, req, 'kubernetesConfigs');
  });
  // endpoint for checking for environment variables necessary for inferred config
  server.get('/:path/check', () => {
    const response = {};
    const status = Math.random() > 0.5 ? 204 : 404;
    if (status === 404) {
      response.errors = [
        'Missing environment variables: KUBERNETES_SERVICE_HOST, KUBERNETES_SERVICE_PORT_HTTPS',
      ];
    }
    return new Response(status, response);
  });
  server.get('/:path/roles', (schema) => {
    return {
      data: {
        keys: schema.db.kubernetesRoles.where({}).map((role) => role.name),
      },
    };
  });
  server.get('/:path/roles/:name', (schema, req) => {
    return getRecord(schema, req, 'kubernetesRoles');
  });
  server.post('/:path/roles/:name', (schema, req) => {
    return createRecord(req, 'kubernetes-role');
  });
  server.delete('/:path/roles/:name', (schema, req) => {
    return deleteRecord(schema, req, 'kubernetesRoles');
  });
  server.post('/:path/creds/:role', (schema, req) => {
    const { role } = req.params;
    const record = schema.db.kubernetesRoles.findBy({ name: role });
    const data = JSON.parse(req.requestBody);
    let errors;
    if (!record) {
      errors = [`role '${role}' does not exist`];
    } else if (!data.kubernetes_namespace) {
      errors = ["'kubernetes_namespace' is required"];
    }
    // creds cannot be fetched after creation so we don't need to store them
    return errors
      ? new Response(400, {}, { errors })
      : {
          request_id: '58fefc6c-5195-c17a-94f2-8f889f3df57c',
          lease_id: 'kubernetes/creds/default-role/aWczfcfJ7NKUdiirJrPXIs38',
          renewable: false,
          lease_duration: 3600,
          data: {
            service_account_name: 'default',
            service_account_namespace: 'default',
            service_account_token: 'eyJhbGciOiJSUzI1NiIsImtpZCI6Imlr',
          },
        };
  });

  server.get('/sys/internal/ui/mounts/kubernetes', () => ({
    data: {
      accessor: 'kubernetes_9f846a87',
      path: 'kubernetes/',
      type: 'kubernetes',
    },
  }));
}
