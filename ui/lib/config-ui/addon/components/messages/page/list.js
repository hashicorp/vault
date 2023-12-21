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
  @service router;
  @service flashMessages;

  get formattedMessages() {
    return this.args.messages.map((message) => {
      const badgeDisplay = {};

      if (message.active) {
        if (message.endTime) {
          badgeDisplay.text = `Active until ${dateFormat([message.endTime, 'MMM d, yyyy hh:mm aaa'], {
            withTimeZone: true,
          })}`;
          badgeDisplay.color = 'success';
        } else {
          badgeDisplay.text = 'Active';
          badgeDisplay.color = 'success';
        }
      } else {
        if (message.isStartTimeAfterToday) {
          badgeDisplay.text = `Scheduled: ${dateFormat([message.startTime, 'MMM d, yyyy hh:mm aaa'], {
            withTimeZone: true,
          })}`;
          badgeDisplay.color = 'highlight';
        } else {
          badgeDisplay.text = `Inactive: ${dateFormat([message.startTime, 'MMM d, yyyy hh:mm aaa'], {
            withTimeZone: true,
          })}`;
          badgeDisplay.color = 'neutral';
        }
      }

      message.badgeDisplay = badgeDisplay;
      return message;
    });
  }

  get breadcrumbs() {
    const label = this.args.authenticated ? 'After User Logs In' : 'On Login Page';
    return [{ label: 'Messages' }, { label }];
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
    this.router.transitionTo('vault.cluster.config-ui.messages');
    this.flashMessages.success(`Successfully deleted ${message.title}.`);
  }
}
