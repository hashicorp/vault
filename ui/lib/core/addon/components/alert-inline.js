/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { later } from '@ember/runloop';
import { tracked } from '@glimmer/tracking';

/**
 * DEPRECATED - use Hds::Alert @type="compact" instead https://helios.hashicorp.design/components/alert?tab=code#type
 * @module AlertInline
 * `AlertInline` components are used to inform users of important messages.
 *
 * @example
 * ```js
 * <AlertInline @type="danger" @message="{{model.keyId}} is not a valid lease ID"/>
 * ```
 *
 * @param {string} type - The alert type maps to the alert @color
 * @param {string} [message=null] - The message to display within the alert.
 * @param {boolean} [mimicRefresh=false] - If true will display a loading icon when attributes change (e.g. when a form submits and the alert message changes).
 */

export default class AlertInlineComponent extends Component {
  @tracked isRefreshing = false;

  get color() {
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
    if (this.args.mimicRefresh) {
      this.isRefreshing = true;
      later(() => {
        this.isRefreshing = false;
      }, 200);
    }
  }
}
