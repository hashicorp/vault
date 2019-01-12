import ApplicationAdapter from './application';
import { inject as service } from '@ember/service';

export default ApplicationAdapter.extend({
  router: service(),

  findRecord(store, type, id) {
    let [path, role] = JSON.parse(id);
    let url = `/v1/auth/${path}/oidc/auth_url`;
    let redirect_uri = `${window.location.origin}${this.router.urlFor('vault.cluster.oidc-callback', {
      auth_path: path,
    })}`;

    return this.ajax(url, 'POST', {
      data: {
        role,
        redirect_uri,
      },
    });
  },
});
