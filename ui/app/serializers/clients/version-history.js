import ApplicationSerializer from '../application';

export default ApplicationSerializer.extend({
  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    let payloadId = payload.keys[0];
    if (payloadId) {
      const normalizedPayload = {
        id: payloadId,
        ...payload,
      };
      return this._super(store, primaryModelClass, normalizedPayload, id, requestType);
    } else {
      return this._super(store, primaryModelClass, payload, id, requestType);
    }
  },
});
