/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { getOwner } from '@ember/application';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';

/**
 * @module Page::MessagesList
 * Page::MessagesList components are used to display list of messages.
 * @example
 * ```js
 * <Page::MessagesList @messages={{this.messages}}  />
 * ```
 * @param {array} messages - array message objects
 */

export default class MessagesList extends Component {
  @service store;

  get mountPoint() {
    // mountPoint tells transition where to start. In this case, mountPoint will always be vault.cluster.config-ui
    return getOwner(this).mountPoint;
  }

  // callback from HDS pagination to set the queryParams page
  get paginationQueryParams() {
    return (page) => {
      return {
        page,
      };
    };
  }

  @task
  *deleteMessage(message) {
    this.store.clearDataset('config-ui/message');
    yield message.destroyRecord(message.id);
  }
}
