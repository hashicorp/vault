import Ember from 'ember';
import DS from 'ember-data';

const { decamelize } = Ember.String;

export default DS.RESTSerializer.extend(DS.EmbeddedRecordsMixin, {
  keyForAttribute: function(attr) {
    return decamelize(attr);
  },

  attrs: {
    nodes: { embedded: 'always' },
  },

  pushPayload(store, payload) {
    const transformedPayload = this.normalizeResponse(
      store,
      store.modelFor('cluster'),
      payload,
      null,
      'findAll'
    );
    return store.push(transformedPayload);
  },

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    // FIXME when multiple clusters lands
    const transformedPayload = {
      clusters: Ember.assign({ id: '1' }, payload.data || payload),
    };

    return this._super(store, primaryModelClass, transformedPayload, id, requestType);
  },
});
