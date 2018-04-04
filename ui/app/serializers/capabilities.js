import DS from 'ember-data';

export default DS.RESTSerializer.extend({
  primaryKey: 'path',

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    // queryRecord will already have set this, and we won't have an id here
    if (id) {
      payload.path = id;
    }
    const response = {
      [primaryModelClass.modelName]: payload,
    };
    return this._super(store, primaryModelClass, response, id, requestType);
  },

  modelNameFromPayloadKey() {
    return 'capabilities';
  },
});
