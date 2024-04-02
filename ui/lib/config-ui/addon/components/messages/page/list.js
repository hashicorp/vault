/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { task } from 'ember-concurrency';
import { dateFormat } from 'core/helpers/date-format';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import errorMessage from 'vault/utils/error-message';
import { next } from '@ember/runloop';

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
  @service namespace;
  @service customMessages;

  @tracked showMaxMessageModal = false;
  @tracked messageToDelete = null;

  // This follows the pattern in sync/addon/components/secrets/page/destinations for FilterInput.
  // Currently, FilterInput doesn't do a full page refresh causing it to lose focus.
  // The work around is to verify that a transition from this route was completed and then focus the input.
  constructor(owner, args) {
    super(owner, args);
    this.router.on('routeDidChange', this.focusNameFilter);
  }

  willDestroy() {
    super.willDestroy();
    this.router.off('routeDidChange', this.focusNameFilter);
  }

  focusNameFilter(transition) {
    const route = 'vault.cluster.config-ui.messages.index';
    if (transition?.from?.name === route && transition?.to?.name === route) {
      next(() => document.getElementById('message-filter')?.focus());
    }
  }

  get formattedMessages() {
    return this.args.messages.map((message) => {
      let badgeDisplayText = '';
      let badgeColor = 'neutral';

      if (message.active) {
        if (message.endTime) {
          badgeDisplayText = `Active until ${dateFormat([message.endTime, 'MMM d, yyyy hh:mm aaa'], {
            withTimeZone: true,
          })}`;
        } else {
          badgeDisplayText = 'Active';
        }
        badgeColor = 'success';
      } else {
        if (message.isStartTimeAfterToday) {
          badgeDisplayText = `Scheduled: ${dateFormat([message.startTime, 'MMM d, yyyy hh:mm aaa'], {
            withTimeZone: true,
          })}`;
          badgeColor = 'highlight';
        } else {
          badgeDisplayText = `Inactive:  ${dateFormat([message.startTime, 'MMM d, yyyy hh:mm aaa'], {
            withTimeZone: true,
          })}`;
          badgeColor = 'neutral';
        }
      }

      message.badgeDisplayText = badgeDisplayText;
      message.badgeColor = badgeColor;
      return message;
    });
  }

  get breadcrumbs() {
    const label = this.args.authenticated ? 'After User Logs In' : 'On Login Page';
    return [{ label: 'Messages' }, { label }];
  }

  get statusFilterOptions() {
    return [
      { id: 'active', name: 'active' },
      { id: 'inactive', name: 'inactive' },
    ];
  }

  get typeFilterOptions() {
    return [
      { id: 'modal', name: 'modal' },
      { id: 'banner', name: 'banner' },
    ];
  }

  // callback from HDS pagination to set the queryParams page
  get paginationQueryParams() {
    return (page) => {
      return {
        page,
      };
    };
  }

  transitionToMessagesWithParams(queryParams) {
    this.router.transitionTo('vault.cluster.config-ui.messages', {
      queryParams,
    });
  }

  @task
  *deleteMessage(message) {
    try {
      this.store.clearDataset('config-ui/message');
      yield message.destroyRecord(message.id);
      this.router.transitionTo('vault.cluster.config-ui.messages');
      this.customMessages.fetchMessages(this.namespace.path);
      this.flashMessages.success(`Successfully deleted ${message.title}.`);
    } catch (e) {
      const message = errorMessage(e);
      this.flashMessages.danger(message);
    } finally {
      this.messageToDelete = null;
    }
  }

  @action
  onFilterInputChange(pageFilter) {
    this.transitionToMessagesWithParams({ pageFilter });
  }

  @action
  onFilterChange(filterType, [filterOption]) {
    const param = {};
    param[filterType] = filterOption;
    param.page = 1;
    this.transitionToMessagesWithParams(param);
  }

  @action
  createMessage() {
    if (this.args.messages?.meta && this.args.messages?.meta.total >= 100) {
      this.showMaxMessageModal = true;
      return;
    }

    this.router.transitionTo('vault.cluster.config-ui.messages.create', {
      queryParams: { authenticated: this.args.authenticated },
    });
  }
}
