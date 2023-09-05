/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';

/**
 * @module DashboardSecretsEnginesCard
 * DashboardSecretsEnginesCard component are used to display 5 secrets engines to the user.
 *
 * @example
 * ```js
 * <DashboardSecretsEnginesCard @secretsEngines={{@model.secretsEngines}} />
 * ```
 * @param {array} secretsEngines - list of secrets engines
 */

export default class DashboardSecretsEnginesCard extends Component {
  get filteredSecretsEngines() {
    return this.args.secretsEngines?.filter((secretEngine) => secretEngine.shouldIncludeInList);
  }

  get firstFiveSecretsEngines() {
    return this.filteredSecretsEngines?.slice(0, 5);
  }
}
