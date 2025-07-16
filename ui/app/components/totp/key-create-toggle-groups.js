/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

export default class KeyCreateToggleGroupsComponent extends Component {
  @tracked showGroup = null;

  @action
  toggleGroup(group, isOpen) {
    this.showGroup = isOpen ? group : null;
  }
}
