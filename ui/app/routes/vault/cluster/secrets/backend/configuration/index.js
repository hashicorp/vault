/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Route from '@ember/routing/route';
import { CONFIGURABLE_SECRET_ENGINES } from 'vault/helpers/mountable-secret-engines';
import { allEngines } from 'vault/helpers/mountable-secret-engines';
import { hash } from 'rsvp';
/**
 * This route is responsible for fetching all configuration model(s).
 * This includes the mount-configuration model attached to the secret-engine model via a belongsTo relationship.
 * As well as any additional configuration models if the engine is a configurable engine.
 */

export default class SecretsBackendConfigurationRoute extends Route {
  @service store;

  async model() {
    const secretEngineModel = this.modelFor('vault.cluster.secrets.backend');
    if (secretEngineModel.isV2KV) {
      const canRead = await this.store
        .findRecord('capabilities', `${secretEngineModel.id}/config`)
        .then((response) => response.canRead);
      // only set these config params if they can read the config endpoint.
      if (canRead) {
        // design wants specific default to show that can't be set in the model
        secretEngineModel.casRequired = secretEngineModel.casRequired
          ? secretEngineModel.casRequired
          : 'False';
        secretEngineModel.deleteVersionAfter = secretEngineModel.deleteVersionAfter
          ? secretEngineModel.deleteVersionAfter
          : 'Never delete';
      } else {
        // remove the default values from the model if they don't have read access otherwise it will display the defaults even if they've been set (because they error on returning config data)
        secretEngineModel.set('casRequired', null);
        secretEngineModel.set('deleteVersionAfter', null);
        secretEngineModel.set('maxVersions', null);
      }
    }
    // If the engine is configurable fetch the config model(s) for the engine and return it alongside the model
    if (CONFIGURABLE_SECRET_ENGINES.includes(secretEngineModel.type)) {
      let configModels = await this.fetchConfig(secretEngineModel.type, secretEngineModel.id);
      configModels = this.standardizeConfigModels(configModels);

      return hash({
        secretEngineModel,
        ...configModels,
      });
    }
    return secretEngineModel;
  }

  standardizeConfigModels(configModels) {
    // standardize the configModels to an array so that the component can handle it correctly
    Array.isArray(configModels) ? configModels : (configModels = [configModels]);
    // make sure no items in the array are null or undefined
    configModels.forEach((configModel) => {
      if (!configModel) {
        configModels.splice(configModels.indexOf(configModel), 1);
      }
    });

    return configModels;
  }

  async fetchConfig(type, id) {
    switch (type) {
      case 'aws':
        return await this.fetchAwsConfigs(id);
      case 'ssh':
        return await this.fetchSshCaConfig(id);
    }
  }

  async fetchAwsConfigs(id) {
    // AWS has two configuration endpoints root and lease, return an array of these responses.
    const configArray = [];
    const configRoot = await this.fetchAwsRoot(id);
    const configLease = await this.fetchAwsLease(id);
    configArray.push(configRoot, configLease);
    return configArray;
  }

  async fetchAwsLease(id) {
    try {
      return await this.store.queryRecord('aws/lease-config', { backend: id });
    } catch (e) {
      if (e.httpStatus === 404) {
        // a 404 error is thrown when the lease config hasn't been set yet.
        return;
      }
      throw e;
    }
  }

  async fetchAwsRoot(id) {
    try {
      return await this.store.queryRecord('aws/root-config', { backend: id });
    } catch (e) {
      if (e.httpStatus === 404) {
        // a 404 error is thrown when the root config hasn't been set yet.
        // return and let the component handle the empty config.
        return;
      }
      throw e;
    }
  }

  async fetchSshCaConfig(id) {
    try {
      return await this.store.queryRecord('ssh/ca-config', { backend: id });
    } catch (e) {
      if (e.httpStatus === 400 && e.errors[0] === `keys haven't been configured yet`) {
        // When first mounting a SSH engine it throws a 400 error with this specific message.
        // We want to catch this situation and return nothing so that the component can handle it correctly.
        return;
      }
      throw e;
    }
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.typeDisplay = allEngines().find(
      (engine) => engine.type === resolvedModel.secretEngineModel.type
    ).displayName;
    controller.isConfigurable = CONFIGURABLE_SECRET_ENGINES.includes(resolvedModel.secretEngineModel.type);
    controller.modelId = resolvedModel.secretEngineModel.id;
    // from the resolvedModel remove the secretEngineModel as it's not needed in the configuration details component
    const configModels = { ...resolvedModel };
    delete configModels.secretEngineModel;
    controller.configModels = Object.values(configModels);
  }
}
