/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { later } from '@ember/runloop';
import { tracked } from '@glimmer/tracking';
import { assert } from '@ember/debug';

/**
 * @module AlertInline
 * * Use Hds::Alert @type="compact" for displaying alert messages.
 * This component renders a compact Hds::Alert that displays a loading icon if the
 * @message arg changes and then re-renders the updated @message text.
 * (Example: submitting a form and displaying the number of errors because on re-submit the number may change)
 *
 * @example
 * ```
 * <AlertInline @type="danger" @message="There are 2 errors with this form."/>
 * ```
 *
 * @deprecated {string} type - color getter maps type to the Hds::Alert @color
 * @param {string} color - Styles alert color and icon, can be one of: critical, warning, success, highlight, neutral
 * @param {string} message - The message to display within the alert.
 */

export default class AlertInlineComponent extends Component {
  @tracked isRefreshing = false;

  constructor() {
    super(...arguments);
    assert('@type arg is deprecated, pass @color="critical" instead', this.args.type !== 'critical');
    if (this.args.color) {
      const possibleColors = ['critical', 'warning', 'success', 'highlight', 'neutral'];
      assert(
        `${this.args.color} is not a valid color. @color must be one of: ${possibleColors.join(', ')}`,
        possibleColors.includes(this.args.color)
      );
    }
  }

  get color() {
    if (this.args.color) return this.args.color;
    // @type arg is deprecated, this is for backward compatibility of old implementation
    switch (this.args.type) {
      case 'danger':
        return 'critical';
      case 'success':
        return 'success';
      case 'warning':
        return 'warning';
      case 'info':
        return 'highlight';
      default:
        return 'neutral';
    }
  }

  @action
  refresh() {
    this.isRefreshing = true;
    later(() => {
      this.isRefreshing = false;
    }, 200);
  }
}
