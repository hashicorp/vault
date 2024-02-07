/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from './application';

export default ApplicationAdapter.extend({
  query() {
    const namespace = this.namespaceService.userRootNamespace || this.namespaceService.path;
    return this.ajax(this.urlForQuery(), 'GET', {
      namespace,
    }).then((resp) => {
      resp.data.validNamespace = true;
      return resp;
    });
  },

  urlForQuery() {
    return this.buildURL() + '/internal/ui/resultant-acl';
  },
});
