/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import AwsConfigForm from 'vault/forms/secrets/aws-config';
import AzureConfigForm from 'vault/forms/secrets/azure-config';
import GcpConfigForm from 'vault/forms/secrets/gcp-config';
import SshConfigForm from 'vault/forms/secrets/ssh-config';
import engineDisplayData from 'vault/helpers/engines-display-data';

import type SecretsEngineResource from 'vault/resources/secrets/engine';
import type ApiService from 'vault/services/api';
import type VersionService from 'vault/services/version';

type SecretsBackendConfigurationModel = {
  secretsEngine: SecretsEngineResource;
  config: Record<string, unknown>;
};

// This route file is reused for all configurable secret engines.
// It returns config forms based on the engine type.
// Saving and updating of form data is done within the engine specific components.

export default class SecretsBackendConfigurationEdit extends Route {
  @service declare readonly api: ApiService;
  @service declare readonly version: VersionService;

  async model() {
    const { backend } = this.paramsFor('vault.cluster.secrets.backend');
    const { type } = this.modelFor('vault.cluster.secrets.backend') as SecretsEngineResource;
    const { config } = this.modelFor(
      'vault.cluster.secrets.backend.configuration'
    ) as SecretsBackendConfigurationModel;

    const formClass = {
      aws: AwsConfigForm,
      azure: AzureConfigForm,
      gcp: GcpConfigForm,
      ssh: SshConfigForm,
    }[type];

    const defaults = {
      ssh: { generateSigningKey: true, issuer: '' },
    }[type] || { issuer: '' };

    // if the engine type is not configurable or a form class does not exist for the type return a 404.
    if (!engineDisplayData(type)?.isConfigurable || !formClass) {
      throw { httpStatus: 404, backend };
    }

    return {
      type,
      id: backend,
      configForm: new formClass(config || defaults, { isNew: !config }),
    };
  }
}
