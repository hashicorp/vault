/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from './application';

export default ApplicationAdapter.extend({
  query() {
    const namespace = this.namespaceService.userRootNamespace || this.namespaceService.path;
    return this.ajax(this.urlForQuery(), 'GET', { namespace });
  },

  urlForQuery() {
    return this.buildURL() + '/internal/ui/resultant-acl';
  },
});
