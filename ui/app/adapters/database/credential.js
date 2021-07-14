import RSVP from 'rsvp';
import ApplicationAdapter from '../application';

export default ApplicationAdapter.extend({
  namespace: 'v1',

  _staticCreds(backend, secret) {
    return this.ajax(
      `${this.buildURL()}/${encodeURIComponent(backend)}/static-creds/${encodeURIComponent(secret)}`,
      'GET'
    ).then(resp => ({ ...resp, roleType: 'static' }));
  },

  _dynamicCreds(backend, secret) {
    return this.ajax(
      `${this.buildURL()}/${encodeURIComponent(backend)}/creds/${encodeURIComponent(secret)}`,
      'GET'
    ).then(resp => ({ ...resp, roleType: 'dynamic' }));
  },

  fetchByQuery(store, query) {
    const { backend, secret } = query;
    return RSVP.allSettled([this._staticCreds(backend, secret), this._dynamicCreds(backend, secret)]).then(
      ([staticResp, dynamicResp]) => {
        // If one comes back with wrapped response from control group, throw it
        const accessor = staticResp.accessor || dynamicResp.accessor;
        if (accessor) {
          throw accessor;
        }
        // if neither has payload, throw reason with highest httpStatus
        if (!staticResp.value && !dynamicResp.value) {
          let reason = dynamicResp.reason;
          if (reason?.httpStatus < staticResp.reason?.httpStatus) {
            reason = staticResp.reason;
          }
          throw reason;
        }
        // Otherwise, return whichever one has a value
        return staticResp.value || dynamicResp.value;
      }
    );
  },

  queryRecord(store, type, query) {
    return this.fetchByQuery(store, query);
  },
});
