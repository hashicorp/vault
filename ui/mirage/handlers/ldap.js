/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Response } from 'miragejs';

export default function (server) {
  const getRecord = (schema, req, dbKey, type) => {
    const { name } = req.params;
    const factoryKey = dbKey === 'ldapRoles' ? 'ldap-role' : 'ldap-library';
    const record = schema.db[dbKey].findBy({ name }) || server.create(factoryKey, type);
    delete record.id;
    delete record.name;
    return record ? { data: record } : new Response(404, {}, { errors: [] });
  };
  const createOrUpdateRecord = (schema, req, dbKey) => {
    const { name } = req.params;
    const data = JSON.parse(req.requestBody);
    const dbCollection = schema.db[dbKey];
    dbCollection.firstOrCreate({ name }, data);
    dbCollection.update({ name }, data);
    return new Response(204);
  };

  server.post('/:backend/static-role/:name', (schema, req) => createOrUpdateRecord(schema, req, 'ldapRoles'));
  server.post('/:backend/role/:name', (schema, req) => createOrUpdateRecord(schema, req, 'ldapRoles'));
  server.get('/:backend/static-role/:name', (schema, req) => getRecord(schema, req, 'ldapRoles', 'static'));
  server.get('/:backend/role/:name', (schema, req) => getRecord(schema, req, 'ldapRoles', 'dynamic'));
}
