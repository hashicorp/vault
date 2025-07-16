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
import { findDestination } from 'core/helpers/sync-destinations';
import apiMethodResolver from 'sync/utils/api-method-resolver';

import type { ValidationMap } from 'vault/app-types';
import type FlashMessageService from 'vault/services/flash-messages';
import type RouterService from '@ember/routing/router-service';
import type ApiService from 'vault/services/api';
import type { DestinationForm, DestinationType } from 'vault/sync';
import type Owner from '@ember/owner';

interface Args {
  type: DestinationType;
  form: DestinationForm;
}

export default class DestinationsCreateForm extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly api: ApiService;

  @tracked modelValidations: ValidationMap | null = null;
  @tracked invalidFormMessage = '';
  @tracked error = '';

  declare readonly initialCustomTags?: Record<string, string>;

  constructor(owner: Owner, args: Args) {
    super(owner, args);
    // cache initial custom tags value to compare against updates
    // tags that are removed when editing need to be added to the payload
    // cast type here since not all types have customTags
    const { customTags } = args.form.data as unknown as Record<string, unknown>;
    if (customTags) {
      this.initialCustomTags = { ...customTags };
    }
  }

  get header() {
    const { type, form } = this.args;
    const { name: typeDisplayName } = findDestination(type);
    const { name } = form.data;

    return form.isNew
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
              model: { name, type },
            },
            { label: 'Edit Destination' },
          ],
        };
  }

  groupSubtext(group: string, isNew: boolean) {
    const dynamicText = isNew
      ? 'used to authenticate with the destination'
      : 'and the value cannot be read. Enable the input to update';
    switch (group) {
      case 'Advanced configuration':
        return 'Configuration options for the destination.';
      case 'Credentials':
        return `Connection credentials are sensitive information ${dynamicText}.`;
      default:
        return '';
    }
  }

  diffCustomTags(payload: Record<string, unknown>) {
    // if tags were removed we need to add them to the payload
    const { isNew } = this.args.form;
    const { customTags } = payload;
    if (!isNew && customTags && this.initialCustomTags) {
      // compare the new and old keys of customTags object to determine which need to be removed
      const oldKeys = Object.keys(this.initialCustomTags).filter((k) => !Object.keys(customTags).includes(k));
      // add tagsToRemove to the payload if there is a diff
      if (oldKeys.length > 0) {
        payload['tagsToRemove'] = oldKeys;
      }
    }
  }

  save = task(
    waitFor(async (event: Event) => {
      event.preventDefault();
      this.error = '';
      // clear out validation warnings
      this.modelValidations = null;

      const { form, type } = this.args;
      const { name } = form.data;
      const { isValid, state, invalidFormMessage, data } = form.toJSON();

      this.modelValidations = isValid ? null : state;
      this.invalidFormMessage = isValid ? '' : invalidFormMessage;

      if (isValid) {
        try {
          const payload = data as unknown as Record<string, unknown>;
          this.diffCustomTags(payload);
          const method = apiMethodResolver(form.isNew ? 'write' : 'patch', type);
          await this.api.sys[method](name, payload);

          this.router.transitionTo('vault.cluster.sync.secrets.destinations.destination.details', type, name);
        } catch (error) {
          const { message } = await this.api.parseError(
            error,
            'Error saving destination. Please try again or contact support.'
          );
          this.error = message;
        }
      }
    })
  );

  @action
  updateWarningValidation() {
    if (this.args.form.isNew) return;
    // check for warnings on change
    const { state } = this.args.form.toJSON();
    this.modelValidations = state;
  }

  @action
  cancel() {
    const route = this.args.form.isNew ? 'create' : 'destination';
    this.router.transitionTo(`vault.cluster.sync.secrets.destinations.${route}`);
  }
}
