import ApplicationSerializer from './application';

export default ApplicationSerializer.extend({
  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    if (payload.data.masking_character) {
      payload.data.masking_character = String.fromCharCode(payload.data.masking_character);
    }

    if (payload.data.templates && payload.data.templates.length === 1) {
      // add space after comma in returned array length one of a string of templates
      payload.data.templates[0] = payload.data.templates[0].replace(/,/g, ', ');
    }

    // TODO: something similar her for roles, however it's a little tough because each role comes back as an item in an array.
    // Also note, the BE is working on a ticket to amend these response, so revisit.
    // TODO: QA to make sure this is only normalizing the response and not changing data on edit (e.g. add a space after the comma when rePOST on edit)

    return this._super(store, primaryModelClass, payload, id, requestType);
  },
});
