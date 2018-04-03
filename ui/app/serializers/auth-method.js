import ApplicationSerializer from './application';

export default ApplicationSerializer.extend({
  normalizeBackend(path, backend) {
    let struct = {};
    for (let attribute in backend) {
      struct[attribute] = backend[attribute];
    }
    // strip the trailing slash off of the path so we
    // can navigate to it without getting `//` in the url
    struct.id = path.slice(0, -1);
    struct.path = path;
    return struct;
  },

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    const isCreate = requestType === 'createRecord';
    const backends = isCreate
      ? payload.data
      : Object.keys(payload.data).map(id => this.normalizeBackend(id, payload[id]));

    return this._super(store, primaryModelClass, backends, id, requestType);
  },
});
