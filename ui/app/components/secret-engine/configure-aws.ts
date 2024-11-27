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
import type IdentityOidcConfigModel from 'vault/models/identity/oidc/config';
import type Router from '@ember/routing/router';
import type StoreService from 'vault/services/store';
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
  issuerConfig: IdentityOidcConfigModel;
  backendPath: string;
  issuer?: string;
}

export default class ConfigureAwsComponent extends Component<Args> {
  @service declare readonly router: Router;
  @service declare readonly store: StoreService;
  @service declare readonly version: VersionService;
  @service declare readonly flashMessages: FlashMessageService;

  @tracked errorMessageRoot: string | null = null;
  @tracked errorMessageLease: string | null = null;
  @tracked invalidFormAlert: string | null = null;
  @tracked modelValidationsLease: ValidationMap | null = null;
  @tracked accessType = 'iam';
  @tracked saveIssuerWarning = '';

  disableAccessType = false;

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    // the following checks are only relevant to enterprise users and those editing an existing root configuration.
    if (this.version.isCommunity || this.args.rootConfig.isNew) return;
    const { roleArn, identityTokenAudience, identityTokenTtl, accessKey } = this.args.rootConfig;
    // do not include issuer in this check. Issuer is a global endpoint and can be set even if we're not editing wif attributes
    const wifAttributesSet = !!roleArn || !!identityTokenAudience || !!identityTokenTtl;
    const iamAttributesSet = !!accessKey;
    // If any WIF attributes have been set in the rootConfig model, set accessType to 'wif'.
    this.accessType = wifAttributesSet ? 'wif' : 'iam';
    // If there are either WIF or IAM attributes set then disable user's ability to change accessType.
    this.disableAccessType = wifAttributesSet || iamAttributesSet;
  }

  @action continueSubmitForm() {
    // called when the user confirms they are okay with the issuer change
    this.saveIssuerWarning = '';
    this.save.perform();
  }

  // on form submit - validate inputs and check for issuer changes
  submitForm = task(
    waitFor(async (event: Event) => {
      event?.preventDefault();
      this.resetErrors();
      const { leaseConfig, issuerConfig } = this.args;
      // Note: only aws/lease-config model has validations
      const isValid = this.validate(leaseConfig);
      if (!isValid) return;
      if (issuerConfig?.hasDirtyAttributes) {
        // if the issuer has changed show modal with warning that the config will change
        // if the modal is shown, the user has to click confirm to continue save
        this.saveIssuerWarning = `You are updating the global issuer config. This will overwrite Vault's current issuer ${
          issuerConfig.queryIssuerError ? 'if it exists ' : ''
        }and may affect other configurations using this value. Continue?`;
        // exit task until user confirms
        return;
      }
      await this.save.perform();
    })
  );

  save = task(
    waitFor(async () => {
      // when we get here, the models have already been validated so just continue with save
      const { leaseConfig, rootConfig, issuerConfig } = this.args;
      // Check if any of the models' attributes have changed.
      // If no changes to either model, transition and notify user.
      // If changes to either model, save the model(s) that changed and notify user.
      // Note: "backend" dirties model state so explicity ignore it here.
      const leaseAttrChanged = Object.keys(leaseConfig?.changedAttributes()).some(
        (item) => item !== 'backend'
      );
      const rootAttrChanged = Object.keys(rootConfig?.changedAttributes()).some((item) => item !== 'backend');
      const issuerAttrChanged = issuerConfig?.hasDirtyAttributes;
      if (!leaseAttrChanged && !rootAttrChanged && !issuerAttrChanged) {
        this.flashMessages.info('No changes detected.');
        this.transition();
        return;
      }
      // Attempt saves of changed models. If at least one of them succeed, transition
      const rootSaved = rootAttrChanged ? await this.saveRoot() : false;
      const leaseSaved = leaseAttrChanged ? await this.saveLease() : false;
      const issuerSaved = issuerAttrChanged ? await this.updateIssuer() : false;

      if (rootSaved || leaseSaved || issuerSaved) {
        this.transition();
      } else {
        // otherwise there was a failure and we should not transition and exit the function.
        return;
      }
    })
  );

  async updateIssuer(): Promise<boolean> {
    try {
      await this.args.issuerConfig.save();
      this.flashMessages.success('Issuer saved successfully');
      return true;
    } catch (e) {
      this.flashMessages.danger(`Issuer was not saved: ${errorMessage(e, 'Check Vault logs for details.')}`);
      return false;
    }
  }

  async saveRoot(): Promise<boolean> {
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

  async saveLease(): Promise<boolean> {
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
    const { rootConfig } = this.args;
    if (accessType === 'iam') {
      // reset all WIF attributes
      rootConfig.roleArn = rootConfig.identityTokenAudience = rootConfig.identityTokenTtl = undefined;
      // for the issuer return to the globally set value (if there is one) on toggle
      this.args.issuerConfig.rollbackAttributes();
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
