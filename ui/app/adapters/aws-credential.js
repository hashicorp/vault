import ApplicationAdapter from './application';

export default ApplicationAdapter.extend({
  createRecord(store, type, snapshot) {
    let ttl = snapshot.attr('ttl');
    let roleArn = snapshot.attr('roleArn');
    let roleType = snapshot.attr('credentialType');
    let method = 'POST';
    let options;
    let data = {};
    if (roleType === 'iam_user') {
      method = 'GET';
    } else {
      if (ttl !== undefined) {
        data.ttl = ttl;
      }
      if (roleType === 'assumed_role' && roleArn) {
        data.role_arn = roleArn;
      }
      options = data.ttl || data.role_arn ? { data } : {};
    }
    let role = snapshot.attr('role');
    let url = `/v1/${role.backend}/creds/${role.name}`;

    return this.ajax(url, method, options).then(response => {
      response.id = snapshot.id;
      response.modelName = type.modelName;
      store.pushPayload(type.modelName, response);
    });
  },
});
