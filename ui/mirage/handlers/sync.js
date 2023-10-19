/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { Response } from 'miragejs';

export default function (server) {
  const base = '/sys/sync/destinations';
  const uri = `${base}/:type/:name`;
  // destinations
  server.get(base, (schema) => {
    const records = schema.db.syncDestinations.where({});
    if (!records.length) {
      return new Response(404, {}, { errors: [] });
    }
    return {
      data: {
        key_info: records.reduce((keyInfo, record) => {
          if (!keyInfo[record.type]) {
            keyInfo[record.type] = [record.name];
          } else {
            keyInfo[record.type].push(record.name);
          }
          return keyInfo;
        }, {}),
        keys: records.map((r) => r.type),
      },
    };
  });
  server.get(uri, (schema, req) => {
    const { type, name } = req.params;
    const record = schema.db.syncDestinations.findBy({ type, name });
    if (record) {
      delete record.type;
      delete record.name;
      delete record.id;
      return {
        data: {
          connection_details: record,
          name,
          type,
        },
      };
    }
    return new Response(404, {}, { errors: [] });
  });
  server.post(uri, (schema, req) => {
    const { type, name } = req.params;
    const data = { ...JSON.parse(req.requestBody), type, name };
    schema.db.syncDestinations.firstOrCreate({ type, name }, data);
    schema.db.syncDestinations.update({ type, name }, data);
    return new Response(204);
  });
  server.delete(uri, (schema, req) => {
    const { type, name } = req.params;
    schema.db.syncDestinations.remove({ type, name });
    return new Response(204);
  });
  // associations
  server.get(`${uri}/associations`, (schema, req) => {
    const { type, name } = req.params;
    const records = schema.db.syncAssociations.where({ type, name });
    if (!records.length) {
      return new Response(404, {}, { errors: [] });
    }
    return {
      data: {
        associated_secrets: records.reduce((associations, association) => {
          const key = `${association.accessor}/${association.secret_name}`;
          delete association.type;
          delete association.name;
          associations[key] = association;
          return associations;
        }, {}),
        store_name: name,
        store_type: type,
      },
    };
  });
  server.post(`${uri}/associations/set`, (schema, req) => {
    const { type, name } = req.params;
    const { secret_name, mount: accessor } = JSON.parse(req.requestBody);
    const data = { type, name, accessor, secret_name };
    schema.db.syncAssociations.firstOrCreate({ type, name }, data);
    schema.db.syncAssociations.update({ type, name }, data);
    return new Response(204);
  });
  server.post(`${uri}/associations/remove`, (schema, req) => {
    const { type, name } = req.params;
    schema.db.syncAssociations.remove({ type, name });
    return new Response(204);
  });
}
