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
    return this.version.isEnterprise
      ? `Vault v${this.version.version.slice(0, this.version.version.indexOf('+'))}`
      : `Vault v${this.version.version}`;
  }

  get namespaceDisplay() {
    if (this.namespace.inRootNamespace) return 'root';
    const parts = this.namespace.path?.split('/');
    return parts[parts.length - 1];
  }
}
