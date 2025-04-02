/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service, { service } from '@ember/service';
import { assert } from '@ember/debug';

import type AdapterError from '@ember-data/adapter/error';
import type CapabilitiesModel from 'vault/vault/models/capabilities';
import type Store from '@ember-data/store';

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
  @service declare readonly store: Store;

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

  async fetchMultiplePaths(paths: string[]): Promise<MultipleCapabilities> {
    // if the request to capabilities-self fails, silently catch
    // all of path capabilities default to "true"
    const resp: CapabilitiesModel[] = await this.request({ paths }).catch(() => []);
    return paths.reduce((obj: MultipleCapabilities, apiPath: string) => {
      // path is the model's primaryKey (id)
      const model = resp.find((m) => m.path === apiPath);
      if (model) {
        const { canCreate, canDelete, canList, canPatch, canRead, canSudo, canUpdate } = model;
        obj[apiPath] = { canCreate, canDelete, canList, canPatch, canRead, canSudo, canUpdate };
      } else {
        // default to true if there is a problem fetching the model
        // since we can rely on the API to gate as a fallback
        obj[apiPath] = {
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
  }

  /*
  this method returns all of the capabilities for a singular path 
  */
  fetchPathCapabilities(path: string): Promise<CapabilitiesModel> | AdapterError {
    try {
      return this.request({ path });
    } catch (error) {
      return error as AdapterError;
    }
  }

  /* 
  internal method for specific capability checks below
  checks the capability model for the passed capability, ie "canRead"
  */
  async _fetchSpecificCapability(
    path: string,
    capability: string
  ): Promise<CapabilitiesModel | AdapterError> {
    try {
      const capabilities = await this.request({ path });
      return capabilities[capability];
    } catch (e) {
      return e as AdapterError;
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
