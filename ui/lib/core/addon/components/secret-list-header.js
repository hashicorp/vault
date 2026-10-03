/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import engineDisplayData from 'vault/helpers/engines-display-data';
import { supportedSecretBackends } from 'vault/helpers/supported-secret-backends';
import { getEffectiveEngineType } from 'vault/utils/external-plugin-helpers';

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
  get breadcrumbs() {
    const breadcrumbs = [
      { label: 'Secrets', route: 'vault.cluster.secrets' },
      {
        label: this.args.model.id,
        route: 'vault.cluster.secrets.backend.list-root',
        model: this.args.model.id,
        current: !this.args.isConfigure,
      },
    ];

    if (this.args.isConfigure) {
      breadcrumbs.push([{ label: 'Configure' }]);

      return breadcrumbs;
    }

    return breadcrumbs;
  }

  get effectiveEngineType() {
    return getEffectiveEngineType(this.args.model.engineType);
  }

  get isKV() {
    const effectiveType = getEffectiveEngineType(this.args.model.engineType);
    return ['kv', 'generic'].includes(effectiveType);
  }

  get showListTab() {
    // only show the list tab if the engine is not a configuration only engine and the UI supports it
    const effectiveType = getEffectiveEngineType(this.args.model.engineType);
    return (
      supportedSecretBackends().includes(effectiveType) && !engineDisplayData(effectiveType)?.isOnlyMountable
    );
  }
}
