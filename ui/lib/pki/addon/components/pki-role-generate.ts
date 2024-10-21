/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { task } from 'ember-concurrency';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import errorMessage from 'vault/utils/error-message';
import type RouterService from '@ember/routing/router';
import type Store from '@ember-data/store';
import type FlashMessageService from 'vault/services/flash-messages';
import type DownloadService from 'vault/services/download';
import type PkiCertificateGenerateModel from 'vault/models/pki/certificate/generate';
import type PkiCertificateSignModel from 'vault/models/pki/certificate/sign';

interface Args {
  onSuccess: CallableFunction;
  model: PkiCertificateGenerateModel | PkiCertificateSignModel;
  type: string;
}

export default class PkiRoleGenerate extends Component<Args> {
  @service declare readonly download: DownloadService;
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly store: Store;
  @service('app-router') declare readonly router: RouterService;

  @tracked errorBanner = '';
  @tracked invalidFormAlert = '';

  get verb() {
    return this.args.type === 'sign' ? 'sign' : 'generate';
  }

  @task
  *save(evt: Event) {
    evt.preventDefault();
    this.errorBanner = '';
    const { model, onSuccess } = this.args;
    try {
      yield model.save();
      onSuccess();
    } catch (err) {
      this.errorBanner = errorMessage(err, `Could not ${this.verb} certificate. See Vault logs for details.`);
      this.invalidFormAlert = 'There was an error submitting this form.';
    }
  }

  @action cancel() {
    this.args.model.unloadRecord();
    this.router.transitionTo('vault.cluster.secrets.backend.pki.roles.role.details');
  }
}
