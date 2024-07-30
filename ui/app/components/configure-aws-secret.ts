/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { task } from 'ember-concurrency';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import errorMessage from 'vault/utils/error-message';

import type LeaseConfigModel from 'vault/models/aws/lease-config';
import type RootConfigModel from 'vault/models/aws/root-config';
import type { TtlEvent, ValidationMap } from 'vault/app-types';
import type Router from '@ember/routing/router';
import type Store from '@ember-data/store';
import type FlashMessageService from 'vault/services/flash-messages';

/**
 * @module ConfigureAwsSecretComponent
 *
 * @example
 * ```js
 * <ConfigureAwsSecret
    @rootConfig={{this.model.rootConfig}}
    @leaseConfig={{this.model.leaseConfig}}
    @tab={{tab}} />
 * ```
 *
 * @param {object} rootConfig - aws secret engine config/root model
 * @param {object} leaseConfig - aws secret engine config/lease model
 * @param {string} tab - current tab selection // ARG TODO remove
 *
 */

interface Args {
  leaseConfig: LeaseConfigModel;
  rootConfig: RootConfigModel;
  path: string;
  tab?: string;
}

export default class ConfigureAwsSecretComponent extends Component<Args> {
  @service declare readonly router: Router;
  @service declare readonly store: Store;
  @service declare readonly flashMessages: FlashMessageService;

  @tracked modelValidations: ValidationMap | null = null;
  @tracked invalidFormMessage = '';
  @tracked error = '';

  // ARG TODO issues with waitFor revisit
  @task
  *save(event: Event) {
    event.preventDefault();
    const { path, leaseConfig, rootConfig } = this.args;
    //  ARG TODO no validations currently on the models so no model validation checks
    try {
      // try saving root config first
      const action = rootConfig.isNew ? 'created' : 'updated';
      yield rootConfig.save();
      // ARG TODO might need to clear out store.
      this.flashMessages.success(`Successfully saved ${path} root configure.`); // arg todo work on message.
    } catch (error) {
      let message = errorMessage(error);
      debugger;
      this.errorMessage = message;
    }

    // even if root config fails, try saving lease config

    try {
      const action = leaseConfig.isNew ? 'created' : 'updated';
      yield leaseConfig.save();
      // ARG TODO might need to clear out store.
      this.flashMessages.success(`Successfully saved ${path} lease configure.`); // arg todo work on message.
    } catch (error) {
      let message = errorMessage(error);
      debugger;
      this.errorMessage = message;
    }

    // allow transition even if both requests fail.
    this.router.transitionTo('vault.cluster.secrets.backend.configuration', path);
  }
  resetErrors() {
    this.flashMessages.clearMessages();
    this.errorMessage = null;
    this.modelValidations = null;
    this.invalidFormAlert = null;
  }

  @action
  handleTtlChange(name: string, ttlObj: TtlEvent) {
    // lease values cannot be undefined, set to 0 to use default
    const valueToSet = ttlObj.enabled ? ttlObj.goSafeTimeString : 0;
    // ARG TODO need to set the value on the model
    // this.args.model.set(name, valueToSet);
  }
}
