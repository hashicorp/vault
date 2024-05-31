/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { task, timeout } from 'ember-concurrency';
import errorMessage from 'vault/utils/error-message';
import { service } from '@ember/service';
import { action } from '@ember/object';
import Ember from 'ember';
import { isAfter } from 'date-fns';

/**
 * @module Page::CreateAndEditMessageForm
 * Page::CreateAndEditMessageForm components are used to display create and edit message form fields.
 * @example
 * ```js
 * <Page::CreateAndEditMessageForm @message={{this.message}}  />
 * ```
 * @param {model} message - message model to pass to form components
 */

export default class MessagesList extends Component {
  @service router;
  @service store;
  @service flashMessages;
  @service customMessages;
  @service namespace;

  @tracked errorBanner = '';
  @tracked modelValidations;
  @tracked invalidFormMessage;
  @tracked showMessagePreviewModal = false;
  @tracked showMultipleModalsMessage = false;
  @tracked userConfirmation = '';

  willDestroy() {
    const noTeardown = this.store && !this.store.isDestroying;
    const { model } = this;
    if (noTeardown && model && model.isDirty && !model.isDestroyed && !model.isDestroying) {
      model.rollbackAttributes();
    }
    super.willDestroy();
  }

  validate() {
    const { isValid, state, invalidFormMessage } = this.args.message.validate();
    this.modelValidations = isValid ? null : state;
    this.invalidFormAlert = invalidFormMessage;
    return isValid;
  }

  @task
  *save(event) {
    event.preventDefault();
    try {
      this.userConfirmation = '';

      const isValid = this.validate();
      const modalMessages = this.args.messages?.filter((message) => message.type === 'modal') || [];
      const hasExpiredModalMessages = modalMessages.every((message) => {
        if (!message.endTime) return false;
        return isAfter(new Date(), new Date(message.endTime));
      });

      if (!hasExpiredModalMessages && this.args.hasSomeActiveModals && this.args.message.type === 'modal') {
        this.showMultipleModalsMessage = true;
        const isConfirmed = yield this.getUserConfirmation.perform();
        if (!isConfirmed) return;
      }

      if (isValid) {
        const { isNew } = this.args.message;
        const { id, title } = yield this.args.message.save();
        this.flashMessages.success(`Successfully ${isNew ? 'created' : 'updated'} ${title} message.`);
        this.store.clearDataset('config-ui/message');
        this.customMessages.fetchMessages(this.namespace.path);
        this.router.transitionTo('vault.cluster.config-ui.messages.message.details', id);
      }
    } catch (error) {
      this.errorBanner = errorMessage(error);
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
