/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';

/**
 * @module HoverCopyButton
 * The `HoverCopyButton` is used on dark backgrounds to show a copy button.
 *
 * @example ```js
 * <HoverCopyButton @copyValue={{stringify this.model.id}} @alwaysShow={{true}} />```
 *
 * @param {string} copyValue - The value to be copied.
 * @param {boolean} [alwaysShow] - Boolean that affects the class.
 */

export default class HoverCopyButton extends Component {
  get alwaysShow() {
    return this.args.alwaysShow || false;
  }
  get copyValue() {
    return this.args.copyValue || false;
  }

  @tracked tooltipText = 'Copy';
}
