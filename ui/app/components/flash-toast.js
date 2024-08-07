/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { capitalize } from '@ember/string';
import Component from '@glimmer/component';

/**
 * FlashToast components are used to translate flash messages into toast notifications.
 * Flash object passed should have a `type` and `message` property at minimum.
 */
export default class FlashToastComponent extends Component {
  get color() {
    switch (this.args.flash.type) {
      case 'info':
        return 'highlight';
      case 'danger':
        return 'critical';
      case 'warning':
      case 'success':
        return this.args.flash.type;
      default:
        return 'neutral';
    }
  }

  get title() {
    if (this.args.title) return this.args.title;
    switch (this.args.flash.type) {
      case 'danger':
        return 'Error';
      default:
        return capitalize(this.args.flash.type);
    }
  }
}
