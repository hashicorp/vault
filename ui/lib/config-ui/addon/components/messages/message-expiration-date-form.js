/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';

/**
 * @module Messages::MessageExpirationDateForm
 * Messages::MessageExpirationDateForm components are used to display list of messages.
 * @example
 * ```js
 * <Messages::MessageExpirationDateForm @message={{this.message}} @attr={{attr}} />
 * ```
 * @param {array} messages - array message objects
 */

export default class MessageExpirationDateForm extends Component {
  @tracked groupValue = 'never';
  @tracked formDateTime = '';

  constructor() {
    super(...arguments);
    if (this.args.message.endTime) {
      this.groupValue = 'specificDate';
    }
  }
}
