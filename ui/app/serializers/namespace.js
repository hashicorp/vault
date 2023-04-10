import ApplicationSerializer from './application';

export default class NamespaceSerializer extends ApplicationSerializer {
  attrs = {
    path: { serialize: false },
  };

  normalizeList(payload) {
    const data = payload.data.keys
      ? payload.data.keys.map((key) => ({
          path: key,
          // remove the trailing slash from the id
          id: key.replace(/\/$/, ''),
        }))
      : payload.data;

    return data;
  }

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    const nullResponses = ['deleteRecord', 'createRecord'];
    const cid = (id || payload.id || '').replace(/\/$/, '');
    const normalizedPayload = nullResponses.includes(requestType)
      ? { id: cid, path: cid }
      : this.normalizeList(payload);
    return super.normalizeResponse(store, primaryModelClass, normalizedPayload, id, requestType);
  }
}
