import ApplicationSerializer from './application';
export default ApplicationSerializer.extend({
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

  normalizeItems() {
    let normalized = this._super(...arguments);
    // most roles will only have one in this array,
    // we'll default to the first, and keep the array on the
    // model and show a warning if there's more than one so that
    // they don't inadvertently save
    if (normalized.credential_types) {
      normalized.credential_type = normalized.credential_types[0];
    }
    return normalized;
  },
});
