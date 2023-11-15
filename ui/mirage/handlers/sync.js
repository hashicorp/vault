/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { Response } from 'miragejs';

export const associationsResponse = (schema, req) => {
  const { type, name } = req.params;
  const records = schema.db.syncAssociations.where({ type, name });

  if (!records.length) {
    return new Response(404, {}, { errors: [] });
  }
  return {
    data: {
      associated_secrets: records.reduce((associations, association) => {
        const key = `${association.mount}/${association.secret_name}`;
        delete association.type;
        delete association.name;
        associations[key] = association;
        return associations;
      }, {}),
      store_name: name,
      store_type: type,
    },
  };
};

export default function (server) {
  const base = '/sys/sync/destinations';
  const uri = `${base}/:type/:name`;

  const destinationResponse = (record) => {
    delete record.id;
    const { name, type, ...connection_details } = record;
    return {
      data: {
        connection_details,
        name,
        type,
      },
    };
  };

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
      return destinationResponse(record);
    }
    return new Response(404, {}, { errors: [] });
  });
  server.post(uri, (schema, req) => {
    const { type, name } = req.params;
    const data = { ...JSON.parse(req.requestBody), type, name };
    schema.db.syncDestinations.firstOrCreate({ type, name }, data);
    const record = schema.db.syncDestinations.update({ type, name }, data);
    return destinationResponse(record);
  });
  server.delete(uri, (schema, req) => {
    const { type, name } = req.params;
    schema.db.syncDestinations.remove({ type, name });
    return new Response(204);
  });
  // associations
  server.get(`${uri}/associations`, (schema, req) => {
    return associationsResponse(schema, req);
  });
  server.post(`${uri}/associations/set`, (schema, req) => {
    const { type, name } = req.params;
    const { secret_name, mount } = JSON.parse(req.requestBody);
    const data = { type, name, mount, secret_name };
    schema.db.syncAssociations.firstOrCreate({ type, name }, data);
    schema.db.syncAssociations.update({ type, name }, { ...data, sync_status: 'SYNCED' });
    return associationsResponse(schema, req);
  });
  server.post(`${uri}/associations/remove`, (schema, req) => {
    const { type, name } = req.params;
    schema.db.syncAssociations.update({ type, name }, { sync_status: 'UNSYNCED' });
    return associationsResponse(schema, req);
  });
}
