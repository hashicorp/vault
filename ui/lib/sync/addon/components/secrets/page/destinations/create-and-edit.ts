/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import { inject as service } from '@ember/service';
import errorMessage from 'vault/utils/error-message';

import type SyncDestinationModel from 'vault/models/sync/destination';
import { ValidationMap } from 'vault/vault/app-types';
import type FlashMessageService from 'vault/services/flash-messages';
import type RouterService from '@ember/routing/router-service';
import type StoreService from 'vault/services/store';

interface Args {
  destination: SyncDestinationModel;
}

export default class DestinationsCreateForm extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly router: RouterService;
  @service declare readonly store: StoreService;

  @tracked modelValidations: ValidationMap | null = null;
  @tracked invalidFormMessage = '';
  @tracked error = '';

  @task
  @waitFor
  *save(event: Event) {
    event.preventDefault();

    const { destination } = this.args;
    const { isValid, state, invalidFormMessage } = destination.validate();

    this.modelValidations = isValid ? null : state;
    this.invalidFormMessage = isValid ? '' : invalidFormMessage;

    if (isValid) {
      try {
        const verb = destination.isNew ? 'created' : 'updated';
        yield destination.save();
        this.flashMessages.success(`Successfully ${verb} the destination ${destination.name}`);

        this.router.transitionTo(
          'vault.cluster.sync.secrets.destinations.destination.details',
          destination.type,
          destination.name
        );
      } catch (error) {
        this.error = errorMessage(error, 'Error saving destination. Please try again or contact support.');
      }
    }
  }

  @action
  cancel() {
    const method = this.args.destination.isNew ? 'unloadRecord' : 'rollbackAttributes';
    this.args.destination[method]();
    this.router.transitionTo('vault.cluster.sync.secrets.destinations.create');
  }
}
