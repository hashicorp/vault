/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AdapterError from '@ember-data/adapter/error';
import { set } from '@ember/object';
import Route from '@ember/routing/route';
import { service } from '@ember/service';
import errorMessage from 'vault/utils/error-message';
import { action } from '@ember/object';

import type Store from '@ember-data/store';
import type SecretEngineModel from 'vault/models/secret-engine';
import type VersionService from 'vault/services/version';
import engineDisplayData from 'vault/helpers/engines-display-data';

// This route file is reused for all configurable secret engines.
// It generates config models based on the engine type.
// Saving and updating of those models are done within the engine specific components.

const MOUNT_CONFIG_MODEL_NAMES: Record<string, string[]> = {
  aws: ['aws/root-config', 'aws/lease-config'],
  azure: ['azure/config'],
  gcp: ['gcp/config'],
  ssh: ['ssh/ca-config'],
};

export default class SecretsBackendConfigurationEdit extends Route {
  @service declare readonly store: Store;
  @service declare readonly version: VersionService;

  standardizedModelName(type: string, modelName: string): string {
    const modelNames = MOUNT_CONFIG_MODEL_NAMES[type];
    // We check the second item in the modelNames array (if it exists)
    // to determine if there's an additional config model for the given type.
    const additionalModel = modelNames?.[1] ?? null;

    return modelName === additionalModel ? 'additional-config-model' : 'mount-config-model';
  }

  async model() {
    const { backend } = this.paramsFor('vault.cluster.secrets.backend');
    const secretEngineRecord = this.modelFor('vault.cluster.secrets.backend') as SecretEngineModel;
    const type = secretEngineRecord.type;

    // if the engine type is not configurable, return a 404.
    if (!secretEngineRecord || !engineDisplayData(type)?.isConfigurable) {
      const error = new AdapterError();
      set(error, 'httpStatus', 404);
      throw error;
    }
    // generate the model based on the engine type.
    // and pre-set model with type and backend e.g. {type: ssh, id: ssh-123}
    const model: Record<string, unknown> = { type, id: backend };
    for (const modelName of MOUNT_CONFIG_MODEL_NAMES[type] as string[]) {
      // create a key that corresponds with the model order
      // ex: modelName = aws/lease-config, convert to: additional-config-model so that you can pass to component @additionalConfigModel={{this.model.additional-config-model}}
      const standardizedKey = this.standardizedModelName(type, modelName);
      try {
        const configModel = await this.store.queryRecord(modelName, {
          backend,
          type,
        });
        // some of the models return a 200 if they are not configured (ex: azure)
        // so instead of checking a catch or httpStatus, we check if the model is configured based on the getter `isConfigured` on the engine's model
        // if the engine is not configured we update the record to get the default values
        if (!configModel.isConfigured && type === 'azure') {
          model[standardizedKey] = await this.store.createRecord(modelName, {
            backend,
            type,
          });
        } else {
          model[standardizedKey] = configModel;
        }
      } catch (error) {
        const e = error as AdapterError;
        // For most models if the adapter returns a 404, we want to create a new record.
        // The ssh secret engine however returns a 400 if the CA is not configured.
        // For ssh's 400 error, we want to create the CA config model.
        if (
          e.httpStatus === 404 ||
          (type === 'ssh' && e.httpStatus === 400 && errorMessage(e) === `keys haven't been configured yet`)
        ) {
          model[standardizedKey] = await this.store.createRecord(modelName, {
            backend,
            type,
          });
        } else {
          throw e;
        }
      }
    }
    // if the type is a WIF engine and it's enterprise, we also fetch the issuer
    // from a global endpoint which has no associated model/adapter
    if (engineDisplayData(type)?.isWIF && this.version.isEnterprise) {
      try {
        const response = await this.store.queryRecord('identity/oidc/config', {});
        model['identity-oidc-config'] = response;
      } catch (e) {
        // return a property called queryIssuerError and let the component handle it.
        model['identity-oidc-config'] = { queryIssuerError: true };
      }
    }
    return model;
  }

  @action
  willTransition() {
    // catch the transition and refresh model so the route shows the most recent model data.
    this.refresh();
  }
}
