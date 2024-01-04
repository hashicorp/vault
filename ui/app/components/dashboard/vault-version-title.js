/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { inject as service } from '@ember/service';

/**
 * @module DashboardVaultVersionTitle
 * DashboardVaultVersionTitle component are use to display the vault version title and the badges
 *
 * @example
 * ```js
 * <Dashboard::VaultVersionTitle @version={{this.versionSvc}} />
 * ```
 */

export default class DashboardVaultVersionTitle extends Component {
  @service namespace;
}
