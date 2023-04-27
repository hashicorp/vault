/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { allSettled } from 'rsvp';
import ApplicationAdapter from '../application';
import ControlGroupError from 'vault/lib/control-group-error';

export default ApplicationAdapter.extend({
  namespace: 'v1',

  _staticCreds(backend, secret) {
    return this.ajax(
      `${this.buildURL()}/${encodeURIComponent(backend)}/static-creds/${encodeURIComponent(secret)}`,
      'GET'
    ).then((resp) => ({ ...resp, roleType: 'static' }));
  },

  _dynamicCreds(backend, secret) {
    return this.ajax(
      `${this.buildURL()}/${encodeURIComponent(backend)}/creds/${encodeURIComponent(secret)}`,
      'GET'
    ).then((resp) => ({ ...resp, roleType: 'dynamic' }));
  },

  fetchByQuery(store, query) {
    const { backend, secret } = query;
    if (query.roleType === 'static') {
      return this._staticCreds(backend, secret);
    } else if (query.roleType === 'dynamic') {
      return this._dynamicCreds(backend, secret);
    }
    return allSettled([this._staticCreds(backend, secret), this._dynamicCreds(backend, secret)]).then(
      ([staticResp, dynamicResp]) => {
        if (staticResp.state === 'rejected' && dynamicResp.state === 'rejected') {
          let reason = staticResp.reason;
          if (dynamicResp.reason instanceof ControlGroupError) {
            throw dynamicResp.reason;
          }
          if (reason?.httpStatus < dynamicResp.reason?.httpStatus) {
            reason = dynamicResp.reason;
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

  rotateRoleCredentials(backend, id) {
    return this.ajax(
      `${this.buildURL()}/${encodeURIComponent(backend)}/rotate-role/${encodeURIComponent(id)}`,
      'POST'
    );
  },
});
