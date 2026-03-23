/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { service } from '@ember/service';
import type SecretsEngineForm from 'vault/forms/secrets/engine';
import type Router from '@ember/routing/router';
import type { EngineVersionInfo } from 'vault/utils/plugin-catalog-helpers';

export default class VaultClusterSecretsEnableCreateController extends Controller {
  @service declare router: Router;

  declare model: {
    form: SecretsEngineForm;
    availableVersions: EngineVersionInfo[];
  };
}
