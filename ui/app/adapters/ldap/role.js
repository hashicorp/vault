/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import NamedPathAdapter from 'vault/adapters/named-path';
import { encodePath } from 'vault/utils/path-encoding-helpers';
import { inject as service } from '@ember/service';

export default class LdapRoleAdapter extends NamedPathAdapter {
  @service flashMessages;

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

  async query(store, type, query, recordArray, options) {
    const { showPartialError } = options.adapterOptions || {};
    const { backend } = query;
    const roles = [];
    const errors = [];

    for (const roleType of ['static', 'dynamic']) {
      const url = this.getURL(backend, this.pathForRoleType(roleType));
      try {
        const models = await this.ajax(url, 'GET', { data: { list: true } }).then((resp) => {
          return resp.data.keys.map((name) => ({ name, backend, type: roleType }));
        });
        roles.addObjects(models);
      } catch (error) {
        if (error.httpStatus !== 404) {
          errors.push(error);
        }
      }
    }

    if (errors.length) {
      const errorMessages = errors.reduce((errors, e) => {
        e.errors.forEach((error) => {
          errors.push(`${e.path}: ${error}`);
        });
        return errors;
      }, []);
      if (errors.length === 2) {
        // throw error as normal if both requests fail
        // ignore status code and concat errors to be displayed in Page::Error component with generic message
        throw { message: 'Error fetching roles:', errors: errorMessages };
      } else if (showPartialError) {
        // if only one request fails, surface the error to the user an info level flash message
        // this may help for permissions errors where a users policy may be incorrect
        this.flashMessages.info(`Error fetching roles from ${errorMessages.join(', ')}`);
      }
    }

    return roles.sortBy('name');
  }
  queryRecord(store, type, query) {
    const { backend, name, type: roleType } = query;
    const url = this.getURL(backend, this.pathForRoleType(roleType), name);
    return this.ajax(url, 'GET').then((resp) => ({ ...resp.data, backend, name, type: roleType }));
  }

  fetchCredentials(backend, type, name) {
    const url = this.getURL(backend, this.pathForRoleType(type, true), name);
    return this.ajax(url, 'GET').then((resp) => {
      if (type === 'dynamic') {
        const { lease_id, lease_duration, renewable } = resp;
        return { ...resp.data, lease_id, lease_duration, renewable, type };
      }
      return { ...resp.data, type };
    });
  }
  rotateStaticPassword(backend, name) {
    const url = this.getURL(backend, 'rotate-role', name);
    return this.ajax(url, 'POST');
  }
}
