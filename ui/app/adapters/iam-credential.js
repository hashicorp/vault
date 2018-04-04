import ApplicationAdapter from './application';

export default ApplicationAdapter.extend({
  namespace: 'v1',

  url(role, isSTS) {
    if (isSTS) {
      return `/v1/${role.backend}/sts/${role.name}`;
    }
    return `/v1/${role.backend}/creds/${role.name}`;
  },

  createRecord(store, type, snapshot) {
    const isSTS = snapshot.attr('withSTS');
    const options = isSTS ? { data: { ttl: snapshot.attr('ttl') } } : {};
    const method = isSTS ? 'POST' : 'GET';
    const role = snapshot.attr('role');
    const url = this.url(role, isSTS);

    return this.ajax(url, method, options).then(response => {
      response.id = snapshot.id;
      response.modelName = type.modelName;
      store.pushPayload(type.modelName, response);
    });
  },
});
