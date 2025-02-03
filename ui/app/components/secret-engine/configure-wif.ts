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
import { capitalize } from '@ember/string';
import errorMessage from 'vault/utils/error-message';

import type MountConfigModel from 'vault/vault/models/secret-engine/mount-config';
import type AdditionalConfigModel from 'vault/vault/models/secret-engine/additional-config';
import type IdentityOidcConfigModel from 'vault/models/identity/oidc/config';
import type Router from '@ember/routing/router';
import type StoreService from 'vault/services/store';
import type VersionService from 'vault/services/version';
import type FlashMessageService from 'vault/services/flash-messages';

/**
 * @module SecretEngineConfigureWif component is used to configure secret engines that allow the WIF (Workload Identity Federation) configuration.
 * The ability to configure WIF fields is an enterprise only feature.
 * If the user is configuring WIF attributes they will also have the option to update the global issuer config, which is a separate endpoint named identity/oidc/config.
 * If a user is on CE, the account configuration fields will display with no ability to select or see wif fields.
 * 
 * @example
 * <SecretEngine::ConfigureWif
    @backendPath={{this.model.id}}
    @displayName="AWS"
    @type="aws"
    @mountConfigModel={{this.model.mount-config-model}}
    @additionalConfigModel={{this.model.additional-config-model}}
    @issuerConfig={{this.model.identity-oidc-config}}
    />
 *
 * @param {string} backendPath - name of the secret engine, ex: 'azure-123'.
 * @param {string} displayName - used for flash messages, subText and labels. ex: 'Azure'.
 * @param {string} type - the type of the engine, ex: 'azure'.
 * @param {object} mountConfigModel - the config model for the engine. The attr `isWifPluginConfigured` must be added to this config model otherwise this component will assert an error. `isWifPluginConfigured` should return true if any required wif fields have been set.
 * @param {object} [additionalConfigModel] - for engines with two config models. Currently, only used by aws
 * @param {object} [issuerConfig] - the identity/oidc/config model. Will be passed in if user has an enterprise license.
 */

interface Args {
  backendPath: string;
  displayName: string;
  type: string;
  mountConfigModel: MountConfigModel;
  additionalConfigModel: AdditionalConfigModel;
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
    if (this.version.isCommunity && this.args.mountConfigModel.isNew) return;
    const { isWifPluginConfigured, isAccountPluginConfigured } = this.args.mountConfigModel;
    assert(
      `'isWifPluginConfigured' is required to be defined on the config model. Must return a boolean.`,
      isWifPluginConfigured !== undefined
    );
    this.accessType = isWifPluginConfigured ? 'wif' : 'account';
    // if wif or account only attributes are defined, disable the user's ability to change the access type
    this.disableAccessType = isWifPluginConfigured || isAccountPluginConfigured;
  }

  get mountConfigModelAttrChanged() {
    // "backend" dirties model state so explicity ignore it here
    return Object.keys(this.args.mountConfigModel?.changedAttributes()).some((item) => item !== 'backend');
  }

  get issuerAttrChanged() {
    return this.args.issuerConfig?.hasDirtyAttributes;
  }

  get additionalConfigModelAttrChanged() {
    const { additionalConfigModel } = this.args;
    // required to check for additional model otherwise Object.keys will have nothing to iterate over and fails
    return additionalConfigModel
      ? Object.keys(additionalConfigModel.changedAttributes()).some((item) => item !== 'backend')
      : false;
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
      // currently we only check validations on the additional model
      if (this.args.additionalConfigModel && !this.isValid(this.args.additionalConfigModel)) {
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
      const mountConfigModelChanged = this.mountConfigModelAttrChanged;
      const additionalModelAttrChanged = this.additionalConfigModelAttrChanged;
      const issuerAttrChanged = this.issuerAttrChanged;
      // check if any of the model(s) or issuer attributes have changed
      // if no changes, transition and notify user
      if (!mountConfigModelChanged && !additionalModelAttrChanged && !issuerAttrChanged) {
        this.flashMessages.info('No changes detected.');
        this.transition();
        return;
      }

      const mountConfigModelSaved = mountConfigModelChanged ? await this.saveMountConfigModel() : false;
      const issuerSaved = issuerAttrChanged ? await this.updateIssuer() : false;

      if (
        mountConfigModelSaved ||
        (!mountConfigModelChanged && issuerSaved) ||
        (!mountConfigModelChanged && additionalModelAttrChanged)
      ) {
        // if there are changes made to the an additional model, attempt to save it. if saving fails, we transition and the failure will surface as a sticky flash message on the configuration details page.
        if (additionalModelAttrChanged) {
          await this.saveAdditionalConfigModel();
        }
        // we only prevent a transition if the mount config model or issuer fail when saving
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

  async saveMountConfigModel(): Promise<boolean> {
    const { backendPath, mountConfigModel } = this.args;
    try {
      await mountConfigModel.save();
      this.flashMessages.success(`Successfully saved ${backendPath}'s configuration.`);
      return true;
    } catch (error) {
      this.errorMessage = errorMessage(error);
      this.invalidFormAlert = 'There was an error submitting this form.';
      return false;
    }
  }

  async saveAdditionalConfigModel() {
    const { backendPath, additionalConfigModel, type } = this.args;
    const additionalConfigModelName = type === 'aws' ? 'lease configuration' : 'additional configuration';
    try {
      await additionalConfigModel.save();
      this.flashMessages.success(`Successfully saved ${backendPath}'s ${additionalConfigModelName}.`);
    } catch (error) {
      // the only error the user sees is a sticky flash message on the next view.
      this.flashMessages.danger(
        `${capitalize(additionalConfigModelName)} was not saved: ${errorMessage(error)}`,
        {
          sticky: true,
        }
      );
    }
  }

  resetErrors() {
    this.flashMessages.clearMessages();
    this.errorMessage = this.invalidFormAlert = '';
    this.modelValidations = null;
  }

  transition() {
    this.router.transitionTo('vault.cluster.secrets.backend.configuration', this.args.backendPath);
  }

  isValid(model: AdditionalConfigModel) {
    const { isValid, state, invalidFormMessage } = model.validate();
    this.modelValidations = isValid ? null : state;
    this.invalidFormAlert = isValid ? '' : invalidFormMessage;
    return isValid;
  }

  @action
  onChangeAccessType(accessType: string) {
    this.accessType = accessType;
    const { mountConfigModel, type } = this.args;
    if (accessType === 'account') {
      // reset all "wif" attributes that are mutually exclusive with "account" attributes
      // these attributes are the same for each engine
      mountConfigModel.identityTokenAudience = mountConfigModel.identityTokenTtl = undefined;
      // return the issuer to the globally set value (if there is one) on toggle
      this.args.issuerConfig.rollbackAttributes();
    }
    if (accessType === 'wif') {
      // reset all "account" attributes that are mutually exclusive with "wif" attributes
      // these attributes are different for each engine
      type === 'azure'
        ? (mountConfigModel.clientSecret = mountConfigModel.rootPasswordTtl = undefined)
        : type === 'aws'
        ? (mountConfigModel.accessKey = undefined)
        : null;
    }
  }

  @action
  onCancel() {
    this.resetErrors();
    this.args.mountConfigModel.unloadRecord();
    this.transition();
  }
}
