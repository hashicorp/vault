/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { Response } from 'miragejs';

export const associationsResponse = (schema, req) => {
  const { type, name } = req.params;
  const records = schema.db.syncAssociations.where({ type, name });

  return {
    data: {
      associated_secrets: records.length
        ? records.reduce((associations, association) => {
            const key = `${association.mount}/${association.secret_name}`;
            delete association.type;
            delete association.name;
            associations[key] = association;
            return associations;
          }, {})
        : {},
      store_name: name,
      store_type: type,
    },
  };
};

export const syncStatusResponse = (schema, req) => {
  const { mount, name: secret_name } = req.params;
  const records = schema.db.syncAssociations.where({ mount, secret_name });
  if (!records.length) {
    return new Response(404, {}, { errors: [] });
  }
  const STATUSES = ['SYNCED', 'SYNCING', 'UNSYNCED', 'UNSYNCING', 'INTERNAL_VAULT_ERROR', 'UNKNOWN'];
  const generatedRecords = records.reduce((records, record, index) => {
    const destinationType = record.type;
    const destinationName = record.name;
    record.sync_status = STATUSES[index];
    const key = `${destinationType}/${destinationName}`;
    records[key] = record;
    return records;
  }, {});
  if (records.length === 5) {
    // create one more record with sync_status = 'UNKNOWN' to mock each status option
    generatedRecords['aws-sm/my-aws-destination'] = {
      ...generatedRecords['aws-sm/destination-aws'],
      sync_status: 'UNKNOWN',
      name: 'my-aws-destination',
      updated_at: new Date().toISOString(),
    };
  }
  return {
    data: {
      associated_destinations: generatedRecords,
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
          const key = `${record.type}/`;
          if (!keyInfo[key]) {
            keyInfo[key] = [record.name];
          } else {
            keyInfo[key].push(record.name);
          }
          return keyInfo;
        }, {}),
        keys: records.map((r) => `${r.type}/`),
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
  server.get('/sys/sync/associations', (schema) => {
    const records = schema.db.syncAssociations.where({});
    if (!records.length) {
      return new Response(404, {}, { errors: [] });
    }
    // for now we only care about the total_associations value
    return {
      data: {
        key_info: {},
        keys: [],
        total_associations: records.length,
        total_secrets: records.reduce((secrets, association) => {
          const secretPath = `${association.mount}/${association.secret_name}`;
          if (!secrets.includes(secretPath)) {
            secrets.push(secretPath);
          }
          return secrets;
        }, []),
      },
    };
  });
  server.get(`${uri}/associations`, (schema, req) => {
    return associationsResponse(schema, req);
  });
  server.post(`${uri}/associations/set`, (schema, req) => {
    const { type, name } = req.params;
    const { secret_name, mount } = JSON.parse(req.requestBody);
    if (secret_name.slice(-1) === '/') {
      return new Response(
        400,
        {},
        { errors: ['Secret not found. Please provide full path to existing secret'] }
      );
    }
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
  server.get('sys/sync/associations/:mount/*name', (schema, req) => {
    return syncStatusResponse(schema, req);
  });
}
