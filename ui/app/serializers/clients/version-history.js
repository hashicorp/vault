import ApplicationSerializer from '../application';

export default ApplicationSerializer.extend({
  normalizeFindAllResponse(store, primaryModelClass, payload, id, requestType) {
    let normalizedPayload = [];
    payload.keys.forEach((key) => {
      normalizedPayload.push({ id: key, ...payload.key_info[key] });
    });
    return this._super(store, primaryModelClass, normalizedPayload, id, requestType);
  },
});
