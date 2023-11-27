/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { inject as service } from '@ember/service';

/**
 * @module DashboardVaultVersionTitle
 * DashboardVaultVersionTitle component are use to display the vault version title and the badges
 *
 * @example
 * ```js
 * <Dashboard::VaultVersionTitle />
 * ```
 */

export default class DashboardVaultVersionTitle extends Component {
  @service version;
  @service namespace;

  get versionHeader() {
    return `Vault ${this.version.versionDisplay}`;
  }
}
