import ApplicationAdapter from './application';

export default ApplicationAdapter.extend({
  findRecord(store, type, id) {
    let [path, role] = JSON.parse(id);
    let url = `/v1/auth/${path}/oidc/auth_url`;

    return this.ajax(url, 'POST', { data: { role } });
  },
});
