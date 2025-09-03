/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { task } from 'ember-concurrency';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';

import type Router from '@ember/routing/router';
import type FlashMessageService from 'vault/services/flash-messages';
import type ApiService from 'vault/services/api';
import type SecretsEngineResource from 'vault/resources/secrets/engine';

/**
 * @module GeneralSettingsComponent is used to configure the SSH secret engine.
 *
 * @example
 * ```js
 * <Secrets:Page:GeneralSettings
 *    @model={{this.model}}
 *  />
 * ```
 *
 * @param {string} secretsEngine - secrets engine resource
 */

interface Args {
  model: {
    secretsEngine: SecretsEngineResource;
  };
}

export default class GeneralSettingsComponent extends Component<Args> {
  @service declare readonly router: Router;
  @service declare readonly api: ApiService;
  @service declare readonly flashMessages: FlashMessageService;

  @tracked errorMessage: string | null = null;
  @tracked invalidFormAlert: string | null = null;
  @tracked showUnsavedChangesModal = false;

  saveGeneralSettings = task(async (event) => {
    event.preventDefault();

    try {
      const fd = new FormData(event.target as HTMLFormElement);
      await this.api.sys.mountsTuneConfigurationParameters(this.args.model.secretsEngine.id, {
        // TODO: add other params when other card components are made
        description: fd.get('description') as string,
      });
      this.flashMessages.success('Engine settings successfully updated.');
    } catch (e) {
      // handle error state
      const { message } = await this.api.parseError(e);
      this.flashMessages.danger(`Try again or check your network connection. ${message}`);
    }
  });
}
