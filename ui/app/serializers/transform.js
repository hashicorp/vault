import ApplicationSerializer from './application';

export default ApplicationSerializer.extend({
  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    if (payload.data.masking_character) {
      payload.data.masking_character = String.fromCharCode(payload.data.masking_character);
    }
    // TODO: the BE is working on a ticket to amend these response, so revisit.
    return this._super(store, primaryModelClass, payload, id, requestType);
  },
});
