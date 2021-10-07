import RESTSerializer from '@ember-data/serializer/rest';

export default RESTSerializer.extend({
  primaryKey: 'name',

  serializeAttribute(snapshot, json, key, attributes) {
    // Don't send values that are undefined
    if (
      undefined !== snapshot.attr(key) &&
      (snapshot.record.get('isNew') || snapshot.changedAttributes()[key])
    ) {
      this._super(snapshot, json, key, attributes);
    }
  },

  normalizeSecrets(payload) {
    if (payload.data.keys && Array.isArray(payload.data.keys)) {
      const connections = payload.data.keys.map(secret => ({ name: secret, backend: payload.backend }));
      return connections;
    }
    // Query single record response:
    let response = {
      id: payload.id,
      name: payload.id,
      backend: payload.backend,
      ...payload.data,
      ...payload.data.connection_details,
    };
    if (payload.data.root_credentials_rotate_statements) {
      response.root_rotation_statements = payload.data.root_credentials_rotate_statements;
    }
    return response;
  },

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    const nullResponses = ['updateRecord', 'createRecord', 'deleteRecord'];
    const connections = nullResponses.includes(requestType)
      ? { name: id, backend: payload.backend }
      : this.normalizeSecrets(payload);
    const { modelName } = primaryModelClass;
    let transformedPayload = { [modelName]: connections };
    if (requestType === 'queryRecord') {
      // comes back as object anyway
      transformedPayload = { [modelName]: { id, ...connections } };
    }
    return this._super(store, primaryModelClass, transformedPayload, id, requestType);
  },

  // add serialize function this._super
  // serialize(snapshot, requestType) {
  //   let data = this._super(snapshot, requestType);
  //   if (data.database) {
  //     const db = data.database[0];
  //     data.db_name = db;
  //     delete data.database;
  //   }
  //   // This is necessary because the input for MongoDB is a json string
  //   // rather than an array, so we transpose that here
  //   if (data.creation_statement) {
  //     const singleStatement = data.creation_statement;
  //     data.creation_statements = [singleStatement];
  //     delete data.creation_statement;
  //   }
  //   if (data.revocation_statement) {
  //     const singleStatement = data.revocation_statement;
  //     data.revocation_statements = [singleStatement];
  //     delete data.revocation_statement;
  //   }

  //   return data;
  // },
});
