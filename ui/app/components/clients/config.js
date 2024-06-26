/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * @module ClientsConfig
 * ClientsConfig components are used to show and edit the client count config information.
 *
 * @example
 * ```js
 * <Clients::Config @model={{model}} @mode="edit" />
 * ```
 * @param {object} model - model is the DS clients/config model which should be passed in
 * @param {string} [mode=show] - mode is either show or edit. Show results in a table with the config, show has a form.
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';

export default class ConfigComponent extends Component {
  @service router;

  @tracked mode = 'show';
  @tracked modalOpen = false;
  @tracked validations;
  @tracked error = null;

  get infoRows() {
    return [
      {
        label: 'Usage data collection',
        helperText: 'Enable or disable collecting data to track clients.',
        valueKey: 'enabled',
      },
      {
        label: 'Retention period',
        helperText: 'The number of months of activity logs to maintain for client tracking.',
        valueKey: 'retentionMonths',
      },
    ];
  }

  get modalTitle() {
    return `Turn usage tracking ${this.args.model.enabled.toLowerCase()}?`;
  }

  @(task(function* () {
    try {
      yield this.args.model.save();
      this.router.transitionTo('vault.cluster.clients.config');
    } catch (err) {
      this.error = err.message;
      this.modalOpen = false;
    }
  }).drop())
  save;

  @action
  toggleEnabled(event) {
    this.args.model.enabled = event.target.checked ? 'On' : 'Off';
  }

  @action
  onSaveChanges(evt) {
    evt.preventDefault();
    const { isValid, state } = this.args.model.validate();
    const changed = this.args.model.changedAttributes();
    if (!isValid) {
      this.validations = state;
    } else if (changed.enabled) {
      this.modalOpen = true;
    } else {
      this.save.perform();
    }
  }
}
