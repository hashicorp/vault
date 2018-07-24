import ApplicationSerializer from './application';

export default ApplicationSerializer.extend({
  normalizeList(payload) {
    let pk = 'path';
    const data = payload.data.keys
      ? payload.data.keys.map(key => ({
          [pk]: key,
          id: key.slice(0, -1),
        }))
      : payload.data;

    return data;
  },

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    const nullResponses = ['deleteRecord', 'createRecord'];
    let cid = id || payload.id;
    let normalizedPayload = nullResponses.includes(requestType) ? { id: cid } : this.normalizeList(payload);
    return this._super(store, primaryModelClass, normalizedPayload, id, requestType);
  },
});
