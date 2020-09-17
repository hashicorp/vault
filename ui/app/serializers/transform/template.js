import ApplicationSerializer from '../application';

export default ApplicationSerializer.extend({
  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    payload.data.name = payload.id;
    if (payload.data.alphabet) {
      payload.data.alphabet = [payload.data.alphabet];
    }
    return this._super(store, primaryModelClass, payload, id, requestType);
  },

  serialize() {
    let json = this._super(...arguments);
    if (json.alphabet && Array.isArray(json.alphabet)) {
      // Templates should only ever have one alphabet
      json.alphabet = json.alphabet[0];
    }
    return json;
  },

  extractLazyPaginatedData(payload) {
    let ret;
    ret = payload.data.keys.map(key => {
      let model = {
        id: key,
      };
      if (payload.backend) {
        model.backend = payload.backend;
      }
      return model;
    });
    return ret;
  },
});
