import { isNone, isBlank } from '@ember/utils';
import { assign } from '@ember/polyfills';
import { decamelize } from '@ember/string';
import DS from 'ember-data';

export default DS.JSONSerializer.extend({
  keyForAttribute: function(attr) {
    return decamelize(attr);
  },

  normalizeItems(payload) {
    if (payload.data && payload.data.keys && Array.isArray(payload.data.keys)) {
      let models = payload.data.keys.map(key => {
        if (typeof key !== 'string') {
          return key;
        }
        let pk = this.get('primaryKey') || 'id';
        return { [pk]: key };
      });
      return models;
    }
    assign(payload, payload.data);
    delete payload.data;
    return payload;
  },

  pushPayload(store, payload) {
    const transformedPayload = this.normalizeResponse(
      store,
      store.modelFor(payload.modelName),
      payload,
      payload.id,
      'findRecord'
    );
    return store.push(transformedPayload);
  },

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    const responseJSON = this.normalizeItems(payload, requestType);
    if (id && !responseJSON.id) {
      responseJSON.id = id;
    }
    return this._super(store, primaryModelClass, responseJSON, id, requestType);
  },

  serializeAttribute(snapshot, json, key, attributes) {
    const val = snapshot.attr(key);
    const valHasNotChanged = isNone(snapshot.changedAttributes()[key]);
    const valIsBlank = isBlank(val);
    if (attributes.options.readOnly) {
      return;
    }
    if (attributes.type === 'object' && val && Object.keys(val).length > 0 && valHasNotChanged) {
      return;
    }
    if (valIsBlank && valHasNotChanged) {
      return;
    }

    this._super(snapshot, json, key, attributes);
  },

  serializeBelongsTo(snapshot, json) {
    return json;
  },
});
