/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';
import { dateFormat } from 'core/helpers/date-format';

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

  get getMessages() {
    return this.args.messages.map((message) => {
      let badgeDisplayText = '';

      if (message.active) {
        if (message.endTime) {
          badgeDisplayText = `Active until ${dateFormat([message.endTime, 'MMM d, yyyy hh:mm aaa'], {
            withTimeZone: true,
          })}`;
        } else {
          badgeDisplayText = 'Active';
        }
      } else {
        if (message.isStartTimeAfterToday) {
          badgeDisplayText = `Scheduled: ${dateFormat([message.startTime, 'MMM d, yyyy hh:mm aaa'], {
            withTimeZone: true,
          })}`;
        } else {
          badgeDisplayText = `Inactive:  ${dateFormat([message.startTime, 'MMM d, yyyy hh:mm aaa'], {
            withTimeZone: true,
          })}`;
        }
      }

      message.badgeDisplayText = badgeDisplayText;
      return message;
    });
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
