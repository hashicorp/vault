/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from 'vault/adapters/application';
import { encodePath } from 'vault/utils/path-encoding-helpers';
import { service } from '@ember/service';
import AdapterError from '@ember-data/adapter/error';
import { addManyToArray } from 'vault/helpers/add-to-array';
import sortObjects from 'vault/utils/sort-objects';

export const ldapRoleID = (type, name) => `type:${type}::name:${name}`;

export default class LdapRoleAdapter extends ApplicationAdapter {
  namespace = 'v1';

  @service flashMessages;

  // we do this in the adapter because query() requests separate endpoints to fetch static and dynamic roles.
  // it also handles some error logic and serializing (some of which is for lazyPaginatedQuery)
  // so for consistency formatting the response here
  _constructRecord({ backend, name, type }) {
    // ID cannot just be the 'name' because static and dynamic roles can have identical names
    return { id: ldapRoleID(type, name), backend, name, type };
  }

  _getURL(backend, path, name) {
    const base = `${this.buildURL()}/${encodePath(backend)}/${path}`;
    return name ? `${base}/${name}` : base;
  }

  _pathForRoleType(type, isCred) {
    const staticPath = isCred ? 'static-cred' : 'static-role';
    const dynamicPath = isCred ? 'creds' : 'role';
    return type === 'static' ? staticPath : dynamicPath;
  }

  _createOrUpdate(store, modelSchema, snapshot) {
    const { backend, name, type } = snapshot.record;
    const data = snapshot.serialize();
    return this.ajax(this._getURL(backend, this._pathForRoleType(type), name), 'POST', {
      data,
    }).then(() => {
      // add ID to response because ember data dislikes 204s...
      return { data: this._constructRecord({ backend, name, type }) };
    });
  }

  createRecord() {
    return this._createOrUpdate(...arguments);
  }

  updateRecord() {
    return this._createOrUpdate(...arguments);
  }

  urlForDeleteRecord(id, modelName, snapshot) {
    const { backend, type, completeRoleName } = snapshot.record;
    return this._getURL(backend, this._pathForRoleType(type), completeRoleName);
  }

  /* 
    roleAncestry: { path_to_role: string; type: string };
  */
  async query(store, type, query, recordArray, options) {
    const { showPartialError, roleAncestry } = options.adapterOptions || {};
    const { backend } = query;

    if (roleAncestry) {
      return this._querySubdirectory(backend, roleAncestry);
    }

    return this._queryAll(backend, showPartialError);
  }

  // LIST request for all roles (static and dynamic)
  async _queryAll(backend, showPartialError) {
    let roles = [];
    const errors = [];

    for (const roleType of ['static', 'dynamic']) {
      const url = this._getURL(backend, this._pathForRoleType(roleType));
      try {
        const models = await this.ajax(url, 'GET', { data: { list: true } }).then((resp) => {
          return resp.data.keys.map((name) => this._constructRecord({ backend, name, type: roleType }));
        });
        roles = addManyToArray(roles, models);
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
        const errorObject = new AdapterError(errorMessages);
        errorObject.message = 'Error fetching roles:';
        throw errorObject;
      } else if (showPartialError) {
        // if only one request fails, surface the error to the user an info level flash message
        // this may help for permissions errors where a users policy may be incorrect
        this.flashMessages.info(`Error fetching roles from ${errorMessages.join(', ')}`);
      }
    }
    // must return an object in this shape for lazyPaginatedQuery to function
    // changing the responsePath or providing the extractLazyPaginatedData serializer method causes normalizeResponse to return data: [undefined]
    return { data: { keys: sortObjects(roles, 'name') } };
  }

  // LIST request for children of a hierarchical role
  async _querySubdirectory(backend, roleAncestry) {
    // path_to_role is the ancestral path
    const { path_to_role, type: roleType } = roleAncestry;
    const url = `${this._getURL(backend, this._pathForRoleType(roleType))}/${path_to_role}`;
    const roles = await this.ajax(url, 'GET', { data: { list: true } }).then((resp) => {
      return resp.data.keys.map((name) => ({
        ...this._constructRecord({ backend, name, type: roleType }),
        path_to_role, // adds path_to_role attr to ldap/role model
      }));
    });
    return { data: { keys: sortObjects(roles, 'name') } };
  }

  queryRecord(store, type, query) {
    const { backend, name, type: roleType } = query;
    const url = this._getURL(backend, this._pathForRoleType(roleType), name);
    return this.ajax(url, 'GET').then((resp) => ({
      ...resp.data,
      ...this._constructRecord({ backend, name, type: roleType }),
    }));
  }

  fetchCredentials(backend, type, name) {
    const url = this._getURL(backend, this._pathForRoleType(type, true), name);
    return this.ajax(url, 'GET').then((resp) => {
      if (type === 'dynamic') {
        const { lease_id, lease_duration, renewable } = resp;
        return { ...resp.data, lease_id, lease_duration, renewable, type };
      }
      return { ...resp.data, type };
    });
  }
  rotateStaticPassword(backend, name) {
    const url = this._getURL(backend, 'rotate-role', name);
    return this.ajax(url, 'POST');
  }
}
