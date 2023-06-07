/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import ApplicationAdapter from '../application';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default class KvDataAdapter extends ApplicationAdapter {
  namespace = 'v1';

  _urlForSecret(backend, path, version) {
    // path is "kv-test/2/my-secret"
    const base = `${this.buildURL()}/${encodePath(backend)}/data/${encodePath(path)}`;
    return version ? base + `?version=${version}` : base;
  }

  createRecord(store, type, snapshot) {
    const backend = snapshot.record.backend;
    const path = snapshot.attr('path');
    const version = snapshot.attr('version');
    const url = this._urlForSecret(backend, path);

    return this.ajax(url, 'POST', { data: this.serialize(snapshot) }).then((resp) => {
      resp.id = `${backend}/${version}/${path}`;
      return resp;
    });
  }

  updateRecord(store, type, snapshot) {
    const { backend, path, version } = snapshot.record;
    const data = this.serialize(snapshot);
    const url = this._urlForSecret(backend, path, version);
    return this.ajax(url, 'POST', { data });
  }

  query(store, type, query) {
    const { path, backend, version } = query;
    return this.ajax(this._urlForSecret(backend, path, version), 'GET');
  }

  queryRecord(store, type, query) {
    const { path, backend, version } = query;
    return this.ajax(this._urlForSecret(backend, path, version), 'GET');
  }

  /* Four type of delete operations */

  // 1. Soft delete the secret's latest version.
  deleteLatestVersion(backend, path) {
    return this.ajax(this._urlForSecret(backend, path), 'DELETE');
  }

  // 2. Soft delete specified version of the secret.
  deleteSpecificVersions(backend, path, versions) {
    return this.ajax(this._urlForSecret(backend, path), 'POST', {
      data: { versions },
    });
  }

  // 3. Permanently remove specific version data of a secret(s).
  destroySpecificVersions(backend, path, versions) {
    const url = `${this.buildURL()}/${encodePath(backend)}/destroy/${encodePath(path)}`;
    return this.ajax(url, 'PUT', {
      data: { versions },
    });
  }

  // 4. Permanently remove specific version data of a secret(s).
  destroyEverything(backend, path) {
    const url = `${this.buildURL()}/${encodePath(backend)}/metadata/${encodePath(path)}`;
    return this.ajax(url, 'DELETE');
  }

  // Undelete
  undeleteSpecificVersions(backend, path, versions) {
    const url = `${this.buildURL()}/${encodePath(backend)}/undelete/${encodePath(path)}`;
    return this.ajax(url, 'POST', {
      data: { versions },
    });
  }
}
