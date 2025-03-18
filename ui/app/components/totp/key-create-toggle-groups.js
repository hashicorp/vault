/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

export default class KeyCreateToggleGroupsComponent extends Component {
  @tracked showGroup = null;

  get groups() {
    const { generate } = this.args.model;

    const groups = {
      'TOTP Code Options': ['algorithm', 'digits', 'period'],
    };

    if (generate) {
      groups['Provider Options'] = ['key_size', 'skew', 'exported', 'qr_size'];
    }

    return groups;
  }

  @action
  toggleGroup(group, isOpen) {
    this.showGroup = isOpen ? group : null;
  }
}
