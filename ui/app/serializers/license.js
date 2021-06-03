import ApplicationSerializer from './application';

export default ApplicationSerializer.extend({
  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    console.log('normalize response', payload);
    let transformedPayload = { autoloaded: payload.autoloading_used, id: 'no-license' };
    if (payload.stored) {
      transformedPayload = {
        ...payload.stored,
        ...transformedPayload,
        id: payload.stored.license_id,
      };
    } else if (payload.autoloaded) {
      transformedPayload = {
        ...payload.autoloaded,
        ...transformedPayload,
        id: payload.autoloaded.license_id,
      };
    }
    return this._super(store, primaryModelClass, transformedPayload, id, requestType);
  },
});
