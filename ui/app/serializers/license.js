import ApplicationSerializer from './application';

export default ApplicationSerializer.extend({
  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    console.log('normalize response', payload);
    let transformedPayload = { autoloaded: payload.autoloading_used, license_id: 'no-license' };
    if (payload.autoloaded) {
      transformedPayload = {
        ...transformedPayload,
        ...payload.autoloaded,
      };
    } else if (payload.stored) {
      transformedPayload = {
        ...transformedPayload,
        ...payload.stored,
      };
    }
    transformedPayload.id = transformedPayload.license_id;
    return this._super(store, primaryModelClass, transformedPayload, id, requestType);
  },
});
