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

import type ConfigModel from 'vault/models/gcp/config';
import type IdentityOidcConfigModel from 'vault/models/identity/oidc/config';
import type Router from '@ember/routing/router';
import type StoreService from 'vault/services/store';
import type VersionService from 'vault/services/version';
import type FlashMessageService from 'vault/services/flash-messages';

/**
 * @module SecretEngineConfigureGcp component is used to configure the Gcp secret engine
 * For enterprise users, they will see an additional option to config WIF attributes in place of GCP credentials.
 * If the user is configuring WIF attributes they will also have the option to update the global issuer config, which is a separate endpoint named identity/oidc/config.
 * @example
 * <SecretEngine::ConfigureGcp
    @model={{this.model.gcp-config}}
    @backendPath={{this.model.id}}
    @issuerConfig={{this.model.identity-oidc-config}}
    />
 * 
 * @param {object} model - GCP config model
 * @param {string} backendPath - name of the Gcp secret engine, ex: 'gcp-123'
 * @param {object} issuerConfigModel - the identity/oidc/config model
 */

interface Args {
  model: ConfigModel;
  issuerConfig: IdentityOidcConfigModel;
  backendPath: string;
}

export default class ConfigureGcpComponent extends Component<Args> {
  @service declare readonly router: Router;
  @service declare readonly store: StoreService;
  @service declare readonly version: VersionService;
  @service declare readonly flashMessages: FlashMessageService;

  @tracked accessType = 'gcp';
  @tracked errorMessage = '';
  @tracked invalidFormAlert = '';
  @tracked saveIssuerWarning = '';

  disableAccessType = false;

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    // the following checks are only relevant to existing enterprise configurations
    if (this.version.isCommunity && this.args.model.isNew) return;
    const { isWifPluginConfigured } = this.args.model;
    this.accessType = isWifPluginConfigured ? 'wif' : 'gcp';
    this.disableAccessType = isWifPluginConfigured;
  }

  get modelAttrChanged() {
    // "backend" dirties model state so explicity ignore it here
    return Object.keys(this.args.model?.changedAttributes()).some((item) => item !== 'backend');
  }

  get issuerAttrChanged() {
    return this.args.issuerConfig?.hasDirtyAttributes;
  }

  @action continueSubmitForm() {
    this.saveIssuerWarning = '';
    this.save.perform();
  }

  // check if the issuer has been changed to show issuer modal
  // continue saving the configuration
  submitForm = task(
    waitFor(async (event: Event) => {
      event?.preventDefault();
      this.resetErrors();

      if (this.issuerAttrChanged) {
        // if the issuer has changed show modal with warning that the config will change
        // if the modal is shown, the user has to click confirm to continue saving
        this.saveIssuerWarning = `You are updating the global issuer config. This will overwrite Vault's current issuer ${
          this.args.issuerConfig.queryIssuerError ? 'if it exists ' : ''
        }and may affect other configurations using this value. Continue?`;
        // exit task until user confirms
        return;
      }
      await this.save.perform();
    })
  );

  save = task(
    waitFor(async () => {
      const modelAttrChanged = this.modelAttrChanged;
      const issuerAttrChanged = this.issuerAttrChanged;
      // check if any of the model or issue attributes have changed
      // if no changes, transition and notify user
      if (!modelAttrChanged && !issuerAttrChanged) {
        this.flashMessages.info('No changes detected.');
        this.transition();
        return;
      }

      const modelSaved = modelAttrChanged ? await this.saveModel() : false;
      const issuerSaved = issuerAttrChanged ? await this.updateIssuer() : false;

      if (modelSaved || (!modelAttrChanged && issuerSaved)) {
        // transition if the model was saved successfully
        // we only prevent a transition if the model is edited and fails saving
        this.transition();
      } else {
        // otherwise there was a failure and we should not transition and exit the function
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
      // remove issuer from the config model if it was not saved
      this.args.issuerConfig.rollbackAttributes();
      return false;
    }
  }

  async saveModel(): Promise<boolean> {
    const { backendPath, model } = this.args;
    try {
      await model.save();
      this.flashMessages.success(`Successfully saved ${backendPath}'s configuration.`);
      return true;
    } catch (error) {
      this.errorMessage = errorMessage(error);
      this.invalidFormAlert = 'There was an error submitting this form.';
      return false;
    }
  }

  resetErrors() {
    this.flashMessages.clearMessages();
    this.errorMessage = '';
    this.invalidFormAlert = '';
  }

  transition() {
    this.router.transitionTo('vault.cluster.secrets.backend.configuration', this.args.backendPath);
  }

  @action
  onChangeAccessType(accessType: string) {
    this.accessType = accessType;
    const { model } = this.args;
    if (accessType === 'gcp') {
      // reset all WIF attributes
      model.identityTokenAudience = model.identityTokenTtl = undefined;
      // return the issuer to the globally set value (if there is one) on toggle
      this.args.issuerConfig.rollbackAttributes();
    }
    if (accessType === 'wif') {
      // reset all Gcp attributes
      model.credentials = undefined;
    }
  }

  @action
  onCancel() {
    this.resetErrors();
    this.args.model.unloadRecord();
    this.transition();
  }
}
