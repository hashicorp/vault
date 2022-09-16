import ApplicationSerializer from '../application';

export default class PkiIssuerEngineSerializer extends ApplicationSerializer {
  primaryKey = 'id';

  // rehydrate each issuer model so all model attributes are accessible from the LIST response
  normalizeItems(payload) {
    // debugger;
    if (payload.data) {
      if (payload.data?.keys && Array.isArray(payload.data.keys)) {
        return payload.data.keys.map((key) => ({ name: key, ...payload.data.key_info[key] }));
      }
      Object.assign(payload, payload.data);
      delete payload.data;
    }
    return payload;
  }
}
