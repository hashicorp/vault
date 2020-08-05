import ApplicationSerializer from './application';

export default ApplicationSerializer.extend({
  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    if (payload.data.masking_character) {
      // TODO: provide catch for when the entered something that is not UTF-16, see slack discussion on FE channel
      payload.data.masking_character = String.fromCharCode(payload.data.masking_character);
    }
    return this._super(store, primaryModelClass, payload, id, requestType);
  },
});
