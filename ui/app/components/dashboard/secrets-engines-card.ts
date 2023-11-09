/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import type SecretEngineModel from 'vault/models/secret-engine';

/**
 * @module DashboardSecretsEnginesCard
 * DashboardSecretsEnginesCard component are used to display 5 secrets engines to the user.
 *
 * @example
 * <DashboardSecretsEnginesCard @secretsEngines={{@model.secretsEngines}} />
 *
 * @param {array} secretsEngines - list of secrets engines
 */

interface Args {
  secretsEngines: SecretEngineModel[];
}

export default class DashboardSecretsEnginesCard extends Component<Args> {
  get filteredSecretsEngines() {
    return this.args.secretsEngines?.filter((secretEngine) => secretEngine.shouldIncludeInList);
  }

  get firstFiveSecretsEngines() {
    return this.filteredSecretsEngines?.slice(0, 5);
  }
}
