/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Route from '@ember/routing/route';
import { CONFIGURABLE_SECRET_ENGINES, allEngines } from 'vault/helpers/mountable-secret-engines';

/**
 * This route is responsible for fetching all configuration model(s).
 * This includes the mount-configuration model attached to the secret-engine model via a belongsTo relationship.
 * As well as any additional configuration models if the engine is a configurable engine.
 */

export default class SecretsBackendConfigurationRoute extends Route {
  @service api;
  @service version;

  async model() {
    const secretsEngine = this.modelFor('vault.cluster.secrets.backend');
    const { type, id } = secretsEngine;
    return {
      secretsEngine,
      config: await this.fetchConfig(type, id), // fetch config for configurable engines (aws, azure, gcp, ssh)
    };
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
    }
  }

  async fetchAwsConfigs(path) {
    // AWS has two configuration endpoints root and lease, as well as a separate endpoint for the issuer.
    // return an array of these responses.
    try {
      const { data: configRoot } = await this.api.secrets.awsReadRootIamCredentialsConfiguration(path);
      const { data: configLease } = await this.api.secrets.awsReadLeaseConfiguration(path);
      const WIF_FIELDS = ['roleArn', 'identityTokenAudience', 'identityTokenTtl'];
      const issuer = await this.checkIssuer(configRoot, WIF_FIELDS);

      return Object.assign({}, configRoot, configLease, issuer);
    } catch (e) {
      if (e.httpStatus === 404) {
        // a 404 error is thrown when the lease config hasn't been set yet.
        return;
      }
      throw e;
    }
  }

  async fetchAzureConfig(path) {
    try {
      const { data: azureConfig } = await this.api.secrets.azureReadConfiguration(path);
      const WIF_FIELDS = ['identityTokenAudience', 'identityTokenTtl'];
      const issuer = await this.checkIssuer(azureConfig, WIF_FIELDS);
      // azure config endpoint returns 200 with default values if engine has not been configured yet
      // all values happen to be falsy so we can just check if any are truthy
      const isConfigured = Object.values(azureConfig).some((value) => value);
      if (isConfigured) {
        return Object.assign({}, azureConfig, issuer);
      }
    } catch (e) {
      if (e.httpStatus === 404) {
        // a 404 error is thrown when Azure's config hasn't been set yet.
        return;
      }
      throw e;
    }
  }

  async fetchGcpConfig(path) {
    try {
      const { data: gcpConfig } = await this.api.secrets.googleCloudReadConfiguration(path);
      const WIF_FIELDS = ['identityTokenAudience', 'identityTokenTtl', 'serviceAccountEmail'];
      const issuer = await this.checkIssuer(gcpConfig, WIF_FIELDS);

      return Object.assign({}, gcpConfig, issuer);
    } catch (e) {
      if (e.httpStatus === 404) {
        // a 404 error is thrown when GCP's config hasn't been set yet.
        return;
      }
      throw e;
    }
  }

  async fetchSshCaConfig(path) {
    try {
      const { data } = await this.api.secrets.sshReadCaConfiguration(path);
      return data;
    } catch (e) {
      if (e.httpStatus === 400 && e.errors[0] === `keys haven't been configured yet`) {
        // When first mounting a SSH engine it throws a 400 error with this specific message.
        // We want to catch this situation and return nothing so that the component can handle it correctly.
        return;
      }
      throw e;
    }
  }

  async checkIssuer(config, fields) {
    // issuer is an enterprise only related feature
    // issuer is also a global endpoint that doesn't mean anything in the AWS secret details context if WIF related fields on the rootConfig have not been set.
    if (this.version.isEnterprise) {
      const shouldFetchIssuer = fields.some((field) => config[field]);

      if (shouldFetchIssuer) {
        try {
          const { data } = this.api.identity.oidcReadConfiguration();
          return data;
        } catch (e) {
          // silently fail if the endpoint is not available or the user doesn't have permission to access it.
        }
      }
    }
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.typeDisplay = allEngines().find(
      (engine) => engine.type === resolvedModel.secretsEngine.type
    )?.displayName;
    controller.isConfigurable = CONFIGURABLE_SECRET_ENGINES.includes(resolvedModel.secretsEngine.type);
    controller.modelId = resolvedModel.secretsEngine.id;
  }
}
