/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import { tracked } from '@glimmer/tracking';

import type RouterService from '@ember/routing/router-service';
import type PkiTidyForm from 'vault/forms/secrets/pki/tidy';
import type { TtlEvent } from 'vault/app-types';
import type ApiService from 'vault/services/api';
import type SecretMountPathService from 'vault/services/secret-mount-path';

interface Args {
  form: PkiTidyForm;
  tidyType: string;
  onSave: CallableFunction;
  onCancel: CallableFunction;
}

export default class PkiTidyFormComponent extends Component<Args> {
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPathService;

  @tracked errorBanner = '';
  @tracked invalidFormAlert = '';

  save = task(
    waitFor(async (event: Event) => {
      event.preventDefault();
      try {
        const { currentPath } = this.secretMountPath;
        const { data } = this.args.form.toJSON();

        if (this.args.tidyType === 'auto') {
          await this.api.secrets.pkiConfigureAutoTidy(currentPath, data);
        } else {
          await this.api.secrets.pkiTidy(currentPath, data);
        }
        this.args.onSave();
      } catch (e) {
        const { message } = await this.api.parseError(e);
        this.errorBanner = message;
        this.invalidFormAlert = 'There was an error submitting this form.';
      }
    })
  );

  @action
  handleAcmeTtl(e: TtlEvent) {
    const { enabled, goSafeTimeString } = e;
    this.args.form.data.acme_account_safety_buffer = goSafeTimeString;
    this.args.form.data.tidy_acme = enabled;
  }
}
