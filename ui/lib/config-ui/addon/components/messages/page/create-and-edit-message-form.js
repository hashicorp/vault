/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { task } from 'ember-concurrency';
import errorMessage from 'vault/utils/error-message';
import { inject as service } from '@ember/service';
import { format, addDays, startOfDay } from 'date-fns';

import { localDateTimeString } from 'vault/models/config-ui/message';

/**
 * @module Page::CreateAndEditMessageForm
 * Page::CreateAndEditMessageForm components are used to display list of messages.
 * @example
 * ```js
 * <Page::CreateAndEditMessageForm @messages={{this.messages}}  />
 * ```
 * @param {array} messages - array message objects
 */

export default class MessagesList extends Component {
  @service router;
  @service flashMessages;

  @tracked startTime = format(addDays(startOfDay(new Date()), 1), localDateTimeString);
  @tracked endTime = '';
  @tracked errorBanner = '';
  @tracked modelValidations;
  @tracked invalidFormMessage;

  get breadcrumbs() {
    const authenticated = this.args.authenticated === undefined ? true : this.args.authenticated;
    return [
      { label: 'Messages', route: 'messages.index', query: { authenticated } },
      { label: 'Create Message' },
    ];
  }

  @action
  updateRadioValue(evt) {
    this.args.messages[evt.target.name] = evt.target.value;
  }

  @action
  updateDateTime(evt) {
    this.args.messages[evt.target.name] = format(new Date(evt.target.value), localDateTimeString);
  }

  @task
  *save(event) {
    event.preventDefault();
    try {
      const { isValid, state, invalidFormMessage } = this.args.messages.validate();
      this.modelValidations = isValid ? null : state;
      this.invalidFormAlert = invalidFormMessage;

      if (isValid) {
        const { isNew } = this.args.messages;
        const { id } = yield this.args.messages.save();
        this.flashMessages.success(`Successfully ${isNew ? 'created' : 'updated'} the message.`);
        this.router.transitionTo('vault.cluster.config-ui.messages.message.details', id);
      }
    } catch (error) {
      this.errorBanner = errorMessage(error);
      this.invalidFormAlert = 'There was an error submitting this form.';
    }
  }
}
