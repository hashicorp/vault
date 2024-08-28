/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service, { service } from '@ember/service';
import { assert } from '@ember/debug';

import type AdapterError from '@ember-data/adapter/error';
import type CapabilitiesModel from 'vault/vault/models/capabilities';
import type StoreService from 'vault/services/store';

interface Capabilities {
  canCreate: boolean;
  canDelete: boolean;
  canList: boolean;
  canPatch: boolean;
  canRead: boolean;
  canSudo: boolean;
  canUpdate: boolean;
}

interface MultipleCapabilities {
  [key: string]: Capabilities;
}

export default class CapabilitiesService extends Service {
  @service declare readonly store: StoreService;

  async request(query: { paths?: string[]; path?: string }) {
    if (query?.paths) {
      const { paths } = query;
      return this.store.query('capabilities', { paths });
    }
    if (query?.path) {
      const { path } = query;
      const storeData = await this.store.peekRecord('capabilities', path);
      return storeData ? storeData : this.store.findRecord('capabilities', path);
    }
    return assert('query object must contain "paths" or "path" key', false);
  }

  async fetchMultiplePaths(paths: string[]): MultipleCapabilities | AdapterError {
    try {
      const resp: Array<CapabilitiesModel> = await this.request({ paths });
      return resp.reduce((obj: MultipleCapabilities, model: CapabilitiesModel) => {
        const path = paths.find((p) => model.path === p);
        if (path) {
          const { canCreate, canDelete, canList, canPatch, canRead, canSudo, canUpdate } = model;
          obj[path] = { canCreate, canDelete, canList, canPatch, canRead, canSudo, canUpdate };
        } else {
          // default to true since we can rely on API to gate as a fallback
          obj[model.path] = {
            canCreate: true,
            canDelete: true,
            canList: true,
            canPatch: true,
            canRead: true,
            canSudo: true,
            canUpdate: true,
          };
        }
        return obj;
      }, {});
    } catch (e) {
      return e;
    }
  }

  /*
  this method returns all of the capabilities for a singular path 
  */
  fetchPathCapabilities(path: string): Promise<CapabilitiesModel> | AdapterError {
    try {
      return this.request({ path });
    } catch (error) {
      return error;
    }
  }

  /* 
  internal method for specific capability checks below
  checks the capability model for the passed capability, ie "canRead"
  */
  async _fetchSpecificCapability(
    path: string,
    capability: string
  ): Promise<CapabilitiesModel> | AdapterError {
    try {
      const capabilities = await this.request({ path });
      return capabilities[capability];
    } catch (e) {
      return e;
    }
  }

  canRead(path: string) {
    try {
      return this._fetchSpecificCapability(path, 'canRead');
    } catch (e) {
      return e;
    }
  }

  canUpdate(path: string) {
    try {
      return this._fetchSpecificCapability(path, 'canUpdate');
    } catch (e) {
      return e;
    }
  }

  canPatch(path: string) {
    try {
      return this._fetchSpecificCapability(path, 'canPatch');
    } catch (e) {
      return e;
    }
  }
}
