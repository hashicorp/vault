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
import type VersionService from 'vault/services/version';
import type FlashMessageService from 'vault/services/flash-messages';

/**
 * @module ConfigureAwsComponent is used to configure the AWS secret engine
 * A user can configure the endpoint root/config and/or lease/config.
 * For enterprise users, they will see an additional option to config WIF attributes in place of IAM attributes. 
 * The fields for these endpoints are on one form.
 *
 * @example
 * ```js
 * <SecretEngine::ConfigureAws
    @rootConfig={{this.model.aws-root-config}}
    @leaseConfig={{this.model.aws-lease-config}} 
    @backendPath={{this.model.id}} 
    />
 * ```
 *
 * @param {object} rootConfig - AWS config/root model
 * @param {object} leaseConfig - AWS config/lease model
 * @param {string} backendPath - name of the AWS secret engine, ex: 'aws-123'
 */

interface Args {
  leaseConfig: LeaseConfigModel;
  rootConfig: RootConfigModel;
  backendPath: string;
}

export default class ConfigureAwsComponent extends Component<Args> {
  @service declare readonly router: Router;
  @service declare readonly store: Store;
  @service declare readonly version: VersionService;
  @service declare readonly flashMessages: FlashMessageService;

  @tracked errorMessageRoot: string | null = null;
  @tracked errorMessageLease: string | null = null;
  @tracked invalidFormAlert: string | null = null;
  @tracked modelValidationsLease: ValidationMap | null = null;
  @tracked accessType = 'iam';
  @tracked disableAccessType = false;

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    const { roleArn, identityTokenAudience, identityTokenTtl, accessKey } = this.args.rootConfig;
    const wifAttributesSet = !!roleArn || !!identityTokenAudience || !!identityTokenTtl;
    const iamAttributesSet = !!accessKey;
    // If any WIF attributes have been set in the rootConfig model, set accessType to 'wif'.
    this.accessType = wifAttributesSet ? 'wif' : 'iam';
    // If there are either WIF or IAM attributes set then disable user's ability to change accessType.
    this.disableAccessType = wifAttributesSet || iamAttributesSet;
  }

  @task
  @waitFor
  *save(event: Event) {
    event.preventDefault();
    this.resetErrors();
    const { leaseConfig, rootConfig } = this.args;
    // Note: aws/root-config model does not have any validations
    const isValid = this.validate(leaseConfig);
    if (!isValid) return;
    // Check if any of the models' attributes have changed.
    // If no changes to either model, transition and notify user.
    // If changes to either model, save the model(s) that changed and notify user.
    // Note: "backend" dirties model state so explicity ignore it here.

    const leaseAttrChanged =
      Object.keys(leaseConfig.changedAttributes()).filter((item) => item !== 'backend').length > 0;
    const rootAttrChanged =
      Object.keys(rootConfig.changedAttributes()).filter((item) => item !== 'backend').length > 0;

    if (!leaseAttrChanged && !rootAttrChanged) {
      this.flashMessages.info('No changes detected.');
      this.transition();
    }

    const rootSaved = rootAttrChanged ? yield this.saveRoot() : false;
    const leaseSaved = leaseAttrChanged ? yield this.saveLease() : false;

    if (rootSaved || leaseSaved) {
      this.transition();
    } else {
      // otherwise there was a failure and we should not transition and exit the function.
      return;
    }
  }

  async saveRoot() {
    const { backendPath, rootConfig } = this.args;
    try {
      await rootConfig.save();
      this.flashMessages.success(`Successfully saved ${backendPath}'s root configuration.`);
      return true;
    } catch (error) {
      this.errorMessageRoot = errorMessage(error);
      this.invalidFormAlert = 'There was an error submitting this form.';
      return false;
    }
  }

  async saveLease() {
    const { backendPath, leaseConfig } = this.args;
    try {
      await leaseConfig.save();
      this.flashMessages.success(`Successfully saved ${backendPath}'s lease configuration.`);
      return true;
    } catch (error) {
      // if lease config fails, but there was no error saving rootConfig: notify user of the lease failure with a flash message, save the root config, and transition.
      if (!this.errorMessageRoot) {
        this.flashMessages.danger(`Lease configuration was not saved: ${errorMessage(error)}`, {
          sticky: true,
        });
        return true;
      } else {
        this.errorMessageLease = errorMessage(error);
        this.flashMessages.danger(
          `Configuration not saved: ${errorMessage(error)}. ${this.errorMessageRoot}`
        );
        return false;
      }
    }
  }

  resetErrors() {
    this.flashMessages.clearMessages();
    this.errorMessageRoot = null;
    this.invalidFormAlert = null;
  }

  transition() {
    this.router.transitionTo('vault.cluster.secrets.backend.configuration', this.args.backendPath);
  }

  validate(model: LeaseConfigModel) {
    const { isValid, state, invalidFormMessage } = model.validate();
    this.modelValidationsLease = isValid ? null : state;
    this.invalidFormAlert = isValid ? '' : invalidFormMessage;
    return isValid;
  }

  unloadModels() {
    this.args.rootConfig.unloadRecord();
    this.args.leaseConfig.unloadRecord();
  }

  @action
  onChangeAccessType(accessType: string) {
    this.accessType = accessType;
    const { rootConfig } = this.args; // cannot destructure attributes and use them below because they're within if statements and javascript doesn't like that.
    if (accessType === 'iam') {
      // reset all WIF attributes
      rootConfig.roleArn = rootConfig.identityTokenAudience = rootConfig.identityTokenTtl = undefined;
    }
    if (accessType === 'wif') {
      // reset all IAM attributes
      rootConfig.accessKey = rootConfig.secretKey = undefined;
    }
  }

  @action
  onCancel() {
    // clear errors because they're canceling out of the workflow.
    this.resetErrors();
    this.unloadModels();
    this.transition();
  }
}
