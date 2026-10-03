/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';

import type ApiService from 'vault/services/api';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type FlashMessageService from 'vault/services/flash-messages';
import type RouterService from '@ember/routing/router-service';
import type { HTMLElementEvent } from 'vault/forms';

export default class KmipScopesCreatePageComponent extends Component {
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPath;
  @service declare readonly flashMessages: FlashMessageService;
  @service('app-router') declare readonly router: RouterService;

  @tracked declare name: string;
  @tracked validationError = false;
  @tracked invalidFormAlert: string | null = null;
  @tracked errorMessage: string | null = null;

  save = task(
    waitFor(async (event: HTMLElementEvent<HTMLFormElement>) => {
      event.preventDefault();

      if (!this.name) {
        this.validationError = true;
        this.invalidFormAlert = 'There is an error with this form.';
      } else {
        this.validationError = false;
        try {
          await this.api.secrets.kmipCreateScope(this.name, this.secretMountPath.currentPath, {});
          this.flashMessages.success(`Successfully created scope ${this.name}`);
          this.router.transitionTo('vault.cluster.secrets.backend.kmip.scopes.index');
        } catch (error) {
          const { message } = await this.api.parseError(error);
          this.errorMessage = message;
          this.invalidFormAlert = 'There was an error submitting this form.';
        }
      }
    })
  );
}
