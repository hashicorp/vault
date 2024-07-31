/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AdapterError from '@ember-data/adapter/error';
import { set } from '@ember/object';
import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { CONFIGURABLE_SECRET_ENGINES } from 'vault/helpers/mountable-secret-engines';

import type Store from '@ember-data/store';

// This route file is reused for all configurable secret engines.
// It's generates models the various models based on the engine type.
// Saving and updating of those models are done within the engine specific components.

const CONFIG_ADAPTERS_PATHS: Record<string, string[]> = {
  aws: ['aws/lease-config', 'aws/root-config'],
  ssh: ['ssh/ca-config'],
};

export default class SecretsBackendConfigurationEdit extends Route {
  @service declare readonly store: Store;

  async model() {
    const { backend } = this.paramsFor('vault.cluster.secrets.backend');
    const record = await this.store.findRecord('secret-engine', backend); // ARG TODO: might be able to do modelFor
    const type = record.type;
    // if the engine type is not configurable, return a 404.
    if (!record || !CONFIGURABLE_SECRET_ENGINES.includes(type)) {
      const error = new AdapterError();
      set(error, 'httpStatus', 404);
      throw error;
    }
    // generate the model based on the engine type.
    const model: Record<string, unknown> = { type, id: backend };
    for (const adapterPath of CONFIG_ADAPTERS_PATHS[type] as string[]) {
      try {
        model[adapterPath] = await this.store.queryRecord(adapterPath, {
          backend,
          type,
        });
      } catch (e: AdapterError) {
        if (e.httpStatus === 404) {
          model[adapterPath] = await this.store.createRecord(adapterPath, {
            backend,
            type,
          });
        } else {
          // ARG TODO figure out error handling, likely in components
          throw e;
        }
      }
    }
    // replace the adapterPath with a useable model name to pass to the components (aws/lease-config -> aws-lease-config)
    for (const key in model) {
      if (key !== 'type' && key !== 'id') {
        const newKey = key.replace(/\//g, '-');
        model[newKey] = model[key] as object | null;
        delete model[key];
      }
    }
    return model;
  }
}
