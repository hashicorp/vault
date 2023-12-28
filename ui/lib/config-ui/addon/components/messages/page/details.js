/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';

/**
 * @module Page::MessageDetails
 * Page::MessageDetails components are used to display a message
 * @example
 * ```js
 * <Page::MessageDetails @message={{this.message}}  />
 * ```
 * @param {model} message - message model
 */

export default class MessageDetails extends Component {
  @service store;
  @service router;
  @service flashMessages;

  @task
  *deleteMessage() {
    this.store.clearDataset('config-ui/message');
    yield this.args.message.destroyRecord(this.args.message.id);
    this.router.transitionTo('vault.cluster.config-ui.messages');
    this.flashMessages.success(`Successfully deleted ${this.args.message.title}.`);
  }
}
