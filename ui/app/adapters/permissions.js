/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from './application';

export default ApplicationAdapter.extend({
  allowNsAccess(userRoot, currentNs) {
    if (!userRoot) return true;
    // current namespace is root but userRoot is not falsy
    if (!currentNs) return false;
    return currentNs.includes(userRoot);
  },
  query() {
    const allowAccess = this.allowNsAccess(this.auth.authData.userRootNamespace, this.namespaceService.path);
    if (!allowAccess) {
      // This triggers the resultant-acl banner to display
      throw new Error("Requested namespace does not include user's root namespace.");
    }
    return this.ajax(this.urlForQuery(), 'GET', { namespace: this.auth.authData.userRootNamespace });
  },

  urlForQuery() {
    return this.buildURL() + '/internal/ui/resultant-acl';
  },
});
