/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { task, timeout } from 'ember-concurrency';
import { service } from '@ember/service';
import { action } from '@ember/object';
import Ember from 'ember';
import { isAfter } from 'date-fns';
import timestamp from 'core/utils/timestamp';

/**
 * @module Page::CreateAndEditMessageForm
 * Page::CreateAndEditMessageForm components are used to display create and edit message form fields.
 * @example
 * ```js
 * <Page::CreateAndEditMessageForm @message={{this.message}}  />
 * ```
 * @param message - message to pass to form component
 * @param messages - array of all created messages
 * @param breadcrumbs - breadcrumbs to pass to the TabPageHeader component
 */

export default class MessagesList extends Component {
  @service('app-router') router;
  @service pagination;
  @service flashMessages;
  @service customMessages;
  @service api;

  @tracked errorBanner = '';
  @tracked modelValidations;
  @tracked invalidFormMessage;
  @tracked showMessagePreviewModal = false;
  @tracked showMultipleModalsMessage = false;
  @tracked userConfirmation = '';

  get hasSomeActiveModals() {
    const { messages } = this.args;
    return messages?.some((message) => message.type === 'modal' && message.active);
  }

  get hasExpiredModalMessages() {
    const modalMessages = this.args.messages?.filter((message) => message.type === 'modal') || [];
    return modalMessages.every((message) => {
      if (!message.endTime) return false;
      return isAfter(timestamp.now(), new Date(message.endTime));
    });
  }

  validate() {
    const { isValid, state, invalidFormMessage } = this.args.message.toJSON();
    this.modelValidations = isValid ? null : state;
    this.invalidFormAlert = invalidFormMessage;
    return isValid;
  }

  @task
  *save(event) {
    event.preventDefault();
    try {
      this.userConfirmation = '';

      const { message } = this.args;
      const isValid = this.validate();

      if (!this.hasExpiredModalMessages && this.hasSomeActiveModals && message.type === 'modal') {
        this.showMultipleModalsMessage = true;
        const isConfirmed = yield this.getUserConfirmation.perform();
        if (!isConfirmed) return;
      }

      if (isValid) {
        const { data } = message.toJSON();
        let id = data.id;

        if (message.isNew) {
          const response = yield this.api.sys.createCustomMessage(data);
          id = response.data.id;
        } else {
          yield this.api.sys.uiConfigUpdateCustomMessage(id, data);
        }

        this.flashMessages.success(`Successfully saved ${data.title} message.`);
        this.customMessages.fetchMessages();
        this.router.transitionTo('vault.cluster.config-ui.messages.message.details', id);
      }
    } catch (error) {
      const { message } = yield this.api.parseError(error);
      this.errorBanner = message;
      this.invalidFormAlert = 'There was an error submitting this form.';
    }
  }

  @task
  *getUserConfirmation() {
    while (true) {
      if (Ember.testing) {
        return;
      }
      if (this.userConfirmation) {
        return this.userConfirmation === 'confirmed';
      }
      yield timeout(500);
    }
  }

  @action
  displayPreviewModal() {
    const isValid = this.validate();
    if (isValid) {
      this.showMessagePreviewModal = true;
    }
  }

  @action
  updateUserConfirmation(userConfirmation) {
    this.userConfirmation = userConfirmation;
    this.showMultipleModalsMessage = false;
  }

  @action
  cancel() {
    if (this.args.message.isNew) {
      this.router.transitionTo('vault.cluster.config-ui.messages');
    } else {
      this.router.transitionTo('vault.cluster.config-ui.messages.message.details', this.args.message.id);
    }
  }
}
