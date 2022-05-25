import ApplicationSerializer from '../application';

export default ApplicationSerializer.extend({
  normalizeItems(payload) {
    if (payload.data.keys && Array.isArray(payload.data.keys)) {
      return payload.data.keys.map((key) => {
        let model = payload.data.key_info[key];
        model.id = key;
        return model;
      });
    }
  },
});
