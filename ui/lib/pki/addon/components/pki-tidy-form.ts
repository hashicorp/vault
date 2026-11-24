/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import { tracked } from '@glimmer/tracking';
import { isDurationOverMax } from 'core/utils/duration-utils';

import type RouterService from '@ember/routing/router-service';
import type PkiTidyForm from 'vault/forms/secrets/pki/tidy';
import type { TtlEvent, ValidationMap } from 'vault/app-types';
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

  // Durations the form loaded with, used to replace a hidden duration that is too large.
  initialData: Record<string, unknown>;

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    this.initialData = { ...args.form.data };
  }

  @tracked errorBanner = '';
  @tracked invalidFormAlert = '';
  @tracked modelValidations: ValidationMap | null = null;

  // Vault rejects a duration it can't store even when it isn't in use, but a hidden or toggled-off picker
  // can't show an error. Replace those values so they neither block saving nor get sent.
  resetInactiveDurations() {
    const { form, tidyType } = this.args;
    const data = form.data as unknown as Record<string, unknown>;
    // turning automatic tidy off hides every other field
    const inactive = tidyType === 'auto' && !data['enabled'] ? Object.keys(data) : [];
    if (!data['tidy_acme']) inactive.push('acme_account_safety_buffer');
    for (const key of inactive) {
      if (key !== 'enabled' && isDurationOverMax(data[key])) {
        data[key] = this.initialData[key];
      }
    }
  }

  save = task(
    waitFor(async (event: Event) => {
      event.preventDefault();
      try {
        this.resetInactiveDurations();
        const { currentPath } = this.secretMountPath;
        const { isValid, state, invalidFormMessage, data } = this.args.form.toJSON();
        this.modelValidations = isValid ? null : state;
        this.invalidFormAlert = invalidFormMessage;
        if (!isValid) return;

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
