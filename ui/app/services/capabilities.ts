/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service, { service } from '@ember/service';

import type StoreService from 'vault/services/store';

export default class CapabilitiesService extends Service {
  @service declare readonly store: StoreService;

  request = (apiPath: string) => {
    return this.store.findRecord('capabilities', apiPath);
  };

  async fetchAll(apiPath: string) {
    try {
      return await this.request(apiPath);
    } catch (e) {
      return e;
    }
  }

  async fetchSpecific(apiPath: string, capability: string) {
    try {
      const capabilities = await this.request(apiPath);
      return capabilities[capability];
    } catch (e) {
      return e;
    }
  }

  async canRead(apiPath: string) {
    try {
      return await this.fetchSpecific(apiPath, 'canRead');
    } catch (e) {
      return e;
    }
  }

  async canUpdate(apiPath: string) {
    try {
      return await this.fetchSpecific(apiPath, 'canUpdate');
    } catch (e) {
      return e;
    }
  }
}
