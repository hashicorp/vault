/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { assign } from '@ember/polyfills';
import { assert } from '@ember/debug';
import ControlGroupError from 'vault/lib/control-group-error';
import ApplicationAdapter from '../application';
import { allSettled } from 'rsvp';
import { addToArray } from 'vault/helpers/add-to-array';
import { removeFromArray } from 'vault/helpers/remove-from-array';

export default ApplicationAdapter.extend({
  namespace: 'v1',
  pathForType() {
    assert('Generate the url dynamically based on role type', false);
  },

  urlFor(backend, id, type = 'dynamic') {
    let role = 'roles';
    if (type === 'static') {
      role = 'static-roles';
    }
    let url = `${this.buildURL()}/${backend}/${role}`;
    if (id) {
      url = `${this.buildURL()}/${backend}/${role}/${id}`;
    }
    return url;
  },

  staticRoles(backend, id) {
    return this.ajax(this.urlFor(backend, id, 'static'), 'GET', this.optionsForQuery(id)).then((resp) => {
      if (id) {
        return {
          ...resp,
          type: 'static',
          backend,
          id,
        };
      }
      return resp;
    });
  },

  dynamicRoles(backend, id) {
    return this.ajax(this.urlFor(backend, id), 'GET', this.optionsForQuery(id)).then((resp) => {
      if (id) {
        return {
          ...resp,
          type: 'dynamic',
          backend,
          id,
        };
      }
      return resp;
    });
  },

  optionsForQuery(id) {
    const data = {};
    if (!id) {
      data['list'] = true;
    }
    return { data };
  },

  queryRecord(store, type, query) {
    const { backend, id } = query;

    if (query.type === 'static') {
      return this.staticRoles(backend, id);
    } else if (query?.type === 'dynamic') {
      return this.dynamicRoles(backend, id);
    }
    // if role type is not defined, try both
    return allSettled([this.staticRoles(backend, id), this.dynamicRoles(backend, id)]).then(
      ([staticResp, dynamicResp]) => {
        if (staticResp.state === 'rejected' && dynamicResp.state === 'rejected') {
          let reason = staticResp.reason;
          if (dynamicResp.reason instanceof ControlGroupError) {
            throw dynamicResp.reason;
          }
          if (reason?.httpStatus < dynamicResp.reason?.httpStatus) {
            reason = dynamicResp.reason;
          }
          throw reason;
        }
        // Names are distinct across both types of role,
        // so only one request should ever come back with value
        const type = staticResp.value ? 'static' : 'dynamic';
        const successful = staticResp.value || dynamicResp.value;
        const resp = {
          data: {},
          backend,
          id,
          type,
        };

        resp.data = assign({}, successful.data);

        return resp;
      }
    );
  },

  query(store, type, query) {
    const { backend } = query;
    const staticReq = this.staticRoles(backend);
    const dynamicReq = this.dynamicRoles(backend);

    return allSettled([staticReq, dynamicReq]).then(([staticResp, dynamicResp]) => {
      const resp = {
        backend,
        data: { keys: [] },
      };

      if (staticResp.reason && dynamicResp.reason) {
        // both failed, throw error
        throw dynamicResp.reason;
      }
      // at least one request has data
      let staticRoles = [];
      let dynamicRoles = [];

      if (staticResp.value) {
        staticRoles = staticResp.value.data.keys;
      }
      if (dynamicResp.value) {
        dynamicRoles = dynamicResp.value.data.keys;
      }

      resp.data = assign(
        {},
        resp.data,
        { keys: [...staticRoles, ...dynamicRoles] },
        { backend },
        { staticRoles, dynamicRoles }
      );

      return resp;
    });
  },

  async _updateAllowedRoles(store, { role, backend, db, type = 'add' }) {
    const connection = await store.queryRecord('database/connection', { backend, id: db });
    const roles = [...connection.allowed_roles];
    const allowedRoles = type === 'add' ? addToArray([roles, role]) : removeFromArray([roles, role]);
    connection.allowed_roles = allowedRoles;
    return connection.save();
  },

  async createRecord(store, type, snapshot) {
    const serializer = store.serializerFor(type.modelName);
    const data = serializer.serialize(snapshot);
    const roleType = snapshot.attr('type');
    const backend = snapshot.attr('backend');
    const id = snapshot.attr('name');
    const db = snapshot.attr('database');
    try {
      await this._updateAllowedRoles(store, {
        role: id,
        backend,
        db: db[0],
      });
    } catch (e) {
      this.checkError(e);
    }

    return this.ajax(this.urlFor(backend, id, roleType), 'POST', { data }).then(() => {
      // ember data doesn't like 204s if it's not a DELETE
      return {
        data: assign({}, data, { id }),
      };
    });
  },

  async deleteRecord(store, type, snapshot) {
    const roleType = snapshot.attr('type');
    const backend = snapshot.attr('backend');
    const id = snapshot.attr('name');
    const db = snapshot.attr('database');
    try {
      await this._updateAllowedRoles(store, {
        role: id,
        backend,
        db: db[0],
        type: 'remove',
      });
    } catch (e) {
      this.checkError(e);
    }

    return this.ajax(this.urlFor(backend, id, roleType), 'DELETE');
  },

  async updateRecord(store, type, snapshot) {
    const serializer = store.serializerFor(type.modelName);
    const data = serializer.serialize(snapshot);
    const roleType = snapshot.attr('type');
    const backend = snapshot.attr('backend');
    const id = snapshot.attr('name');

    return this.ajax(this.urlFor(backend, id, roleType), 'POST', { data }).then(() => data);
  },

  checkError(e) {
    if (e.httpStatus === 403) {
      // The user does not have the permission to update the connection. This
      // can happen if their permissions are limited to the role. In that case
      // we ignore the error and continue updating the role.
      return;
    }
    throw new Error(`Could not update allowed roles for selected database: ${e.errors.join(', ')}`);
  },
});
