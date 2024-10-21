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
import errorMessage from 'vault/utils/error-message';

import type ConfigModel from 'vault/models/azure/config';
import type IdentityOidcConfigModel from 'vault/models/identity/oidc/config';
import type Router from '@ember/routing/router';
import type StoreService from 'vault/services/store';
import type VersionService from 'vault/services/version';
import type FlashMessageService from 'vault/services/flash-messages';

/**
 * @module ConfigureAzureComponent is used to configure the Azure secret engine
 * A user can configure the endpoint config.
 * For enterprise users, they will see an additional option to config WIF attributes in place of Azure account attributes.
 * The fields for these endpoints are on one form.
 *
 * @example
 * ```js
 * <SecretEngine::ConfigureAzure
    @model={{this.model.azure-config}}
    @backendPath={{this.model.id}}
    @issuerConfig={{this.model.identity-oidc-config}}
    />
 * ```
 *
 * @param {object} model - Azure config model
 * @param {string} backendPath - name of the Azure secret engine, ex: 'azure-123'
 * @param {object} issuerConfigModel - the identity/oidc/config model
 */

interface Args {
  model: ConfigModel;
  issuerConfig: IdentityOidcConfigModel;
  backendPath: string;
}

export default class ConfigureAwsComponent extends Component<Args> {
  @service declare readonly router: Router;
  @service declare readonly store: StoreService;
  @service declare readonly version: VersionService;
  @service declare readonly flashMessages: FlashMessageService;

  @tracked errorMessage: string | null = null;
  @tracked invalidFormAlert: string | null = null;
  @tracked accessType = 'azure';
  @tracked saveIssuerWarning = '';

  disableAccessType = false;

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    // the following checks are only relevant to enterprise users and those editing an existing configuration.
    if (this.version.isCommunity || this.args.model.isNew) return;
    const { identityTokenAudience, identityTokenTtl, subscriptionId } = this.args.model;
    // do not include issuer in this check. Issuer is a global endpoint and can bet set even if we're not editing wif attributes
    const wifAttributesSet = !!identityTokenAudience || !!identityTokenTtl;
    const azureAttributesSet = !!subscriptionId;
    // If any WIF attributes have been set in the model, set accessType to 'wif'.
    this.accessType = wifAttributesSet ? 'wif' : 'azure';
    // If there are either WIF or azure attributes set then disable user's ability to change accessType.
    this.disableAccessType = wifAttributesSet || azureAttributesSet;
  }

  @action continueSubmitForm() {
    // called when the user confirms they are okay with the issuer change
    this.saveIssuerWarning = '';
    this.save.perform();
  }

  // on form submit check for issuer changes then save model and issuer
  submitForm = task(
    waitFor(async (event: Event) => {
      event?.preventDefault();
      this.resetErrors();
      const { issuerConfig } = this.args;
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
      const { model, issuerConfig } = this.args;
      // Check if any of the model attributes have changed.
      // If no changes, transition and notify user.
      // If changes, save the model and notify user.
      // Note: "backend" dirties model state so explicity ignore it here.
      const attrChanged = Object.keys(model?.changedAttributes()).some((item) => item !== 'backend');
      const issuerAttrChanged = issuerConfig?.hasDirtyAttributes;
      if (!attrChanged && !issuerAttrChanged) {
        this.flashMessages.info('No changes detected.');
        this.transition();
        return;
      }

      const modelSaved = attrChanged ? await this.saveModel() : false;
      const issuerSaved = issuerAttrChanged ? await this.updateIssuer() : false;

      if (modelSaved || issuerSaved) {
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

  async saveModel(): Promise<boolean> {
    const { backendPath, model } = this.args;
    try {
      await model.save();
      this.flashMessages.success(`Successfully saved ${backendPath}'s root configuration.`);
      return true;
    } catch (error) {
      this.errorMessage = errorMessage(error);
      this.invalidFormAlert = 'There was an error submitting this form.';
      return false;
    }
  }

  resetErrors() {
    this.flashMessages.clearMessages();
    this.errorMessage = null;
    this.invalidFormAlert = null;
  }

  transition() {
    this.router.transitionTo('vault.cluster.secrets.backend.configuration', this.args.backendPath);
  }

  @action
  onChangeAccessType(accessType: string) {
    this.accessType = accessType;
    const { model } = this.args;
    if (accessType === 'azure') {
      // reset all WIF attributes
      model.identityTokenAudience = model.identityTokenTtl = undefined;
      // for the issuer return to the globally set value (if there is one) on toggle
      this.args.issuerConfig.rollbackAttributes();
    }
    if (accessType === 'wif') {
      // reset all IAM attributes
      model.subscriptionId =
        model.clientId =
        model.tenantId =
        model.clientSecret =
        model.environment =
        model.rootPasswordTtl =
          undefined;
    }
  }

  @action
  onCancel() {
    // clear errors because they're canceling out of the workflow.
    this.resetErrors();
    this.args.model.unloadRecord();
    this.transition();
  }
}
