/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';

/**
 * @module SecretListHeader
 * SecretListHeader component is by Secret Engine routes to show a routes breadcrumb, title with an icon, and menu with tabs.
 *
 * Example is wrapped in back ticks because this component relies on routing and cannot render an isolated sample, so just rendering template sample
 * @example
 * ```
 * <SecretListHeader @model={{this.model}} @breadcrumbs={{this.breadcrumbs}} />
 * ```
 *
 * @param {object} model - Model used to pull information about icon and title and backend type for navigation.
 * @param {array} breadcrumbs - An array of objects which represent the breadcrumbs for the current path. Breadcrumbs should be set on the controller by the route.
 */

export default class SecretListHeader extends Component {
  get isKV() {
    return ['kv', 'generic'].includes(this.args.model.engineType);
  }
}
