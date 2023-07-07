/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import NamedPathAdapter from 'vault/adapters/named-path';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default class LdapRoleAdapter extends NamedPathAdapter {
  getURL(backend, path, name) {
    const base = `${this.buildURL()}/${encodePath(backend)}/${path}`;
    return name ? `${base}/${name}` : base;
  }
  pathForRoleType(type, isCred) {
    const staticPath = isCred ? 'static-cred' : 'static-role';
    const dynamicPath = isCred ? 'creds' : 'role';
    return type === 'static' ? staticPath : dynamicPath;
  }

  urlForUpdateRecord(name, modelName, snapshot) {
    const { backend, type } = snapshot.record;
    return this.getURL(backend, this.pathForRoleType(type), name);
  }
  urlForDeleteRecord(name, modelName, snapshot) {
    const { backend, type } = snapshot.record;
    return this.getURL(backend, this.pathForRoleType(type), name);
  }

  query(store, type, query) {
    const { backend, type: roleType } = query;
    const url = this.getURL(backend, this.pathForRoleType(roleType));
    return this.ajax(url, 'GET', { data: { list: true } }).then((resp) => {
      return resp.data.keys.map((name) => ({ name, backend }));
    });
  }
  queryRecord(store, type, query) {
    const { backend, name, type: roleType } = query;
    const url = this.getURL(backend, this.pathForRoleType(roleType), name);
    return this.ajax(url, 'GET').then((resp) => {
      resp.data.backend = backend;
      resp.data.name = name;
      return resp.data;
    });
  }

  fetchCredentials(backend, type, name) {
    const url = this.getURL(backend, this.pathForRoleType(type, true), name);
    return this.ajax(url, 'GET');
  }
  rotateStaticPassword(backend, name) {
    const url = this.getURL(backend, 'rotate-role', name);
    return this.ajax(url, 'POST');
  }
}
