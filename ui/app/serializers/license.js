import ApplicationSerializer from './application';

export default ApplicationSerializer.extend({
  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    let transformedPayload = { autoloaded: payload.data.autoloading_used, license_id: 'no-license' };
    if (payload.data.autoloaded) {
      transformedPayload = {
        ...transformedPayload,
        ...payload.data.autoloaded,
      };
    }
    transformedPayload.id = transformedPayload.license_id;
    return this._super(store, primaryModelClass, transformedPayload, id, requestType);
  },
});
