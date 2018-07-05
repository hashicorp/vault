import DS from 'ember-data';

export default DS.RESTSerializer.extend({
  primaryKey: 'accessor',

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    // queryRecord will already have set this, and we won't have an id here
    return this._super(store, primaryModelClass, response, id, requestType);
  },

  modelNameFromPayloadKey() {
    return 'control-group';
  },
});
