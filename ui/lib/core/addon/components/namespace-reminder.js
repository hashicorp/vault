/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';

export default class NamespaceReminder extends Component {
  @service namespace;

  get showMessage() {
    return !this.namespace.inRootNamespace;
  }

  get mode() {
    return this.args.mode || 'edit';
  }

  get modeVerb() {
    if (!this.mode) {
      return '';
    }
    return this.mode.endsWith('e') ? `${this.mode}d` : `${this.mode}ed`;
  }
}
