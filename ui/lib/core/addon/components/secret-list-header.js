/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';

/**
 * @module SecretListHeader
 * SecretListHeader component is breadcrumb, title with icon and menu with tabs component. Hello
 *
 * Example is wrapped in back ticks because this component relies on routing and cannot render an isolated sample, so just rendering template sample
 * @example
 * ```
 * <SecretListHeader @isCertTab={{eq this.tab "cert"}} @model={{this.model}} @baseKey={{this.baseKey}} @backendCrumb={{this.backendCrumb}} @filter={{this.filter}} />
 * ```
 *
 * @param {object} model - Model used to pull information about icon and title and backend type for navigation.
 */

export default class SecretListHeader extends Component {
  get isKV() {
    return ['kv', 'generic'].includes(this.args.model.engineType);
  }
}
