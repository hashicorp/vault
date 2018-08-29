import ApplicationSerializer from './application';

export default ApplicationSerializer.extend({
  normalizeList(payload) {
    const data = payload.data.keys
      ? payload.data.keys.map(key => ({
          path: key,
          // remove the trailing slash from the id
          id: key.replace(/\/$/, ''),
        }))
      : payload.data;

    return data;
  },

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    const nullResponses = ['deleteRecord', 'createRecord'];
    let cid = (id || payload.id || '').replace(/\/$/, '');
    let normalizedPayload = nullResponses.includes(requestType)
      ? { id: cid, path: cid }
      : this.normalizeList(payload);
    return this._super(store, primaryModelClass, normalizedPayload, id, requestType);
  },
});
