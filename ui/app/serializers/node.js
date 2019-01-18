import { assign } from '@ember/polyfills';
import { decamelize } from '@ember/string';
import DS from 'ember-data';

export default DS.RESTSerializer.extend(DS.EmbeddedRecordsMixin, {
  keyForAttribute: function(attr) {
    return decamelize(attr);
  },

  pushPayload(store, payload) {
    const transformedPayload = this.normalizeResponse(
      store,
      store.modelFor('node'),
      payload,
      null,
      'findAll'
    );
    return store.push(transformedPayload);
  },

  nodeFromObject(name, payload) {
    const nodeObj = payload.nodes[name];
    return assign(nodeObj, {
      name,
      id: name,
    });
  },

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    const nodes = payload.nodes
      ? Object.keys(payload.nodes).map(name => this.nodeFromObject(name, payload))
      : [assign(payload, { id: '1' })];

    const transformedPayload = { nodes: nodes };

    return this._super(store, primaryModelClass, transformedPayload, id, requestType);
  },

  normalize(model, hash, prop) {
    hash.id = '1';
    return this._super(model, hash, prop);
  },
});
