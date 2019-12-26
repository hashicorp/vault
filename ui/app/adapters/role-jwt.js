import ApplicationAdapter from './application';
import { inject as service } from '@ember/service';
import { get } from '@ember/object';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default ApplicationAdapter.extend({
  router: service(),

  findRecord(store, type, id, snapshot) {
    let [path, role] = JSON.parse(id);
    path = encodePath(path);

    let namespace = get(snapshot, 'adapterOptions.namespace');
    let url = `/v1/auth/${path}/oidc/auth_url`;
    let redirect_uri = `${window.location.origin}${this.router.urlFor('vault.cluster.oidc-callback', {
      auth_path: path,
    })}`;

    if (namespace) {
      redirect_uri = `${window.location.origin}${this.router.urlFor(
        'vault.cluster.oidc-callback',
        { auth_path: path },
        { queryParams: { namespace } }
      )}`;
    }

    return this.ajax(url, 'POST', {
      data: {
        role,
        redirect_uri,
      },
    });
  },
});
