/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import ApplicationAdapter from '../application';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default class KvDataAdapter extends ApplicationAdapter {
  namespace = 'v1';

  _urlForSecret(backend, path, version) {
    const base = `${this.buildURL()}/${encodePath(backend)}/data/${encodePath(path)}`;
    return version ? base + `?version=${version}` : base;
  }

  createRecord(store, type, snapshot) {
    const { backend, path, version } = snapshot.record;
    const url = this._urlForSecret(backend, path);
    return this.ajax(url, 'POST', { data: this.serialize(snapshot) }).then((resp) => {
      resp.id = `${encodePath(backend)}/${version}/${encodePath(path)}`;
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

  /* Four types of delete operations */
  // make typescript
  // ARG TODO use destroyRecord({deleteType}) and conditionally call the method.
  // util can be the names of these methods.  https://api.emberjs.com/ember-data/4.12/classes/Model/methods?anchor=destroyRecord

  // 1. Soft delete the secret's latest version.
  deleteLatestVersion(backend, path) {
    return this.ajax(this._urlForSecret(backend, path), 'DELETE');
  }

  // 2. Soft delete specific version(s) of the secret.
  deleteSpecificVersions(backend, path, versions) {
    return this.ajax(this._urlForSecret(backend, path), 'POST', {
      data: { versions },
    });
  }

  // 3. Permanently remove specific version(s) of a secret.
  destroySpecificVersions(backend, path, versions) {
    const url = `${this.buildURL()}/${encodePath(backend)}/destroy/${encodePath(path)}`;
    return this.ajax(url, 'PUT', {
      data: { versions },
    });
  }

  // 4. Permanently remove a secret's data and metadata.
  destroyEverything(backend, path) {
    const url = `${this.buildURL()}/${encodePath(backend)}/metadata/${encodePath(path)}`;
    return this.ajax(url, 'DELETE');
  }

  // Undelete a specific version(s)
  undeleteSpecificVersions(backend, path, versions) {
    const url = `${this.buildURL()}/${encodePath(backend)}/undelete/${encodePath(path)}`;
    return this.ajax(url, 'POST', {
      data: { versions },
    });
  }
}
