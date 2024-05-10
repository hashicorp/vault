/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Ember from 'ember';
import Component from '@glimmer/component';
import { service } from '@ember/service';
import { task, timeout } from 'ember-concurrency';
import { dateFormat } from 'core/helpers/date-format';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import errorMessage from 'vault/utils/error-message';

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
      // always reset back to page 1 when changing filters
      queryParams: { ...queryParams, page: 1 },
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

  @task
  *handleSearch(evt) {
    evt.preventDefault();
    const formData = new FormData(evt.target);
    // shows loader to indicate that the search was executed
    yield timeout(Ember.testing ? 0 : 250);
    const params = {};
    for (const key of formData.keys()) {
      const valDefault = key === 'pageFilter' ? '' : null;
      const val = formData.get(key) || valDefault;
      params[key] = val;
    }
    this.transitionToMessagesWithParams(params);
  }

  @action
  resetFilters() {
    this.transitionToMessagesWithParams({ pageFilter: '', status: null, type: null });
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
