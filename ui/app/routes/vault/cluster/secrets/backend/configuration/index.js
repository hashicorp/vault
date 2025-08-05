/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Route from '@ember/routing/route';
import { CONFIGURABLE_SECRET_ENGINES, allEngines } from 'vault/helpers/mountable-secret-engines';
import { reject } from 'rsvp';

/**
 * This route is responsible for fetching all configuration model(s).
 * This includes the mount-configuration model attached to the secret-engine model via a belongsTo relationship.
 * As well as any additional configuration models if the engine is a configurable engine.
 */

export default class SecretsBackendConfigurationRoute extends Route {
  @service store;
  @service version;

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

      return {
        secretEngineModel,
        configModels,
      };
    }
    return { secretEngineModel };
  }

  standardizeConfigModels(configModels) {
    // standardize the configModels to an array so that the component can handle it correctly
    Array.isArray(configModels) ? configModels : (configModels = [configModels]);
    // make sure no items in the array are null or undefined
    return configModels.filter((configModel) => {
      return !!configModel;
    });
  }

  fetchConfig(type, id) {
    // id is the path where the backend is mounted since there's only one config per engine (often this path is referred to just as backend)
    switch (type) {
      case 'aws':
        return this.fetchAwsConfigs(id);
      case 'azure':
        return this.fetchAzureConfig(id);
      case 'gcp':
        return this.fetchGcpConfig(id);
      case 'ssh':
        return this.fetchSshCaConfig(id);
      default:
        return reject({ httpStatus: 404, message: 'not found', path: id });
    }
  }

  async fetchAwsConfigs(id) {
    // AWS has two configuration endpoints root and lease, as well as a separate endpoint for the issuer.
    // return an array of these responses.
    const configArray = [];
    const configRoot = await this.fetchAwsConfig(id, 'aws/root-config');
    const configLease = await this.fetchAwsConfig(id, 'aws/lease-config');
    let issuer = null;
    if (this.version.isEnterprise && configRoot) {
      // issuer is an enterprise only related feature
      // issuer is also a global endpoint that doesn't mean anything in the AWS secret details context if WIF related fields on the rootConfig have not been set.
      const WIF_FIELDS = ['roleArn', 'identityTokenAudience', 'identityTokenTtl'];
      WIF_FIELDS.some((field) => configRoot[field]) ? (issuer = await this.fetchIssuer()) : null;
    }
    configArray.push(configRoot, configLease, issuer);
    return configArray;
  }

  async fetchAwsConfig(id, modelPath) {
    try {
      return await this.store.queryRecord(modelPath, { backend: id });
    } catch (e) {
      if (e.httpStatus === 404) {
        // a 404 error is thrown when the lease config hasn't been set yet.
        return;
      }
      throw e;
    }
  }

  async fetchAzureConfig(id) {
    try {
      const azureModel = await this.store.queryRecord('azure/config', { backend: id });
      let issuer = null;
      if (this.version.isEnterprise) {
        // Issuer is an enterprise only related feature
        // Issuer is also a global endpoint that doesn't mean anything in the Azure secret details context if WIF related fields on the azureConfig have not been set.
        const WIF_FIELDS = ['identityTokenAudience', 'identityTokenTtl'];
        WIF_FIELDS.some((field) => azureModel[field]) ? (issuer = await this.fetchIssuer()) : null;
      }
      const configArray = [];
      if (azureModel.isConfigured) configArray.push(azureModel);
      if (issuer) configArray.push(issuer);
      return configArray;
    } catch (e) {
      if (e.httpStatus === 404) {
        // a 404 error is thrown when Azure's config hasn't been set yet.
        return;
      }
      throw e;
    }
  }

  async fetchGcpConfig(id) {
    try {
      const gcpModel = await this.store.queryRecord('gcp/config', { backend: id });
      let issuer = null;
      if (this.version.isEnterprise) {
        const WIF_FIELDS = ['identityTokenAudience', 'identityTokenTtl', 'serviceAccountEmail'];
        WIF_FIELDS.some((field) => gcpModel[field]) ? (issuer = await this.fetchIssuer()) : null;
      }
      const configArray = [];
      if (gcpModel) configArray.push(gcpModel);
      if (issuer) configArray.push(issuer);
      return configArray;
    } catch (e) {
      if (e.httpStatus === 404) {
        // a 404 error is thrown when GCP's config hasn't been set yet.
        return;
      }
      throw e;
    }
  }

  async fetchIssuer() {
    try {
      return await this.store.queryRecord('identity/oidc/config', {});
    } catch (e) {
      // silently fail if the endpoint is not available or the user doesn't have permission to access it.
      return;
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
    )?.displayName;
    controller.isConfigurable = CONFIGURABLE_SECRET_ENGINES.includes(resolvedModel.secretEngineModel.type);
    controller.modelId = resolvedModel.secretEngineModel.id;
  }
}
