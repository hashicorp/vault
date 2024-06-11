/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { Response } from 'miragejs';
import { camelize } from '@ember/string';
import { findDestination } from 'core/helpers/sync-destinations';
import clientsHandler from './clients';
import modifyPassthroughResponse from '../helpers/modify-passthrough-response';

export const associationsResponse = (schema, req) => {
  const { type, name } = req.params;
  const [destination] = schema.db.syncDestinations.where({ type, name });
  const records = schema.db.syncAssociations.where({ type, name });
  const associations = records.length
    ? records.reduce((associations, association) => {
        const key = `${association.mount}_12345/${association.secret_name}`;
        delete association.type;
        delete association.name;
        associations[key] = association;
        return associations;
      }, {})
    : {};

  // if a destination has granularity: 'secret-key' keys of the secret
  // are added to the association response but they are not individual associations
  // the secret itself is still a single association
  const subKeys = {
    'my-kv_12345/my-granular-secret/foo': {
      mount: 'my-kv',
      secret_name: 'my-granular-secret',
      sync_status: 'SYNCED',
      updated_at: '2023-09-20T10:51:53.961861096-04:00',
      sub_key: 'foo',
    },
    'my-kv_12345/my-granular-secret/bar': {
      mount: 'my-kv',
      secret_name: 'my-granular-secret',
      sync_status: 'SYNCED',
      updated_at: '2023-09-20T10:51:53.961861096-04:00',
      sub_key: 'bar',
    },
    'my-kv_12345/my-granular-secret/baz': {
      mount: 'my-kv',
      secret_name: 'my-granular-secret',
      sync_status: 'SYNCED',
      updated_at: '2023-09-20T10:51:53.961861096-04:00',
      sub_key: 'baz',
    },
  };

  return {
    data: {
      associated_secrets: destination.granularity === 'secret-path' ? associations : subKeys,
      store_name: name,
      store_type: type,
    },
  };
};

