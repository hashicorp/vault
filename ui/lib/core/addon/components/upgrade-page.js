/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';

export default class UpgradePage extends Component {
  get minimumEdition() {
    return this.args.minimumEdition || 'Vault Enterprise';
  }
  get title() {
    return this.args.title || 'Vault Enterprise';
  }

  get featureName() {
    return this.title === 'Vault Enterprise' ? 'this feature' : this.title;
  }
}
