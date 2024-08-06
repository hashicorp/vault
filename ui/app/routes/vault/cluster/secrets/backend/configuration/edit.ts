/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AdapterError from '@ember-data/adapter/error';
import { set } from '@ember/object';
import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { CONFIGURABLE_SECRET_ENGINES } from 'vault/helpers/mountable-secret-engines';
import errorMessage from 'vault/utils/error-message';
import { action } from '@ember/object';

import type Store from '@ember-data/store';
import type SecretEngineModel from 'vault/models/secret-engine';

// This route file is reused for all configurable secret engines.
// It's generates models the various models based on the engine type.
// Saving and updating of those models are done within the engine specific components.

const CONFIG_ADAPTERS_PATHS: Record<string, string[]> = {
  // aws: ['aws/lease-config', 'aws/root-config'],
  ssh: ['ssh/ca-config'],
};

export default class SecretsBackendConfigurationEdit extends Route {
  @service declare readonly store: Store;

  async model() {
    const { backend } = this.paramsFor('vault.cluster.secrets.backend');
    const record = this.modelFor('vault.cluster.secrets.backend') as { type: SecretEngineModel };
    const type = record.type as string;

    // if the engine type is not configurable, return a 404.
    if (!record || !CONFIGURABLE_SECRET_ENGINES.includes(type)) {
      const error = new AdapterError();
      set(error, 'httpStatus', 404);
      throw error;
    }

    if (type !== 'aws') {
      // generate the model based on the engine type.
      // and pre-set with the type and backend (e.g. type: ssh, id: ssh-123)
      const model: Record<string, unknown> = { type, id: backend };
      for (const adapterPath of CONFIG_ADAPTERS_PATHS[type] as string[]) {
        try {
          model[adapterPath] = await this.store.queryRecord(adapterPath, {
            backend,
            type,
          });
        } catch (e: AdapterError) {
          // For most models if the adapter returns a 404, we want to create a new record.
          // The ssh secret engine however returns a 400 if the CA is not configured.
          // For ssh's 400 error, we want to create the CA config model.
          if (
            e.httpStatus === 404 ||
            (type === 'ssh' && e.httpStatus === 400 && errorMessage(e) === `keys haven't been configured yet`)
          ) {
            model[adapterPath] = await this.store.createRecord(adapterPath, {
              backend,
              type,
            });
          } else {
            throw e;
          }
        }
      }
      // TODO this conditional will be removed when we handle AWS

      // replace the adapterPath with a useable model name to pass to the components (ssh/ca-config -> ssh-ca-config)
      for (const key in model) {
        if (key !== 'type' && key !== 'id') {
          const newKey = key.replace(/\//g, '-');
          model[newKey] = model[key] as object | null;
          delete model[key];
        }
      }
      return model;
    } else {
      return await this.store.findRecord('secret-engine', backend);
    }
  }
  // TODO everything below line will be removed once we handle AWS
  afterModel(model: Record<string, unknown>) {
    const type = model.type;

    if (type === 'aws') {
      return this.store
        .queryRecord('secret-engine', {
          backend: model.id,
          type,
        })
        .then(
          () => model,
          () => model
        );
    }
    return model;
  }

  resetController(controller, isExiting) {
    if (controller.model.type === 'aws') {
      if (isExiting) {
        controller.reset();
      }
    }
  }

  @action
  refreshRoute() {
    this.refresh();
  }
}