export const syncStatusResponse = (schema, req) => {
  const { mount, secret_name } = req.queryParams;
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

const createOrUpdateDestination = (schema, req) => {
  const { type, name } = req.params;
  const request = JSON.parse(req.requestBody);
  const apiResponse = {};
  for (const attr in request) {
    // API returns ***** for credentials sent in a request
    // and returns nothing if empty (assume using environment variables)
    const { maskedParams } = findDestination(type);
    if (maskedParams.includes(camelize(attr))) {
      apiResponse[attr] = request[attr] === '' ? '' : '*****';
    } else {
      apiResponse[attr] = request[attr];
    }
  }
  const data = { ...apiResponse, type, name };
  // issue with mirages' update method not returning an id on the payload which causes ember data to error after 4.12.x upgrade.
  // to work around this, determine if we're creating or updating a record first
  const records = schema.db.syncDestinations.where({ type, name });

  if (!records.length) {
    return schema.db.syncDestinations.firstOrCreate({ type, name }, data);
  } else {
    return schema.db.syncDestinations.update({ type, name }, data);
  }
};

export default function (server) {
  // default to enterprise with Secrets Sync on the license and activated
  server.get('sys/health', (schema, req) => modifyPassthroughResponse(req, { enterprise: true }));
  server.get('/sys/license/features', () => ({ features: ['Secrets Sync'] }));
  server.get('/sys/activation-flags', () => {
    return {
      data: {
        activated: ['secrets-sync'],
        unactivated: [''],
      },
    };
  });

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
    const record = createOrUpdateDestination(schema, req);
    return destinationResponse(record);
  });
  server.patch(uri, (schema, req) => {
    const record = createOrUpdateDestination(schema, req);
    return destinationResponse(record);
  });
  server.delete(uri, (schema, req) => {
    const { type, name } = req.params;
    schema.db.syncDestinations.update(
      { type, name },
      // these parameters are added after a purge delete is initiated
      // if only `purge_initiated_at` exists the delete progress banner renders
      // if `purge_error` also has a value then delete failed banner renders
      {
        purge_initiated_at: '2024-01-09T16:54:28.463879-07:00',
        // WIP (backend hasn't added yet) update when we have a realistic error message)
        // purge_error: '1 error occurred: association could for some confusing reason not be un-synced!',
      }
    );
    const record = schema.db.syncDestinations.findBy({ type, name });
    return destinationResponse(record);
    // return the following instead to test immediate deletion
    // schema.db.syncDestinations.remove({ type, name });
    // return new Response(204);
  });
  // associations
  server.get('/sys/sync/associations', (schema) => {
    const associations = schema.db.syncAssociations.where({});
    if (!associations.length) {
      return new Response(404, {}, { errors: [] });
    }

    const secrets = associations.reduce((secrets, association) => {
      const secretPath = `${association.mount}/${association.secret_name}`;
      if (!secrets.includes(secretPath)) {
        secrets.push(secretPath);
      }
      return secrets;
    }, []);

    return {
      data: {
        key_info: {},
        keys: [],
        total_associations: associations.length, // link between a secret and a destination
        total_secrets: secrets.length, // number of secrets synced from vault (one secret can be synced to multiple destinations)
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
    schema.db.syncAssociations.update(
      { type, name },
      { ...data, sync_status: 'SYNCED', updated_at: new Date().toISOString() }
    );
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

  // SYNC CLIENTS ACTIVITY RESPONSE

  // DYNAMIC RESPONSE (with date querying)
  clientsHandler(server); // imports all of the endpoints defined in mirage/handlers/clients file

  // STATIC RESPONSE (0 entity/non-entity clients)
  /*
  server.get('/sys/internal/counters/activity', (schema, req) => {
    let { start_time, end_time } = req.queryParams;
    // backend returns a timestamp if given unix time, so first convert to timestamp string here
    if (!start_time.includes('T')) start_time = fromUnixTime(start_time).toISOString();
    if (!end_time.includes('T')) end_time = fromUnixTime(end_time).toISOString();
    return {
      request_id: 'some-activity-id',
      lease_id: '',
      renewable: false,
      lease_duration: 0,
      data: {
        start_time, // set by query params
        end_time, // set by query params
        total: {
          clients: 15,
          distinct_entities: 0,
          entity_clients: 0,
          non_entity_clients: 0,
          non_entity_tokens: 0,
          secret_syncs: 15,
        },
        by_namespace: [
          {
            counts: {
              clients: 15,
              distinct_entities: 0,
              entity_clients: 0,
              non_entity_clients: 0,
              non_entity_tokens: 0,
              secret_syncs: 15,
            },
            mounts: [
              {
                counts: {
                  clients: 15,
                  distinct_entities: 0,
                  entity_clients: 0,
                  non_entity_clients: 0,
                  non_entity_tokens: 0,
                  secret_syncs: 15,
                },
                mount_path: 'sys/',
              },
            ],
            namespace_id: 'root',
            namespace_path: '',
          },
        ],
        months: [
          { counts: null, namespaces: null, new_clients: null, timestamp: '2023-09-01T00:00:00Z' },
          {
            counts: {
              clients: 10,
              distinct_entities: 0,
              entity_clients: 0,
              non_entity_clients: 0,
              non_entity_tokens: 0,
              secret_syncs: 10,
            },
            namespaces: [
              {
                counts: {
                  clients: 10,
                  distinct_entities: 0,
                  entity_clients: 0,
                  non_entity_clients: 0,
                  non_entity_tokens: 0,
                  secret_syncs: 10,
                },
                mounts: [
                  {
                    counts: {
                      clients: 10,
                      distinct_entities: 0,
                      entity_clients: 0,
                      non_entity_clients: 0,
                      non_entity_tokens: 0,
                      secret_syncs: 10,
                    },
                    mount_path: 'sys/',
                  },
                ],
                namespace_id: 'root',
                namespace_path: '',
              },
            ],
            new_clients: {
              counts: {
                clients: 10,
                distinct_entities: 0,
                entity_clients: 0,
                non_entity_clients: 0,
                non_entity_tokens: 0,
                secret_syncs: 10,
              },
              namespaces: [
                {
                  counts: {
                    clients: 10,
                    distinct_entities: 0,
                    entity_clients: 0,
                    non_entity_clients: 0,
                    non_entity_tokens: 0,
                    secret_syncs: 10,
                  },
                  mounts: [
                    {
                      counts: {
                        clients: 10,
                        distinct_entities: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                        non_entity_tokens: 0,
                        secret_syncs: 10,
                      },
                      mount_path: 'sys/',
                    },
                  ],
                  namespace_id: 'root',
                  namespace_path: '',
                },
              ],
            },
            timestamp: '2023-10-01T00:00:00Z',
          },
          {
            counts: {
              clients: 7,
              distinct_entities: 0,
              entity_clients: 0,
              non_entity_clients: 0,
              non_entity_tokens: 0,
              secret_syncs: 7,
            },
            namespaces: [
              {
                counts: {
                  clients: 7,
                  distinct_entities: 0,
                  entity_clients: 0,
                  non_entity_clients: 0,
                  non_entity_tokens: 0,
                  secret_syncs: 7,
                },
                mounts: [
                  {
                    counts: {
                      clients: 7,
                      distinct_entities: 0,
                      entity_clients: 0,
                      non_entity_clients: 0,
                      non_entity_tokens: 0,
                      secret_syncs: 7,
                    },
                    mount_path: 'sys/',
                  },
                ],
                namespace_id: 'root',
                namespace_path: '',
              },
            ],
            new_clients: {
              counts: {
                clients: 3,
                distinct_entities: 0,
                entity_clients: 0,
                non_entity_clients: 0,
                non_entity_tokens: 0,
                secret_syncs: 3,
              },
              namespaces: [
                {
                  counts: {
                    clients: 3,
                    distinct_entities: 0,
                    entity_clients: 0,
                    non_entity_clients: 0,
                    non_entity_tokens: 0,
                    secret_syncs: 3,
                  },
                  mounts: [
                    {
                      counts: {
                        clients: 3,
                        distinct_entities: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                        non_entity_tokens: 0,
                        secret_syncs: 3,
                      },
                      mount_path: 'sys/',
                    },
                  ],
                  namespace_id: 'root',
                  namespace_path: '',
                },
              ],
            },
            timestamp: '2023-11-01T00:00:00Z',
          },
          {
            counts: {
              clients: 7,
              distinct_entities: 0,
              entity_clients: 0,
              non_entity_clients: 0,
              non_entity_tokens: 0,
              secret_syncs: 7,
            },
            namespaces: [
              {
                counts: {
                  clients: 7,
                  distinct_entities: 0,
                  entity_clients: 0,
                  non_entity_clients: 0,
                  non_entity_tokens: 0,
                  secret_syncs: 7,
                },
                mounts: [
                  {
                    counts: {
                      clients: 7,
                      distinct_entities: 0,
                      entity_clients: 0,
                      non_entity_clients: 0,
                      non_entity_tokens: 0,
                      secret_syncs: 7,
                    },
                    mount_path: 'sys/',
                  },
                ],
                namespace_id: 'root',
                namespace_path: '',
              },
            ],
            new_clients: {
              counts: {
                clients: 2,
                distinct_entities: 0,
                entity_clients: 0,
                non_entity_clients: 0,
                non_entity_tokens: 0,
                secret_syncs: 2,
              },
              namespaces: [
                {
                  counts: {
                    clients: 2,
                    distinct_entities: 0,
                    entity_clients: 0,
                    non_entity_clients: 0,
                    non_entity_tokens: 0,
                    secret_syncs: 2,
                  },
                  mounts: [
                    {
                      counts: {
                        clients: 2,
                        distinct_entities: 0,
                        entity_clients: 0,
                        non_entity_clients: 0,
                        non_entity_tokens: 0,
                        secret_syncs: 2,
                      },
                      mount_path: 'sys/',
                    },
                  ],
                  namespace_id: 'root',
                  namespace_path: '',
                },
              ],
            },
            timestamp: '2023-12-01T00:00:00Z',
          },
        ],
      },
      wrap_info: null,
      warnings: null,
      auth: null,
    };
  });

  */
}
