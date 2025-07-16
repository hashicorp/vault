/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AdapterError from '@ember-data/adapter/error';
import { set } from '@ember/object';
import ApplicationAdapter from './application';
import { sanitizePath, sanitizeStart } from 'core/utils/sanitize-path';

export default class CapabilitiesAdapter extends ApplicationAdapter {
  pathForType() {
    return 'capabilities-self';
  }

  /* 
  users don't always have access to the capabilities-self endpoint in the current namespace,
  this can happen when logging in to a namespace and then navigating to a child namespace.
  adding "relativeNamespace" to the path and/or "this.namespaceService.userRootNamespace"
  to the request header ensures we are querying capabilities-self in the user's root namespace,
  which is where they are most likely to have their policy/permissions.
  */
  _formatPath(path) {
    const { relativeNamespace } = this.namespaceService;
    if (!relativeNamespace) {
      return path;
    }
    // ensure original path doesn't have leading slash
    return `${relativeNamespace}/${sanitizeStart(path)}`;
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
    const pathMap = query?.paths.reduce((mapping, path) => {
      const withNs = this._formatPath(path);
      if (withNs) {
        mapping[withNs] = path;
      }
      return mapping;
    }, {});

    return this.ajax(this.buildURL(type), 'POST', {
      data: { paths: Object.keys(pathMap) },
      namespace: sanitizePath(this.namespaceService.userRootNamespace),
    })
      .then((queryResult) => {
        if (queryResult) {
          // send the pathMap with the response so the serializer can normalize the paths to be relative to the namespace
          queryResult.pathMap = pathMap;
        }
        return queryResult;
      })
      .catch((e) => {
        if (e instanceof AdapterError) {
          set(e, 'policyPath', 'sys/capabilities-self');
        }
        throw e;
      });
  }
}
