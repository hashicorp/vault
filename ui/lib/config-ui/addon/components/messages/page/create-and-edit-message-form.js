/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';
import errorMessage from 'vault/utils/error-message';
import { inject as service } from '@ember/service';
import { format } from 'date-fns';
import { action } from '@ember/object';

import { localDateTimeString } from 'vault/models/config-ui/message';

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
  @service flashMessages;

  @tracked errorBanner = '';
  @tracked modelValidations;
  @tracked invalidFormMessage;
  @tracked endTime;

  constructor() {
    super(...arguments);
    if (this.args.message.endTime) {
      this.endTime = format(new Date(this.args.message.endTime), localDateTimeString);
    }
  }

  willDestroy() {
    super.willDestroy();
    const { isNew } = this.args.message;
    const method = isNew ? 'unloadRecord' : 'rollbackAttributes';
    this.args.message[method]();
  }

  get breadcrumbs() {
    const authenticated =
      this.args.message.authenticated === undefined ? true : this.args.message.authenticated;
    return [
      { label: 'Messages', route: 'messages.index', query: { authenticated } },
      { label: 'Create Message' },
    ];
  }

  @task
  *save(event) {
    event.preventDefault();
    try {
      const { isValid, state, invalidFormMessage } = this.args.message.validate();
      this.modelValidations = isValid ? null : state;
      this.invalidFormAlert = invalidFormMessage;

      if (isValid) {
        const { isNew } = this.args.message;

        // We do these checks here since there could be a scenario where startTime and endTime are strings.
        // The model expects these attrs to be a date object, so we will need to update these attrs to be in
        // date object format.
        if (typeof this.args.message.startTime === 'string')
          this.args.message.startTime = new Date(this.args.message.startTime);
        if (typeof this.args.message.endTime === 'string')
          this.args.message.endTime = new Date(this.args.message.endTime);

        const { id } = yield this.args.message.save();
        this.flashMessages.success(`Successfully ${isNew ? 'created' : 'updated'} the message.`);
        this.router.transitionTo('vault.cluster.config-ui.messages.message.details', id);
      }
    } catch (error) {
      this.errorBanner = errorMessage(error);
      this.invalidFormAlert = 'There was an error submitting this form.';
    }
  }

  @action
  cancel() {
    const { authenticated, isNew } = this.args.message;
    const method = isNew ? 'unloadRecord' : 'rollbackAttributes';
    this.args.message[method]();
    this.router.transitionTo('vault.cluster.config-ui.messages.index', { queryParams: { authenticated } });
  }
}
