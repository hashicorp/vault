import ApplicationSerializer from '../application';
import Ember from 'ember';

export default ApplicationSerializer.extend({
  normalizeItems(payload) {
    if (payload.data.keys && Array.isArray(payload.data.keys)) {
      return payload.data.keys;
    }
    Ember.assign(payload, payload.data);
    delete payload.data;
    return payload;
  },

  extractLazyPaginatedData(payload) {
    let list;
    list = payload.data.keys.map(key => {
      let model = payload.data.key_info[key];
      model.id = key;
      return model;
    });
    delete payload.data.key_info;
    return list.sort((a, b) => {
      return a.name.localeCompare(b.name);
    });
  },
});
