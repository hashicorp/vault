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
import { next } from '@ember/runloop';

import type Router from '@ember/routing/router';
import type ApiService from 'vault/services/api';
import type VersionService from 'vault/services/version';
import type FlashMessageService from 'vault/services/flash-messages';
import type Owner from '@ember/owner';
import type AwsConfigForm from 'vault/forms/secrets/aws-config';
import type AzureConfigForm from 'vault/forms/secrets/azure-config';
import type GcpConfigForm from 'vault/forms/secrets/gcp-config';

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
    @configForm={{this.model.configForm}}
    />
 *
 * @param {string} backendPath - name of the secret engine, ex: 'azure-123'.
 * @param {string} displayName - used for flash messages, subText and labels. ex: 'Azure'.
 * @param {string} type - the type of the engine, ex: 'azure'.
 * @param {object} configForm - the config form for the engine. The field `isWifPluginConfigured` must be added to the form, otherwise this component will assert an error. `isWifPluginConfigured` should return true if any required wif fields have been set.
 */

type ConfigForm = AwsConfigForm | AzureConfigForm | GcpConfigForm;
interface Args {
  backendPath: string;
  displayName: string;
  type: string;
  configForm: ConfigForm;
}

export default class ConfigureWif extends Component<Args> {
  @service declare readonly router: Router;
  @service declare readonly api: ApiService;
  @service declare readonly version: VersionService;
  @service declare readonly flashMessages: FlashMessageService;

  @tracked errorMessage = '';
  @tracked invalidFormAlert = '';
  @tracked saveIssuerWarning = '';
  @tracked modelValidations: ValidationMap | null = null;

  disableAccessType = false;
  originalIssuer: string | undefined;

  constructor(owner: Owner, args: Args) {
    super(owner, args);
    // the following checks are only relevant to existing enterprise configurations
    const { isNew, data, isWifPluginConfigured, isAccountPluginConfigured } = this.args.configForm;

    if (this.version.isEnterprise && !isNew) {
      assert(
        `'isWifPluginConfigured' is required to be defined on the config model. Must return a boolean.`,
        isWifPluginConfigured !== undefined
      );
      next(() => {
        this.args.configForm.accessType = isWifPluginConfigured ? 'wif' : 'account';
      });
      // if wif or account only attributes are defined, disable the user's ability to change the access type
      this.disableAccessType = isWifPluginConfigured || isAccountPluginConfigured;
    }

    // cache the issuer to check if it has been changed later
    this.originalIssuer = data.issuer;
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
      const { isValid, state, invalidFormMessage, data } = this.args.configForm.toJSON();

      if (!isValid) {
        this.modelValidations = isValid ? null : state;
        this.invalidFormAlert = isValid ? '' : invalidFormMessage;
        return;
      }

      if (this.originalIssuer !== data.issuer) {
        // if the issuer has changed show modal with warning that the config will change
        // if the modal is shown, the user has to click confirm to continue saving
        this.saveIssuerWarning = `You are updating the global issuer config. This will overwrite Vault's current issuer ${
          !this.originalIssuer ? 'if it exists ' : ''
        }and may affect other configurations using this value. Continue?`;
        // exit task until user confirms
        return;
      }

      await this.save.perform();
    })
  );

  save = task(
    waitFor(async () => {
      try {
        const { data } = this.args.configForm.toJSON();
        const { issuer } = data;
        await this.saveConfig(data);
        if (this.originalIssuer !== issuer) {
          await this.updateIssuer(issuer as string);
        }
        this.flashMessages.success(`Successfully saved ${this.args.backendPath}'s configuration.`);
        this.transition();
      } catch (e) {
        const { message } = await this.api.parseError(e);
        this.errorMessage = message;
        this.invalidFormAlert = 'There was an error submitting this form.';
      }
    })
  );

  async updateIssuer(issuer: string) {
    try {
      await this.api.identity.oidcConfigure({ issuer });
    } catch (e) {
      const { message } = await this.api.parseError(e, 'Check Vault logs for details.');
      this.flashMessages.danger(`Issuer was not saved: ${message}`);
    }
  }

  async saveConfig(data: ConfigForm['data']) {
    const { backendPath, type } = this.args;
    if (type === 'aws') {
      await this.api.secrets.awsConfigureRootIamCredentials(backendPath, data);
      try {
        const { lease, leaseMax } = data as { lease: string; leaseMax: string };
        await this.api.secrets.awsConfigureLease(backendPath, { lease, leaseMax });
      } catch (e) {
        const { message } = await this.api.parseError(e);
        this.flashMessages.danger(`Error saving lease configuration: ${message}`);
      }
    } else if (type === 'azure') {
      await this.api.secrets.azureConfigure(backendPath, data);
    } else if (type === 'gcp') {
      await this.api.secrets.googleCloudConfigure(backendPath, data);
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

  @action
  onChangeAccessType(accessType: 'account' | 'wif') {
    const { configForm, type } = this.args;
    configForm.accessType = accessType;

    if (accessType === 'account') {
      // reset all "wif" attributes that are mutually exclusive with "account" attributes
      // these attributes are the same for each engine
      configForm.data.identityTokenAudience = configForm.data.identityTokenTtl = undefined;
    } else if (accessType === 'wif') {
      // reset all "account" attributes that are mutually exclusive with "wif" attributes
      // these attributes are different for each engine
      if (type === 'azure') {
        (configForm as AzureConfigForm).data.clientSecret = undefined;
      } else if (type === 'aws') {
        (configForm as AwsConfigForm).data.accessKey = undefined;
      }
    }
  }

  @action
  onCancel() {
    this.resetErrors();
    this.transition();
  }
}
