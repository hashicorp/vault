/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { task } from 'ember-concurrency';
import errorMessage from 'vault/utils/error-message';

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
  @tracked endTime = '';
}

export default class MessagesList extends Component {
  @tracked state = new MessageState();

  @action
  updateRadioValue(evt) {
    this.state[evt.target.name] = evt.target.value;
  }

  @task
  *save(event) {
    event.preventDefault();
    try {
      const { isNew } = this.args.messages;
      yield this.args.messages.save();
      this.flashMessages.success(`Successfully ${isNew ? 'created' : 'updated'} the message.`);
    } catch (error) {
      this.errorBanner = errorMessage(error);
      this.invalidFormAlert = 'There was an error submitting this form.';
    }
  }
}
