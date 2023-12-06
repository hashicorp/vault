/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

/**
 * @module Page::CreateAndEditMessageForm
 * Page::CreateAndEditMessageForm components are used to display list of messages.
 * @example
 * ```js
 * <Page::CreateAndEditMessageForm @messages={{this.messages}}  />
 * ```
 * @param {array} messages - array message objects
 */

class MessageState {
  @tracked authenticated = true;
  @tracked type = 'banner';
  @tracked title = '';
  @tracked message = '';
  @tracked linkTitle = '';
  @tracked linkHref = '';
}

export default class MessagesList extends Component {
  @tracked state = new MessageState();

  @action
  updateRadioValue(evt) {
    this.state[evt.target.name] = evt.target.value;
  }
}
