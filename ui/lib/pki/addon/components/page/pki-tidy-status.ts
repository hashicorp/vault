/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { action } from '@ember/object';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';

export default class PkiTidyStatusComponent extends Component {
  @tracked tidyOptionsModal = false;
  @tracked confirmCancelTidy = false;

  @action
  cancelTidy() {
    // do the thing
  }
}
