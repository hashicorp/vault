/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */
import Component from '@glimmer/component';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';

import type SecretMountPath from 'vault/services/secret-mount-path';
import type ApiService from 'vault/services/api';
import type RouterService from '@ember/routing/router-service';
import type FlashMessageService from 'vault/services/flash-messages';
import type { ValidationMap } from 'vault/app-types';
import type { HTMLElementEvent } from 'vault/forms';
import type KmipConfigForm from 'vault/forms/secrets/kmip/config';

interface Args {
  form: KmipConfigForm;
}

export default class KmipConfigurePageComponent extends Component<Args> {
  @service declare readonly secretMountPath: SecretMountPath;
  @service declare readonly api: ApiService;
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly flashMessages: FlashMessageService;

  @tracked modelValidations: ValidationMap | null = null;
  @tracked invalidFormAlert: string | null = null;
  @tracked errorMessage: string | null = null;

  save = task(
    waitFor(async (event: HTMLElementEvent<HTMLFormElement>) => {
      event.preventDefault();
      this.errorMessage = null;

      try {
        const { isValid, state, invalidFormMessage, data } = this.args.form.toJSON();
        this.modelValidations = isValid ? null : state;
        this.invalidFormAlert = isValid ? '' : invalidFormMessage;

        if (isValid) {
          await this.api.secrets.kmipConfigure(this.secretMountPath.currentPath, data);
          this.flashMessages.success('Successfully configured KMIP engine');
          this.router.transitionTo('vault.cluster.secrets.backend.kmip.configuration');
        }
      } catch (error) {
        const { message } = await this.api.parseError(error);
        this.errorMessage = message;
        this.invalidFormAlert = 'There was an error submitting this form.';
      }
    })
  );
}
