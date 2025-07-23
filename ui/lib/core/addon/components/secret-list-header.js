/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { supportedSecretBackends } from 'vault/helpers/supported-secret-backends';
import { CONFIGURATION_ONLY } from 'vault/helpers/mountable-secret-engines';

/**
 * @module SecretListHeader
 * SecretListHeader component is breadcrumb, title with icon and menu with tabs component.
 *
 * Example is wrapped in back ticks because this component relies on routing and cannot render an isolated sample, so just rendering template sample
 * @example
 * ```
 * <SecretListHeader @model={{this.model}} />
 * ```
 *
 * @param {object} model - Model used to pull information about icon and title and backend type for navigation.
 * @param {boolean} [isConfigure=false] - Boolean to determine if the configure tab should be shown.
 */

export default class SecretListHeader extends Component {
  get isKV() {
    return ['kv', 'generic'].includes(this.args.model.engineType);
  }

  get showListTab() {
    // only show the list tab if the engine is not a configuration only engine and the UI supports it
    const { engineType } = this.args.model;
    return supportedSecretBackends().includes(engineType) && !CONFIGURATION_ONLY.includes(engineType);
  }
}
