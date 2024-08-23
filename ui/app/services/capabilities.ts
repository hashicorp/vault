/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service, { service } from '@ember/service';
import { assert } from '@ember/debug';

import type AdapterError from '@ember-data/adapter/error';
import type CapabilitiesModel from 'vault/vault/models/capabilities';
import type StoreService from 'vault/services/store';

interface Query {
  paths?: string[];
  path?: string;
}

export default class CapabilitiesService extends Service {
  @service declare readonly store: StoreService;

  async request(query: Query) {
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

  /*
  this method returns a capabilities model for each path in the array of paths
  */
  async fetchMultiplePaths(paths: string[]): Promise<Array<CapabilitiesModel>> | AdapterError {
    try {
      return await this.request({ paths });
    } catch (e) {
      return e;
    }
  }

  /*
  this method returns all of the capabilities for a singular path 
  */
  async fetchPathCapabilities(path: string): Promise<CapabilitiesModel> | AdapterError {
    try {
      return await this.request({ path });
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

  async canRead(path: string) {
    try {
      return await this._fetchSpecificCapability(path, 'canRead');
    } catch (e) {
      return e;
    }
  }

  async canUpdate(path: string) {
    try {
      return await this._fetchSpecificCapability(path, 'canUpdate');
    } catch (e) {
      return e;
    }
  }
}
