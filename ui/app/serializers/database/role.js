/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import RESTSerializer from '@ember-data/serializer/rest';

export default RESTSerializer.extend({
  primaryKey: 'name',

  normalizeSecrets(payload) {
    if (payload.data.keys && Array.isArray(payload.data.keys)) {
      const roles = payload.data.keys.map((secret) => {
        let type = 'dynamic';
        let path = 'roles';
        if (payload.data.staticRoles.includes(secret)) {
          type = 'static';
          path = 'static-roles';
        }
        return { name: secret, backend: payload.backend, type, path };
      });
      return roles;
    }
    let path = 'roles';
    if (payload.type === 'static') {
      path = 'static-roles';
    }
    let database = [];
    if (payload.data.db_name) {
      database = [payload.data.db_name];
    }
    // Copy to singular for MongoDB
    let creation_statement = '';
    let revocation_statement = '';
    if (payload.data.creation_statements) {
      creation_statement = payload.data.creation_statements[0];
    }
    if (payload.data.revocation_statements) {
      revocation_statement = payload.data.revocation_statements[0];
    }
    return {
      id: payload.id,
      backend: payload.backend,
      name: payload.id,
      type: payload.type,
      database,
      path,
      creation_statement,
      revocation_statement,
      ...payload.data,
    };
  },

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    const nullResponses = ['updateRecord', 'createRecord', 'deleteRecord'];
    const roles = nullResponses.includes(requestType)
      ? { name: id, backend: payload.backend }
      : this.normalizeSecrets(payload);
    const { modelName } = primaryModelClass;
    let transformedPayload = { [modelName]: roles };
    if (requestType === 'queryRecord') {
      transformedPayload = { [modelName]: roles };
    }
    return this._super(store, primaryModelClass, transformedPayload, id, requestType);
  },

  serializeAttribute(snapshot, json, key, attributes) {
    // Don't send values that are undefined
    if (snapshot.attr(key) !== undefined && (snapshot.record.isNew || snapshot.changedAttributes()[key])) {
      this._super(snapshot, json, key, attributes);
    }
  },

  serialize(snapshot, requestType) {
    const data = this._super(snapshot, requestType);
    if (data.database) {
      const db = data.database[0];
      data.db_name = db;
      delete data.database;
    }
    // This is necessary because the input for MongoDB is a json string
    // rather than an array, so we transpose that here
    if (data.creation_statement) {
      const singleStatement = data.creation_statement;
      data.creation_statements = [singleStatement];
      delete data.creation_statement;
    }
    if (data.revocation_statement) {
      const singleStatement = data.revocation_statement;
      data.revocation_statements = [singleStatement];
      delete data.revocation_statement;
    }

    return data;
  },
});
