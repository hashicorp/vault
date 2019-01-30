import ApplicationSerializer from './application';

export default ApplicationSerializer.extend({
  extractLazyPaginatedData(payload) {
    let ret;

    if (payload.zero_address_roles) {
      payload.zero_address_roles.forEach(role => {
        // mutate key_info object to add zero_address info
        payload.data.key_info[role].zero_address = true;
      });
    }
    if (!payload.data.key_info) {
      return payload.data.keys.map(key => {
        let model = {
          id: key,
        };
        if (payload.backend) {
          model.backend = payload.backend;
        }
        return model;
      });
    }

    ret = payload.data.keys.map(key => {
      let model = {
        id: key,
        key_type: payload.data.key_info[key].key_type,
        zero_address: payload.data.key_info[key].zero_address,
      };
      if (payload.backend) {
        model.backend = payload.backend;
      }
      return model;
    });
    delete payload.data.key_info;
    return ret;
  },
});
