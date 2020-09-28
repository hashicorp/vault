import ApplicationSerializer from '../application';

export default ApplicationSerializer.extend({
  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    const normalizedPayload = {
      id: payload.id,
      total: payload.data.total,
      endTime: new Date(payload.data.end_time),
      startTime: new Date(payload.data.start_time),
    };
    return this._super(store, primaryModelClass, normalizedPayload, id, requestType);
  },
});
