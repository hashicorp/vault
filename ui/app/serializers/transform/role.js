import ApplicationSerializer from '../application';
export default ApplicationSerializer.extend({
  extractLazyPaginatedData(payload) {
    // TODO: do this for transform too?
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
