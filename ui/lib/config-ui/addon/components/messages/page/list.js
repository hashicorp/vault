/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Ember from 'ember';
import Component from '@glimmer/component';
import { service } from '@ember/service';
import { task, timeout } from 'ember-concurrency';
import { dateFormat } from 'core/helpers/date-format';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { isAfter } from 'date-fns';
import timestamp from 'core/utils/timestamp';

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
  @service customMessages;
  @service flashMessages;
  @service namespace;
  @service('app-router') router;
  @service api;

  @tracked showMaxMessageModal = false;
  @tracked messageToDelete = null;

  isStartTimeAfterToday = (message) => {
    return isAfter(message.start_time, timestamp.now());
  };

  get formattedMessages() {
    return this.args.messages.map((message) => {
      let badgeDisplayText = '';
      let badgeColor = 'neutral';

      if (message.active) {
        if (message.end_time) {
          badgeDisplayText = `Active until ${dateFormat([message.end_time, 'MMM d, yyyy hh:mm aaa'], {
            withTimeZone: true,
          })}`;
        } else {
          badgeDisplayText = 'Active';
        }
        badgeColor = 'success';
      } else {
        if (this.isStartTimeAfterToday(message)) {
          badgeDisplayText = `Scheduled: ${dateFormat([message.start_time, 'MMM d, yyyy hh:mm aaa'], {
            withTimeZone: true,
          })}`;
          badgeColor = 'highlight';
        } else {
          badgeDisplayText = `Inactive:  ${dateFormat([message.start_time, 'MMM d, yyyy hh:mm aaa'], {
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
    return [
      { label: 'Vault', route: 'vault', icon: 'vault', linkExternal: true },
      { label: 'Custom messages' },
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
      // always reset back to page 1 when changing filters
      queryParams: { ...queryParams, page: 1 },
    });
  }

  @task
  *deleteMessage(message) {
    try {
      yield this.api.sys.uiConfigDeleteCustomMessage(message.id);
      this.router.transitionTo('vault.cluster.config-ui.messages');
      this.customMessages.fetchMessages();
      this.flashMessages.success(`Successfully deleted ${message.title}.`);
    } catch (e) {
      const { message } = yield this.api.parseError(e);
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
