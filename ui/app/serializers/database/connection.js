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
    return {
      id: payload.id,
      name: payload.id,
      backend: payload.backend,
      ...payload.data,
      ...payload.data.connection_details,
    };
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
});
