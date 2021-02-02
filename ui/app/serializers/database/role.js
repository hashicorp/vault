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
    return {
      id: payload.secret,
      name: payload.secret,
      backend: payload.backend,
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
});
