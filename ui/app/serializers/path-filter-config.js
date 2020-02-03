import { decamelize } from '@ember/string';
import DS from 'ember-data';

export default DS.RESTSerializer.extend({
  keyForAttribute: function(attr) {
    return decamelize(attr);
  },

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    const { modelName } = primaryModelClass;
    payload.data.id = id;
    const transformedPayload = { [modelName]: payload.data };
    return this._super(store, primaryModelClass, transformedPayload, id, requestType);
  },
});
