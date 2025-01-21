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
import { assert } from '@ember/debug';
import errorMessage from 'vault/utils/error-message';

import type ConfigModel from 'vault/models/secret-engine/config';
import type SecondConfigModel from 'vault/models/secret-engine/second-config';
import type IdentityOidcConfigModel from 'vault/models/identity/oidc/config';
import type Router from '@ember/routing/router';
import type StoreService from 'vault/services/store';
import type VersionService from 'vault/services/version';
import type FlashMessageService from 'vault/services/flash-messages';

/**
 * @module SecretEngineConfigureWif component is used to configure secret engines that allow the WIF configuration.
 * The ability to configure WIF fields is an enterprise only feature.
 * If the user is configuring WIF attributes they will also have the option to update the global issuer config, which is a separate endpoint named identity/oidc/config.
 * If a user is on OSS, the account configuration fields will display with no ability to select or see wif fields.
 * 
 * @example
 * <SecretEngine::ConfigureWif
    @backendPath={{this.model.id}}
    @displayName="AWS"
    @type="aws"
    @model={{this.model.root-config}}
    @secondModel={{this.model.lease-config}}
    @issuerConfig={{this.model.identity-oidc-config}}
    />
 *
 * @param {string} backendPath - name of the secret engine, ex: 'azure-123'.
 * @param {string} displayName - used for flash messages, subText and labels. ex: 'Azure'.
 * @param {string} type - the type of the engine, ex: 'azure'.
 * @param {object} model - the config model for the engine.
 * @param {object} [secondModel] - for engines with two config models. Currently, only used by aws
 * @param {object} [issuerConfig] - the identity/oidc/config model. Will be passed in if user has an enterprise license.
 */

interface Args {
  backendPath: string;
  displayName: string;
  type: string;
  model: ConfigModel;
  secondModel: SecondConfigModel;
  issuerConfig: IdentityOidcConfigModel;
}

export default class ConfigureWif extends Component<Args> {
  @service declare readonly router: Router;
  @service declare readonly store: StoreService;
  @service declare readonly version: VersionService;
  @service declare readonly flashMessages: FlashMessageService;

  @tracked accessType = 'account'; // for community users they will not be able to change this. for enterprise users, they will have the option to select "wif".
  @tracked errorMessage = '';
  @tracked invalidFormAlert = '';
  @tracked saveIssuerWarning = '';
  @tracked modelValidations: ValidationMap | null = null;

  disableAccessType = false;

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    // the following checks are only relevant to existing enterprise configurations
    if (this.version.isCommunity && this.args.model.isNew) return;
    const { isWifPluginConfigured, isAccountPluginConfigured } = this.args.model;
    assert(
      `'isWifPluginConfigured' is required to be defined on the config model. Must return a boolean.`,
      isWifPluginConfigured !== undefined
    );
    this.accessType = isWifPluginConfigured ? 'wif' : 'account';
    // if wif or account only attributes are defined, disable the user's ability to change the access type
    this.disableAccessType = isWifPluginConfigured || isAccountPluginConfigured;
  }

  get modelAttrChanged() {
    // "backend" dirties model state so explicity ignore it here
    return Object.keys(this.args.model?.changedAttributes()).some((item) => item !== 'backend');
  }

  get issuerAttrChanged() {
    return this.args.issuerConfig?.hasDirtyAttributes;
  }

  get secondModelAttrChanged() {
    const { secondModel } = this.args;
    // required to check for model first otherwise Object.keys will have nothing to iterate over and fails
    if (!secondModel) return;
    return Object.keys(secondModel.changedAttributes()).some((item) => item !== 'backend');
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
      // currently we only check validations on the second model (for AWS lease config).
      if (this.args.secondModel && !this.validate(this.args.secondModel)) {
        return;
      }
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
      const secondModelAttrChanged = this.secondModelAttrChanged;
      // check if any of the model(s) or issuer attributes have changed
      // if no changes, transition and notify user
      if (!modelAttrChanged && !issuerAttrChanged && !secondModelAttrChanged) {
        this.flashMessages.info('No changes detected.');
        this.transition();
        return;
      }

      const modelSaved = modelAttrChanged ? await this.saveModel() : false;
      const issuerSaved = issuerAttrChanged ? await this.updateIssuer() : false;

      if (modelSaved || (!modelAttrChanged && issuerSaved)) {
        // if there is a second model, attempt to save it. if saving fails, we transition and the failure will surface as a sticky flash message on the configuration details page.
        if (secondModelAttrChanged) {
          await this.saveSecondModel();
        }
        // we only prevent a transition if the first model or issuer are edited and fail when saving
        this.transition();
      } else {
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

  async saveSecondModel(): Promise<boolean> {
    const { backendPath, secondModel, type } = this.args;
    const secondModelName = type === 'aws' ? 'Lease configuration' : 'additional configuration';
    try {
      await secondModel.save();
      this.flashMessages.success(`Successfully saved ${backendPath}'s ${secondModelName}.`);
      return true;
    } catch (error) {
      this.errorMessage = errorMessage(error);
      // we transition even if the second model fails. surface a sticky flash message so the user can see it on the next view.
      this.flashMessages.danger(`${secondModelName} was not saved: ${this.errorMessage}`, {
        sticky: true,
      });
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

  validate(model: SecondConfigModel) {
    const { isValid, state, invalidFormMessage } = model.validate();
    this.modelValidations = isValid ? null : state;
    this.invalidFormAlert = isValid ? '' : invalidFormMessage;
    return isValid;
  }

  @action
  onChangeAccessType(accessType: string) {
    this.accessType = accessType;
    const { model, type } = this.args;
    if (accessType === 'account') {
      // reset all "wif" attributes that are mutually exclusive with "account" attributes
      // these attributes are the same for each engine
      model.identityTokenAudience = model.identityTokenTtl = undefined;
      // return the issuer to the globally set value (if there is one) on toggle
      this.args.issuerConfig.rollbackAttributes();
    }
    if (accessType === 'wif') {
      // reset all "account" attributes that are mutually exclusive with "wif" attributes
      // these attributes are different for each engine
      type === 'azure'
        ? (model.clientSecret = model.rootPasswordTtl = undefined)
        : type === 'aws'
        ? (model.accessKey = undefined)
        : null;
    }
  }

  @action
  onCancel() {
    this.resetErrors();
    this.args.model.unloadRecord();
    this.transition();
  }
}
