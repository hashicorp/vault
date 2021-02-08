import RESTSerializer from '@ember-data/serializer/rest';

export default RESTSerializer.extend({
  primaryKey: 'name',

  normalizeSecrets(payload) {
    if (payload.data.keys && Array.isArray(payload.data.keys)) {
      const roles = payload.data.keys.map(secret => {
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
    if (payload.data.type === 'static') {
      path = 'static-roles';
    }
    let database = [];
    if (payload.data.db_name) {
      database = [payload.data.db_name];
    }
    return {
      id: payload.secret,
      name: payload.secret,
      backend: payload.backend,
      database,
      path,
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
    if (
      undefined !== snapshot.attr(key) &&
      (snapshot.record.get('isNew') || snapshot.changedAttributes()[key])
    ) {
      this._super(snapshot, json, key, attributes);
    }
  },

  serialize(snapshot, requestType) {
    let data = this._super(snapshot, requestType);
    if (data.database) {
      const db = data.database[0];
      data.db_name = db;
      delete data.database;
    }

    return data;
  },
});
