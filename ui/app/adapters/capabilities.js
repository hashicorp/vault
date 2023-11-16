/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AdapterError from '@ember-data/adapter/error';
import { set } from '@ember/object';
import ApplicationAdapter from './application';
import { sanitizePath } from 'core/utils/sanitize-path';

export default ApplicationAdapter.extend({
  pathForType() {
    return 'capabilities-self';
  },

  formatPaths(path) {
    const { path: ns, userRootNamespace, relativeNamespace } = this.namespaceService;
    if (userRootNamespace === ns) {
      // If they match, the user is in their root namespace and we can fetch normally
      return [path];
    }
    // ensure original path doesn't have leading slash
    return [`/${relativeNamespace}/${path.replace(/^\//, '')}`];
  },

  findRecord(store, type, id) {
    const paths = this.formatPaths(id);
    return this.ajax(this.buildURL(type), 'POST', {
      data: { paths },
      namespace: sanitizePath(this.namespaceService.userRootNamespace),
    }).catch((e) => {
      if (e instanceof AdapterError) {
        set(e, 'policyPath', 'sys/capabilities-self');
      }
      throw e;
    });
  },

  queryRecord(store, type, query) {
    const { id } = query;
    if (!id) {
      return;
    }
    return this.findRecord(store, type, id).then((resp) => {
      resp.path = id;
      return resp;
    });
  },
});
