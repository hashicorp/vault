/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';

/**
 * @module SecretListHeader
 * SecretListHeader component is breadcrumb, title with icon and menu with tabs component.
 *
 * @example
 * <SecretListHeader @model={{hash id="my-secret" icon="person" engineType="kv"}} @@isConfigure={{true}} />
 *
 * @param {object} model - Model used to pull information about icon and title and backend type for navigation.
 */

export default class SecretListHeader extends Component {
  get isKV() {
    return ['kv', 'generic'].includes(this.args.model.engineType);
  }
}
