import ApplicationSerializer from './application';

export default ApplicationSerializer.extend({
  normalizeFindAllResponse(store, primaryModelClass, payload, id, requestType) {
    let normalizedPayload = [];
    payload.forEach((data) => {
      data.keys.forEach((key) => {
        normalizedPayload.push({ id: key, ...data.key_info[key] });
      });
    });
    return this._super(store, primaryModelClass, normalizedPayload, id, requestType);
  },
});
