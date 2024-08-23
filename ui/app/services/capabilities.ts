/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service, { service } from '@ember/service';
import { assert } from '@ember/debug';

import type AdapterError from '@ember-data/adapter/error';
import type ArrayProxy from '@ember/array/proxy';
import type CapabilitiesModel from 'vault/vault/models/capabilities';
import type StoreService from 'vault/services/store';

interface Query {
  paths?: string[];
  path?: string;
}

interface ComputedCapabilities {
  canSudo: string;
  canRead: string;
  canCreate: string;
  canUpdate: string;
  canDelete: string;
  canList: string;
}
export default class CapabilitiesService extends Service {
  @service declare readonly store: StoreService;

  request = async (query: Query) => {
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
  };

  /*
  this method returns all of the capabilities for a particular path or paths
  */
  async fetchAll(query: Query): Promise<CapabilitiesModel | ArrayProxy<CapabilitiesModel>> | AdapterError {
    try {
      return await this.request(query);
    } catch (error) {
      return error;
    }
  }

  /* 
  internal method for specific capability checks below
  checks each capability model (one for each path, if multiple)
  for the passed capability, ie "canRead"
  */
  async _fetchSpecific(
    query: Query,
    capability: string
  ): Promise<CapabilitiesModel | ArrayProxy<CapabilitiesModel>> | AdapterError {
    try {
      const capabilities = await this.request(query);
      if (query?.path) {
        return capabilities[capability];
      }
      if (query?.paths) {
        return capabilities.every((c: CapabilitiesModel) => c[capability as keyof ComputedCapabilities]);
      }
    } catch (e) {
      return e;
    }
  }

  async canRead(query: Query) {
    try {
      return await this._fetchSpecific(query, 'canRead');
    } catch (e) {
      return e;
    }
  }

  async canUpdate(query: Query) {
    try {
      return await this._fetchSpecific(query, 'canUpdate');
    } catch (e) {
      return e;
    }
  }
}
