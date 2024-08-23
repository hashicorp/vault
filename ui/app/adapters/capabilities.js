/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AdapterError from '@ember-data/adapter/error';
import { set } from '@ember/object';
import ApplicationAdapter from './application';
import { sanitizePath } from 'core/utils/sanitize-path';

export default class CapabilitiesAdapter extends ApplicationAdapter {
  pathForType() {
    return 'capabilities-self';
  }

  _formatPath(path) {
    const { relativeNamespace } = this.namespaceService;
    if (!relativeNamespace) {
      return path;
    }
    // ensure original path doesn't have leading slash
    return `${relativeNamespace}/${path.replace(/^\//, '')}`;
  }

  async findRecord(store, type, id) {
    const paths = [this._formatPath(id)];
    return this.ajax(this.buildURL(type), 'POST', {
      data: { paths },
      namespace: sanitizePath(this.namespaceService.userRootNamespace),
    }).catch((e) => {
      if (e instanceof AdapterError) {
        set(e, 'policyPath', 'sys/capabilities-self');
      }
      throw e;
    });
  }

  queryRecord(store, type, query) {
    const { id } = query;
    if (!id) {
      return;
    }
    return this.findRecord(store, type, id).then((resp) => {
      resp.path = id;
      return resp;
    });
  }

  query(store, type, query) {
    const paths = query?.paths.map((p) => this._formatPath(p));
    return this.ajax(this.buildURL(type), 'POST', {
      data: { paths },
      namespace: sanitizePath(this.namespaceService.userRootNamespace),
    }).catch((e) => {
      if (e instanceof AdapterError) {
        set(e, 'policyPath', 'sys/capabilities-self');
      }
      throw e;
    });
  }
}
