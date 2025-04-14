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

import type MountConfigModel from 'vault/vault/models/secret-engine/mount-config';
import type Router from '@ember/routing/router';
import type Store from '@ember-data/store';
import type VersionService from 'vault/services/version';
import type FlashMessageService from 'vault/services/flash-messages';

/**
 * @module TODO
 * @example
 * <SecretEngine::TuneMount
    todo
    />
 *
 */

interface Args {
  backendPath: string;
  displayName: string;
  type: string;
  mountConfigModel: MountConfigModel;
  tuneMountModel: TuneMountModel;
}

export default class TuneMount extends Component<Args> {
  @service declare readonly router: Router;
  @service declare readonly store: Store;
  @service declare readonly version: VersionService;
  @service declare readonly flashMessages: FlashMessageService;

  @tracked errorMessage = '';
  @tracked invalidFormAlert = '';
  @tracked modelValidations: ValidationMap | null = null;

  save = task(
    waitFor(async () => {
      // const mountConfigModelChanged = this.mountConfigModelAttrChanged;
      // const additionalModelAttrChanged = this.additionalConfigModelAttrChanged;
      // const issuerAttrChanged = this.issuerAttrChanged;
      // // check if any of the model(s) or issuer attributes have changed
      // // if no changes, transition and notify user
      // if (!mountConfigModelChanged && !additionalModelAttrChanged && !issuerAttrChanged) {
      //   this.flashMessages.info('No changes detected.');
      //   this.transition();
      //   return;
      // }
      // const mountConfigModelSaved = mountConfigModelChanged ? await this.saveMountConfigModel() : false;
      // const issuerSaved = issuerAttrChanged ? await this.updateIssuer() : false;
      // if (
      //   mountConfigModelSaved ||
      //   (!mountConfigModelChanged && issuerSaved) ||
      //   (!mountConfigModelChanged && additionalModelAttrChanged)
      // ) {
      //   // if there are changes made to the an additional model, attempt to save it. if saving fails, we transition and the failure will surface as a sticky flash message on the configuration details page.
      //   if (additionalModelAttrChanged) {
      //     await this.saveAdditionalConfigModel();
      //   }
      //   // we only prevent a transition if the mount config model or issuer fail when saving
      //   this.transition();
      // } else {
      //   return;
      // }
    })
  );

  async saveTuneMount(): Promise<boolean> {
    const { backendPath, mountConfigModel, tuneMountModel } = this.args;
    try {
      await tuneMountModel.save();
      this.flashMessages.success(`Successfully tuned ${backendPath}'s secret engine mount.`);
      return true;
    } catch (error) {
      this.errorMessage = errorMessage(error);
      this.invalidFormAlert = 'There was an error submitting this form.';
      return false;
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

  isValid(model: tuneMountModel) {
    const { isValid, state, invalidFormMessage } = model.validate();
    this.modelValidations = isValid ? null : state;
    this.invalidFormAlert = isValid ? '' : invalidFormMessage;
    return isValid;
  }

  @action
  onCancel() {
    this.resetErrors();
    this.transition();
  }
}
