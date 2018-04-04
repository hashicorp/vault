import DS from 'ember-data';
import Ember from 'ember';
const { decamelize } = Ember.String;

export default DS.RESTSerializer.extend({
  keyForAttribute: function(attr) {
    return decamelize(attr);
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

  normalizeItems(payload) {
    Ember.assign(payload, payload.data);
    delete payload.data;
    return payload;
  },

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    const responseJSON = this.normalizeItems(payload);
    const { modelName } = primaryModelClass;
    let transformedPayload = { [modelName]: responseJSON };
    let ret = this._super(store, primaryModelClass, transformedPayload, id, requestType);
    return ret;
  },

  serializeAttribute(snapshot, json, key, attributes) {
    const val = snapshot.attr(key);
    if (attributes.options.readOnly) {
      return;
    }
    if (
      attributes.type === 'object' &&
      val &&
      Object.keys(val).length > 0 &&
      Ember.isNone(snapshot.changedAttributes()[key])
    ) {
      return;
    }
    if (Ember.isBlank(val) && Ember.isNone(snapshot.changedAttributes()[key])) {
      return;
    }

    this._super(snapshot, json, key, attributes);
  },
});
