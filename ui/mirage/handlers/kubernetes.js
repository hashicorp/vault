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
    const record = getRecord(schema, req, dbKey);
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
  server.get('/:path/roles', (schema) => {
    return {
      data: {
        keys: schema.db.kubernetesRoles.where({}).mapBy('name'),
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
    return errors ? new Response(400, {}, { errors }) : new Response(204);
  });
}
