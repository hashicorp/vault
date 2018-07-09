import ApplicationSerializer from './application';

export default ApplicationSerializer.extend({
  normalizeBackend(path, backend) {
    let struct = { ...backend };
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
      : Object.keys(payload.data).map(path => this.normalizeBackend(path, payload.data[path]));

    return this._super(store, primaryModelClass, backends, id, requestType);
  },
});
