import ApplicationSerializer from '../application';

export default class PkiOverviewSerializer extends ApplicationSerializer {
  normalizeItems(payload) {
    return payload;
    // if (payload.data && payload.data.keys && Array.isArray(payload.data.keys)) {
    //   let ret = payload.data.keys.map((key) => {
    //     let model = {
    //       id_for_nav: `cert/${key}`,
    //       id: key,
    //     };
    //     if (payload.backend) {
    //       model.backend = payload.backend;
    //     }
    //     return model;
    //   });
    //   return ret;
    // }
    // assign(payload, payload.data);
    // delete payload.data;
    // return payload;
  }
}
