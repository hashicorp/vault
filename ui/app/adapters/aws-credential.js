import ApplicationAdapter from './application';

export default ApplicationAdapter.extend({
  createRecord(store, type, snapshot) {
    let ttl = snapshot.attr('ttl');
    let roleArn = snapshot.attr('roleArn');
    let data = {};
    if (ttl) {
      data.ttl = ttl;
    }
    if (roleArn) {
      data.role_arn = roleArn;
    }
    let method = ttl || roleArn ? 'POST' : 'GET';
    let options = ttl || roleArn ? { data } : {};
    let role = snapshot.attr('role');
    let url = `/v1/${role.backend}/creds/${role.name}`;

    return this.ajax(url, method, options).then(response => {
      response.id = snapshot.id;
      response.modelName = type.modelName;
      store.pushPayload(type.modelName, response);
    });
  },
});
