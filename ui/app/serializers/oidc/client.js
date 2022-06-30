import ApplicationSerializer from '../application';

export default class OidcClientSerializer extends ApplicationSerializer {
  primaryKey = 'name';

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    let transformedPayload = {
      data: { ...payload.data, key: [payload.data.key] },
    };
    return super.normalizeResponse(store, primaryModelClass, transformedPayload, id, requestType);
  }
}
