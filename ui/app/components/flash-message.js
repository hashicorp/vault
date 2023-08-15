/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import FlashMessage from 'ember-cli-flash/components/flash-message';

export default class FlashMessageComponent extends FlashMessage {
  // override alertType to get Bulma specific prefix
  //https://github.com/poteto/ember-cli-flash/blob/master/addon/components/flash-message.js#L55
  get alertType() {
    const flashType = this.args.flash.type || '';
    return `is-${flashType}`;
  }
}
