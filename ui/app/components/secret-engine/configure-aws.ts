/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { ValidationMap } from 'vault/vault/app-types';
import errorMessage from 'vault/utils/error-message';

import type LeaseConfigModel from 'vault/models/aws/lease-config';
import type RootConfigModel from 'vault/models/aws/root-config';
import type Router from '@ember/routing/router';
import type Store from '@ember-data/store';
import type FlashMessageService from 'vault/services/flash-messages';

/**
 * @module ConfigureAwsComponent is used to configure the AWS secret engine
 * A user can configure the endpoint root/config and/or lease/config. 
 * The fields for these endpoints are on one form.
 *
 * @example
 * ```js
 * <SecretEngine::ConfigureAws
    @rootConfig={{this.model.aws-root-config}}
    @leaseConfig={{this.model.aws-lease-config}} 
    @id={{this.model.id}} 
    />
 * ```
 *
 * @param {object} rootConfig - AWS config/root model
 * @param {object} leaseConfig - AWS config/lease model
 * @param {string} id - name of the AWS secret engine, ex: 'aws-123'
 */

interface Args {
  leaseConfig: LeaseConfigModel;
  rootConfig: RootConfigModel;
  id: string;
}

export default class ConfigureAwsComponent extends Component<Args> {
  @service declare readonly router: Router;
  @service declare readonly store: Store;
  @service declare readonly flashMessages: FlashMessageService;

  @tracked errorMessageRoot: string | null = null;
  @tracked errorMessageLease: string | null = null;
  @tracked invalidFormAlert: string | null = null;
  @tracked modelValidationsLease: ValidationMap | null = null;

  @task
  @waitFor
  *save(event: Event) {
    event.preventDefault();
    this.resetErrors();
    const { id, leaseConfig, rootConfig } = this.args;
    // Note: aws/root-config model does not have any validations
    const isValid = this.validate(leaseConfig);
    if (!isValid) return;
    // Check if any of the models' attributes have changed.
    // If no changes to either model, transition and notify user.
    // If changes to one model, save the model that changed and notify user.
    // Otherwise save both models.
    // Note: "backend" dirties model state so explicity ignore it here.
    const leaseAttrChanged =
      Object.keys(leaseConfig.changedAttributes()).filter((item) => item !== 'backend').length > 0;
    const rootAttrChanged =
      Object.keys(rootConfig.changedAttributes()).filter((item) => item !== 'backend').length > 0;

    if (!leaseAttrChanged && !rootAttrChanged) {
      this.flashMessages.info('No changes detected.');
      this.transition();
    }
    if (rootAttrChanged) {
      try {
        yield rootConfig.save();
        this.flashMessages.success(`Successfully saved ${id}'s root configuration.`);
        this.transition();
      } catch (error) {
        this.errorMessageRoot = errorMessage(error);
        this.invalidFormAlert = 'There was an error submitting this form.';
      }
    }
    if (leaseAttrChanged) {
      try {
        yield leaseConfig.save();
        this.flashMessages.success(`Successfully saved ${id}'s lease configuration.`);
        this.transition();
      } catch (error) {
        // if lease config fails, but there was no error saving rootConfig: notify user of the lease failure with a flash message, save the root config, and transition.
        if (!this.errorMessageRoot) {
          this.flashMessages.danger(`Lease configuration was not saved: ${errorMessage(error)}`, {
            sticky: true,
          });
          this.transition();
        } else {
          this.errorMessageLease = errorMessage(error);
          this.flashMessages.danger(
            `Configuration not saved: ${errorMessage(error)}. ${this.errorMessageRoot}`
          );
        }
      }
    }
  }

  resetErrors() {
    this.flashMessages.clearMessages();
    this.errorMessageRoot = null;
    this.invalidFormAlert = null;
  }

  transition() {
    this.router.transitionTo('vault.cluster.secrets.backend.configuration', this.args.id);
  }

  validate(model: LeaseConfigModel) {
    const { isValid, state, invalidFormMessage } = model.validate();
    this.modelValidationsLease = isValid ? null : state;
    this.invalidFormAlert = isValid ? '' : invalidFormMessage;
    return isValid;
  }

  @action
  onCancel() {
    // clear errors because they're canceling out of the workflow.
    this.resetErrors();
    this.transition();
  }
}
