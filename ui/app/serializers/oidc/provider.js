import ApplicationSerializer from '../application';

export default class OidcProviderSerializer extends ApplicationSerializer {
  primaryKey = 'name';

  // need to normalize to get issuer metadata for provider's list view
  normalizeItems(payload) {
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
