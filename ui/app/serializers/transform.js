import ApplicationSerializer from './application';

export default ApplicationSerializer.extend({
  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    if (payload.data.masking_character) {
      payload.data.masking_character = String.fromCharCode(payload.data.masking_character);
    }
    return this._super(store, primaryModelClass, payload, id, requestType);
  },

  serialize() {
    let json = this._super(...arguments);
    if (json.template && Array.isArray(json.template)) {
      // Transformations should only ever have one template
      json.template = json.template[0];
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
