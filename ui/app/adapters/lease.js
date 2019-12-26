import ApplicationAdapter from './application';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default ApplicationAdapter.extend({
  revokePrefix(prefix) {
    let url = this.buildURL() + '/leases/revoke-prefix/' + encodePath(prefix);
    url = url.replace(/\/$/, '');
    return this.ajax(url, 'PUT');
  },

  forceRevokePrefix(prefix) {
    let url = this.buildURL() + '/leases/revoke-prefix/' + encodePath(prefix);
    url = url.replace(/\/$/, '');
    return this.ajax(url, 'PUT');
  },

  renew(lease_id, interval) {
    let url = this.buildURL() + '/leases/renew';
    return this.ajax(url, 'PUT', {
      data: {
        lease_id,
        interval,
      },
    });
  },

  deleteRecord(store, type, snapshot) {
    const lease_id = snapshot.id;
    return this.ajax(this.buildURL() + '/leases/revoke', 'PUT', {
      data: {
        lease_id,
      },
    });
  },

  queryRecord(store, type, query) {
    const { lease_id } = query;
    return this.ajax(this.buildURL() + '/leases/lookup', 'PUT', {
      data: {
        lease_id,
      },
    });
  },

  query(store, type, query) {
    const prefix = query.prefix || '';
    return this.ajax(this.buildURL() + '/leases/lookup/' + encodePath(prefix), 'GET', {
      data: {
        list: true,
      },
    }).then(resp => {
      if (prefix) {
        resp.prefix = prefix;
      }
      return resp;
    });
  },
});
