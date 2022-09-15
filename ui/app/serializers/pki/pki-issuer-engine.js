import ApplicationSerializer from '../application';

export default class PkiIssuerEngineSerializer extends ApplicationSerializer {
  primaryKey = 'id';

  // rehydrate each issuer model so all model attributes are accessible from the LIST response
  // normalizeItems(payload) {
  //   payload.data.backend = 'blah';
  //   Object.assign(payload, payload.data);
  //   delete payload.data;
  //   // if (payload.data) {
  //   //   if (payload.data?.key_info) {
  //   //     // return payload.data.keys.map((key) => ({ id: key, ...payload.data.key_info[key] }));
  //   //   }
  //   //   Object.assign(payload, payload.data);
  //   //   delete payload.data;
  //   // }
  //   return payload;
  // }
}
