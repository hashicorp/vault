import ApplicationSerializer from './application';

export default ApplicationSerializer.extend({
  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    const normalizedPayload = {
      id: payload.id,
      data: payload.data.counters,
    };
    return this._super(store, primaryModelClass, normalizedPayload, id, requestType);
  },
});
