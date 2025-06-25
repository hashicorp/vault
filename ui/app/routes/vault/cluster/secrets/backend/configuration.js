/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Route from '@ember/routing/route';

/**
 * This route is responsible for fetching all configuration data.
 * This includes the configuration attached to the secret engine.
 * In addition, any configuration data associated with configurable engines (aws, azure, gcp, ssh).
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
    const handleError = async (e) => {
      const error = await this.parseApiError(e);
      if (error.httpStatus === 404) {
        // a 404 error is thrown when the lease config hasn't been set yet.
        return {};
      }
      throw error;
    };

    const { data: configRoot } = await this.api.secrets
      .awsReadRootIamCredentialsConfiguration(path)
      .catch(handleError);
    const { data: configLease } = await this.api.secrets.awsReadLeaseConfiguration(path).catch(handleError);

    const WIF_FIELDS = ['roleArn', 'identityTokenAudience', 'identityTokenTtl'];
    const issuer = await this.checkIssuer(configRoot, WIF_FIELDS);

    if (configRoot) {
      return Object.assign(configRoot, configLease, issuer);
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
      const error = await this.parseApiError(e);
      if (error.httpStatus === 404) {
        // a 404 error is thrown when Azure's config hasn't been set yet.
        return;
      }
      throw error;
    }
  }

  async fetchGcpConfig(path) {
    try {
      const { data: gcpConfig } = await this.api.secrets.googleCloudReadConfiguration(path);
      const WIF_FIELDS = ['identityTokenAudience', 'identityTokenTtl', 'serviceAccountEmail'];
      const issuer = await this.checkIssuer(gcpConfig, WIF_FIELDS);

      if (gcpConfig) {
        return Object.assign(gcpConfig, issuer);
      }
    } catch (e) {
      const error = await this.parseApiError(e);
      if (error.httpStatus === 404) {
        // a 404 error is thrown when GCP's config hasn't been set yet.
        return;
      }
      throw error;
    }
  }

  async fetchSshCaConfig(path) {
    try {
      const { data } = await this.api.secrets.sshReadCaConfiguration(path);
      return data;
    } catch (e) {
      const error = await this.parseApiError(e);
      if (
        error.httpStatus === 400 &&
        error.errors?.some((str) => str.includes(`keys haven't been configured yet`))
      ) {
        // When first mounting a SSH engine it throws a 400 error with this specific message.
        // We want to catch this situation and return nothing so that the component can handle it correctly.
        return;
      }
      throw error;
    }
  }

  async checkIssuer(config, fields) {
    // issuer is an enterprise only related feature
    // issuer is also a global endpoint that doesn't mean anything in the secret details context if WIF related fields on the engine's config have not been set.
    if (this.version.isEnterprise && config) {
      const shouldFetchIssuer = fields.some((field) => config[field]);

      if (shouldFetchIssuer) {
        try {
          const { data } = await this.api.identity.oidcReadConfiguration();
          return data;
        } catch (e) {
          // silently fail if the endpoint is not available or the user doesn't have permission to access it.
        }
      }
    }
  }

  async parseApiError(e) {
    const { backend } = this.paramsFor('vault.cluster.secrets.backend');
    const error = await this.api.parseError(e);
    return {
      backend,
      httpStatus: error.status,
      ...error.response,
    };
  }
}
