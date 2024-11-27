/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AdapterError from '@ember-data/adapter/error';
import { set } from '@ember/object';
import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { CONFIGURABLE_SECRET_ENGINES, WIF_ENGINES } from 'vault/helpers/mountable-secret-engines';
import errorMessage from 'vault/utils/error-message';
import { action } from '@ember/object';

import type Store from '@ember-data/store';
import type SecretEngineModel from 'vault/models/secret-engine';
import type VersionService from 'vault/services/version';

// This route file is reused for all configurable secret engines.
// It generates config models based on the engine type.
// Saving and updating of those models are done within the engine specific components.

const CONFIG_ADAPTERS_PATHS: Record<string, string[]> = {
  aws: ['aws/lease-config', 'aws/root-config'],
  ssh: ['ssh/ca-config'],
};

export default class SecretsBackendConfigurationEdit extends Route {
  @service declare readonly store: Store;
  @service declare readonly version: VersionService;

  async model() {
    const { backend } = this.paramsFor('vault.cluster.secrets.backend');
    const secretEngineRecord = this.modelFor('vault.cluster.secrets.backend') as SecretEngineModel;
    const type = secretEngineRecord.type;

    // if the engine type is not configurable, return a 404.
    if (!secretEngineRecord || !CONFIGURABLE_SECRET_ENGINES.includes(type)) {
      const error = new AdapterError();
      set(error, 'httpStatus', 404);
      throw error;
    }
    // generate the model based on the engine type.
    // and pre-set model with type and backend e.g. {type: ssh, id: ssh-123}
    const model: Record<string, unknown> = { type, id: backend };
    for (const adapterPath of CONFIG_ADAPTERS_PATHS[type] as string[]) {
      // convert the adapterPath with a name that can be passed to the components
      // ex: adapterPath = ssh/ca-config, convert to: ssh-ca-config so that you can pass to component @model={{this.model.ssh-ca-config}}
      const standardizedKey = adapterPath.replace(/\//g, '-');
      try {
        model[standardizedKey] = await this.store.queryRecord(adapterPath, {
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
          model[standardizedKey] = await this.store.createRecord(adapterPath, {
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
    if (WIF_ENGINES.includes(type) && this.version.isEnterprise) {
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
