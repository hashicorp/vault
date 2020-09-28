import ApplicationSerializer from '../application';

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
});
