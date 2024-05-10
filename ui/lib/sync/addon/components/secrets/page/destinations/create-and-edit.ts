/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import { service } from '@ember/service';
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

  get header() {
    const { isNew, typeDisplayName, name } = this.args.destination;
    return isNew
      ? {
          title: `Create Destination for ${typeDisplayName}`,
          breadcrumbs: [
            { label: 'Secrets Sync', route: 'secrets.overview' },
            { label: 'Select Destination', route: 'secrets.destinations.create' },
            { label: 'Create Destination' },
          ],
        }
      : {
          title: `Edit ${name}`,
          breadcrumbs: [
            { label: 'Secrets Sync', route: 'secrets.overview' },
            { label: 'Destinations', route: 'secrets.destinations' },
            {
              label: 'Destination',
              route: 'secrets.destinations.destination.secrets',
              model: this.args.destination,
            },
            { label: 'Edit Destination' },
          ],
        };
  }

  @task
  @waitFor
  *save(event: Event) {
    event.preventDefault();
    this.error = '';
    // clear out validation warnings
    this.modelValidations = null;
    const { destination } = this.args;
    const { isValid, state, invalidFormMessage } = destination.validate();

    this.modelValidations = isValid ? null : state;
    this.invalidFormMessage = isValid ? '' : invalidFormMessage;

    if (isValid) {
      try {
        // we only want to save if there are changes
        if (destination.dirtyType as unknown as string) {
          const verb = destination.isNew ? 'created' : 'updated';
          yield destination.save();
          this.flashMessages.success(`Successfully ${verb} the destination ${destination.name}`);
          // when saving a record the server returns all credentials as ******
          // Ember Data observes this as a change, marks the model as dirty and the field will be returned from changedAttributes
          // if the user then attempts to update the record the credential will get overwritten with the masked placeholder value
          // since the record will be fetched from the details route we can safely unload it to avoid the aforementioned issue
          destination.unloadRecord();
          this.store.clearDataset('sync/destination');
        }
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
  updateWarningValidation() {
    if (this.args.destination.isNew) return;
    // check for warnings on change
    const { state } = this.args.destination.validate();
    this.modelValidations = state;
  }

  @action
  cancel() {
    const { isNew } = this.args.destination;
    const method = isNew ? 'unloadRecord' : 'rollbackAttributes';
    this.args.destination[method]();
    this.router.transitionTo(`vault.cluster.sync.secrets.destinations.${isNew ? 'create' : 'destination'}`);
  }
}
