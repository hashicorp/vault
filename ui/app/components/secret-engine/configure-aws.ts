/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { task } from 'ember-concurrency';
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
 * @param {object} rootConfig - AWS secret engine config/root model
 * @param {object} leaseConfig - AWS secret engine config/lease model
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
  @tracked invalidFormAlert: string | null = null;
  @tracked modelValidationsLease: ValidationMap | null = null;
  @tracked modelValidationsRoot: ValidationMap | null = null;

  @task
  *save(event: Event) {
    event.preventDefault();
    this.resetErrors();
    const { id, leaseConfig, rootConfig } = this.args;
    // root does not have any validations, yet, but it could in the future so designing flow to accommodate to model validations.
    const isValidLease = this.validate(leaseConfig, 'leaseConfig');

    if (!isValidLease) {
      this.flashMessages.danger('Please correct the errors in the form before submitting.');
      return;
    }
    // Check if any of the models attributes have changed.
    // (note: "backend" dirties model state so explicity ignore it here.)
    // If no changes to either model, transition and notify user.
    // If changes to one model, save the model that changed and notify user.
    // Otherwise save both models.
    const leaseAttrChanged =
      Object.keys(leaseConfig.changedAttributes()).filter((item) => item !== 'backend').length > 0;
    const rootAttrChanged =
      Object.keys(rootConfig.changedAttributes()).filter((item) => item !== 'backend').length > 0;

    if (!leaseAttrChanged && !rootAttrChanged) {
      this.flashMessages.danger('No changes detected.');
      this.transition(id);
    }
    if (rootAttrChanged) {
      try {
        yield rootConfig.save();
        this.flashMessages.success(`Successfully saved ${id}'s root configuration.`);
      } catch (error) {
        this.errorMessageRoot = errorMessage(error);
      }
    }
    if (leaseAttrChanged) {
      try {
        yield leaseConfig.save();
        this.flashMessages.success(`Successfully saved ${id}'s lease configuration.`);
      } catch (error) {
        // if lease config fails, notify with a flash message but still allow users to save the root config.
        if (!this.errorMessageRoot) {
          this.flashMessages.danger(`Lease configuration was not saved: ${errorMessage(error)}`, {
            sticky: true,
          });
        } else {
          this.flashMessages.danger(
            `Configuration not saved: ${errorMessage(error)}. ${this.errorMessageRoot}`,
            {
              sticky: true,
            }
          );
        }
      }
    }
    this.transition(id);
  }

  transition(id: string) {
    // prevent transition if there are errors with root configuration
    if (this.errorMessageRoot) {
      this.invalidFormAlert = 'There was an error submitting this form.';
    } else {
      this.unloadModels();
      this.router.transitionTo('vault.cluster.secrets.backend.configuration', id);
    }
  }

  resetErrors() {
    this.flashMessages.clearMessages();
    this.errorMessageRoot = null;
    this.invalidFormAlert = null;
  }

  validate(model: LeaseConfigModel, modelName: string) {
    const { isValid, state, invalidFormMessage } = model.validate();
    // cannot use a tracked object for modelValidations because it will not update the form.
    if (modelName === 'leaseConfig') {
      this.modelValidationsLease = isValid ? null : state;
    }
    this.invalidFormAlert = isValid ? '' : invalidFormMessage;
    return isValid;
  }

  unloadModels() {
    this.args.rootConfig.unloadRecord();
    this.args.leaseConfig.unloadRecord();
  }

  @action
  onCancel() {
    // clear errors because they're canceling out of the workflow.
    this.resetErrors();
    this.transition(this.args.id);
  }
}
