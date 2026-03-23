/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { getEffectiveEngineType } from 'vault/utils/external-plugin-helpers';

import type SecretsEngineResource from 'vault/resources/secrets/engine';
import type { Breadcrumb } from 'vault/vault/app-types';
import type Controller from '@ember/controller';
import type SecretEngineModel from 'vault/models/secret-engine';
import type Transition from '@ember/routing/transition';
import type ApiService from 'vault/services/api';
import type VersionService from 'vault/services/version';

import engineDisplayData from 'vault/helpers/engines-display-data';
import RouterService from '@ember/routing/router-service';

interface RouteModel {
  secretsEngine: SecretEngineModel;
  versions: string[];
  config: Record<string, unknown>;
}

interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: RouteModel;
}

export default class SecretsBackendConfigurationPluginSettingsRoute extends Route {
  @service declare readonly router: RouterService;
  @service declare readonly api: ApiService;
  @service declare readonly version: VersionService;

  async model() {
    const secretsEngine = this.modelFor('vault.cluster.secrets.backend') as SecretsEngineResource;
    const config = await this.fetchConfig(secretsEngine.type, secretsEngine.id);

    return { secretsEngine, config };
  }

  afterModel(resolvedModel: RouteModel) {
    // If there is no config and no custom config route when nav to plugin-settings tab redirect to edit page.
    if (!resolvedModel.config && !engineDisplayData(resolvedModel.secretsEngine.type).configRoute) {
      return this.router.replaceWith('vault.cluster.secrets.backend.configuration.edit');
    } else {
      return;
    }
  }

  setupController(controller: RouteController, resolvedModel: RouteModel, transition: Transition) {
    super.setupController(controller, resolvedModel, transition);
    const { secretsEngine } = resolvedModel;

    const breadcrumbs = [
      { label: 'Vault', route: 'vault.cluster', icon: 'vault' },
      { label: 'Secrets engines', route: 'vault.cluster.secrets' },
      { label: secretsEngine.id, route: 'vault.cluster.secrets.backend.list-root', model: secretsEngine.id },
      { label: 'Configuration' },
    ];

    controller.set('breadcrumbs', breadcrumbs);
  }

  fetchConfig(type: string, id: string) {
    // id is the path where the backend is mounted since there's only one config per engine (often this path is referred to just as backend)
    // Use effective type to handle external plugin mappings
    const effectiveType = getEffectiveEngineType(type);
    switch (effectiveType) {
      case 'aws':
        return this.fetchAwsConfigs(id);
      case 'azure':
        return this.fetchAzureConfig(id);
      case 'gcp':
        return this.fetchGcpConfig(id);
      case 'ssh':
        return this.fetchSshCaConfig(id);
    }
    return; // no config for this engine type
  }

  async fetchAwsConfigs(path: string) {
    // AWS has two configuration endpoints root and lease, as well as a separate endpoint for the issuer.
    const handleError = async (e: Error) => {
      const error = await this.parseApiError(e);
      if (error.httpStatus === 404) {
        // a 404 error is thrown when the lease config hasn't been set yet.
        return {};
      }
      throw error;
    };

    const { data: configRoot } = (await this.api.secrets
      .awsReadRootIamCredentialsConfiguration(path)
      .catch(handleError)) as Record<string, unknown>;
    const { data: configLease } = (await this.api.secrets
      .awsReadLeaseConfiguration(path)
      .catch(handleError)) as Record<string, unknown>;

    const WIF_FIELDS = ['role_arn', 'identity_token_audience', 'identity_token_ttl'];
    const issuer = await this.checkIssuer(configRoot, WIF_FIELDS);

    if (configRoot) {
      return Object.assign(configRoot, configLease, issuer);
    }
    return;
  }

  async fetchAzureConfig(path: string) {
    try {
      const { data: azureConfig } = await this.api.secrets.azureReadConfiguration(path);
      const WIF_FIELDS = ['identity_token_audience', 'identity_token_ttl'];
      const issuer = await this.checkIssuer(azureConfig, WIF_FIELDS);
      // azure config endpoint returns 200 with default values if engine has not been configured yet
      // all values happen to be falsy so we can just check if any are truthy
      const isConfigured =
        azureConfig && Object.values(azureConfig as Record<string, unknown>).some((value) => value);
      if (isConfigured) {
        return Object.assign({}, azureConfig, issuer);
      }
    } catch (e) {
      const error = await this.parseApiError(e as Error);
      if (error.httpStatus === 404) {
        // a 404 error is thrown when Azure's config hasn't been set yet.
        return;
      }
      throw error;
    }
    return;
  }

  async fetchGcpConfig(path: string) {
    try {
      const { data: gcpConfig } = await this.api.secrets.googleCloudReadConfiguration(path);
      const WIF_FIELDS = ['identity_token_audience', 'identity_token_ttl', 'service_account_email'];
      const issuer = await this.checkIssuer(gcpConfig, WIF_FIELDS);

      if (gcpConfig) {
        return Object.assign(gcpConfig, issuer);
      }
    } catch (e) {
      const error = await this.parseApiError(e as Error);
      if (error.httpStatus === 404) {
        // a 404 error is thrown when GCP's config hasn't been set yet.
        return;
      }
      throw error;
    }
    return;
  }

  async fetchSshCaConfig(path: string) {
    try {
      const { data } = await this.api.secrets.sshReadCaConfiguration(path);
      return data;
    } catch (e: any) {
      const error = await this.parseApiError(e);
      if (
        error.httpStatus === 400 &&
        error.errors?.some((str: string) => str.includes(`keys haven't been configured yet`))
      ) {
        // When first mounting a SSH engine it throws a 400 error with this specific message.
        // We want to catch this situation and return nothing so that the component can handle it correctly.
        return;
      }
      throw error;
    }
  }

  async checkIssuer(config: any, fields: string[]) {
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
    return;
  }

  async parseApiError(e: Error) {
    const { backend } = this.paramsFor('vault.cluster.secrets.backend');
    const error = await this.api.parseError(e);
    return {
      backend,
      httpStatus: error.status,
      ...error.response,
    };
  }
}
