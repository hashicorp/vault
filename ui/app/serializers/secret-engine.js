import DS from 'ember-data';
import Ember from 'ember';
const { decamelize } = Ember.String;

export default DS.RESTSerializer.extend({
  keyForAttribute: function(attr) {
    return decamelize(attr);
  },
  normalizeBackend(path, backend) {
    let struct = {};
    for (let attribute in backend) {
      struct[attribute] = backend[attribute];
    }
    //queryRecord adds path to the response
    if (path !== null && !struct.path) {
      struct.path = path;
    }

    if (struct.data) {
      struct = Ember.assign({}, struct, struct.data);
      delete struct.data;
    }
    // strip the trailing slash off of the path so we
    // can navigate to it without getting `//` in the url
    struct.id = struct.path.slice(0, -1);
    return struct;
  },

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    const isCreate = requestType === 'createRecord';
    const isFind = requestType === 'findRecord';
    const isQueryRecord = requestType === 'queryRecord';
    let backends;
    if (isCreate) {
      backends = payload.data;
    } else if (isFind) {
      backends = this.normalizeBackend(id + '/', payload.data);
    } else if (isQueryRecord) {
      backends = this.normalizeBackend(null, payload);
    } else {
      backends = Object.keys(payload.data).map(id => this.normalizeBackend(id, payload[id]));
    }

    const transformedPayload = { [primaryModelClass.modelName]: backends };
    return this._super(store, primaryModelClass, transformedPayload, id, requestType);
  },

  serialize() {
    return this._super(...arguments);
  },
});
