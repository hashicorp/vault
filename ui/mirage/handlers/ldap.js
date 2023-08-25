/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Response } from 'miragejs';

export default function (server) {
  const query = (req) => {
    const { name, backend } = req.params;
    return name ? { name } : { backend };
  };
  const getRecord = (schema, req, dbKey) => {
    const record = schema.db[dbKey].findBy(query(req));
    if (record) {
      delete record.id;
      delete record.name;
      delete record.backend;
      delete record.type;
      return { data: record };
    }
    return new Response(404, {}, { errors: [] });
  };
  const createOrUpdateRecord = (schema, req, dbKey) => {
    const data = JSON.parse(req.requestBody);
    const dbCollection = schema.db[dbKey];
    dbCollection.firstOrCreate(query(req), data);
    dbCollection.update(query(req), data);
    return new Response(204);
  };
  const listRecords = (schema, dbKey, query = {}) => {
    const records = schema.db[dbKey].where(query);
    return {
      data: { keys: records.map((record) => record.name) },
    };
  };

  // mount
  server.post('/sys/mounts/:path', () => new Response(204));
  server.get('/sys/internal/ui/mounts/:path', () => ({
    data: {
      accessor: 'ldap_ade94329',
      type: 'ldap',
      path: 'ldap-test/',
      uuid: '35e9119d-5708-4b6b-58d2-f913e27f242d',
      config: {},
    },
  }));
  // config
  server.post('/:backend/config', (schema, req) => createOrUpdateRecord(schema, req, 'ldapConfigs'));
  server.get('/:backend/config', (schema, req) => getRecord(schema, req, 'ldapConfigs'));
  // roles
  server.post('/:backend/static-role/:name', (schema, req) => createOrUpdateRecord(schema, req, 'ldapRoles'));
  server.post('/:backend/role/:name', (schema, req) => createOrUpdateRecord(schema, req, 'ldapRoles'));
  server.get('/:backend/static-role/:name', (schema, req) => getRecord(schema, req, 'ldapRoles', 'static'));
  server.get('/:backend/role/:name', (schema, req) => getRecord(schema, req, 'ldapRoles', 'dynamic'));
  server.get('/:backend/static-role', (schema) => listRecords(schema, 'ldapRoles', { type: 'static' }));
  server.get('/:backend/role', (schema) => listRecords(schema, 'ldapRoles', { type: 'dynamic' }));
  // role credentials
  server.get('/:backend/static-cred/:name', (schema) => ({
    data: schema.db.ldapCredentials.firstOrCreate({ type: 'static' }),
  }));
  server.get('/:backend/creds/:name', (schema) => ({
    data: schema.db.ldapCredentials.firstOrCreate({ type: 'dynamic' }),
  }));
  // libraries
  server.post('/:backend/library/:name', (schema, req) => createOrUpdateRecord(schema, req, 'ldapLibraries'));
  server.get('/:backend/library/:name', (schema, req) => getRecord(schema, req, 'ldapLibraries'));
  server.get('/:backend/library', (schema) => listRecords(schema, 'ldapLibraries'));
  server.get('/:backend/library/:name/status', () => ({
    'bob.johnson': { available: false, borrower_client_token: '8b80c305eb3a7dbd161ef98f10ea60a116ce0910' },
    'mary.smith': { available: true },
  }));
}
