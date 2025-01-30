/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
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
    const keys = records.map(({ name }) => {
      if (name.includes('/')) {
        const [parent, child] = name.split('/');
        // query.name is only passed by listOrGetRecord and means we want to list children of admin/
        // otherwise this is the request for all roles in an engine so we return the top-level paths
        return query?.name ? child : `${parent}/`;
      }
      return name;
    });

    return {
      data: { keys },
    };
  };

  const listOrGetRecord = (schema, req, type) => {
    const dbKey = type ? 'ldapRoles' : 'ldapLibraries';
    const query = type ? { type, name: `admin/child-${type}-role` } : { name: 'admin/test-library' };
    if (req.queryParams.list) {
      // the mirage database has setup all hierarchical names to be prefixed with "admin/"
      // while passing a query with specific name is not flexible, for simplicity
      // we only seeded the mirage db with one hierarchical resource for each role and a library
      return listRecords(schema, dbKey, query);
    }
    // otherwise we want to view details for a specific resource
    return getRecord(schema, req, dbKey);
  };

  // config
  server.post('/:backend/config', (schema, req) => createOrUpdateRecord(schema, req, 'ldapConfigs'));
  server.get('/:backend/config', (schema, req) => getRecord(schema, req, 'ldapConfigs'));
  // roles
  server.post('/:backend/static-role/:name', (schema, req) => createOrUpdateRecord(schema, req, 'ldapRoles'));
  server.post('/:backend/role/:name', (schema, req) => createOrUpdateRecord(schema, req, 'ldapRoles'));
  // if the role is hierarchical the name ends in a forward slash so we make a list request
  server.get('/:backend/static-role/*name', (schema, req) => listOrGetRecord(schema, req, 'static'));
  server.get('/:backend/role/*name', (schema, req) => listOrGetRecord(schema, req, 'dynamic'));

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
  server.get('/:backend/library/*name', (schema, req) => listOrGetRecord(schema, req));
  server.get('/:backend/library', (schema) => listRecords(schema, 'ldapLibraries'));
  server.get('/:backend/library/*name/status', (schema) => {
    const data = schema.db['ldapAccountStatuses'].reduce((prev, curr) => {
      prev[curr.account] = {
        available: curr.available,
        borrower_client_token: curr.borrower_client_token,
      };
      return prev;
    }, {});
    return { data };
  });
  // check-out / check-in
  server.post('/:backend/library/:set_name/check-in', (schema, req) => {
    // Check-in makes an unavailable account available again
    const { service_account_names } = JSON.parse(req.requestBody);
    const dbCollection = schema.db['ldapAccountStatuses'];
    const updated = dbCollection.find(service_account_names).map((f) => ({
      ...f,
      available: true,
      borrower_client_token: undefined,
    }));
    updated.forEach((u) => {
      dbCollection.update(u.id, u);
    });
    return {
      data: {
        check_ins: service_account_names,
      },
    };
  });
  server.post('/:backend/library/:set_name/check-out', (schema, req) => {
    const { set_name, backend } = req.params;
    const dbCollection = schema.db['ldapAccountStatuses'];
    const available = dbCollection.where({ available: true });
    if (available) {
      return Response(404, {}, { errors: ['no accounts available to check out'] });
    }
    const checkOut = {
      ...available[0],
      available: false,
      borrower_client_token: crypto.randomUUID(),
    };
    dbCollection.update(checkOut.id, checkOut);
    return {
      request_id: '364a17d4-e5ab-998b-ceee-b49929229e0c',
      lease_id: `${backend}/library/${set_name}/check-out/aoBsaBEI4PK96VnukubvYDlZ`,
      renewable: true,
      lease_duration: 36000,
      data: {
        password: crypto.randomUUID(),
        service_account_name: checkOut.account,
      },
      wrap_info: null,
      warnings: null,
      auth: null,
    };
  });
}
