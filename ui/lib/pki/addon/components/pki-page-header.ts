/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import Component from '@glimmer/component';
import { task } from 'ember-concurrency';

import type SecretsEngineResource from 'vault/resources/secrets/engine';
import type RouterService from '@ember/routing/router-service';
import type FlashMessageService from 'vault/services/flash-messages';
import type ApiService from 'vault/services/api';

/**
 * @module PkiPageHeader
 * The `PkiPageHeader` is used to display pki page headers.
 *
 * @example ```js
 * <PkiPageHeader @backend="exampleBackend" />
 * ```
 */

interface Args {
  backend: { id: string };
}

export default class PkiPageHeader extends Component<Args> {
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly api: ApiService;
  @service declare readonly flashMessages: FlashMessageService;

  @tracked engineToDisable = undefined;
  get breadcrumbs() {
    return [
      {
        label: 'Secrets',
        route: 'secrets',
        linkExternal: true,
      },
      {
        label: this.args?.backend?.id,
      },
    ];
  }

  @task
  *disableEngine(engine: SecretsEngineResource) {
    const { engineType, id, path } = engine;

    try {
      yield this.api.sys.mountsDisableSecretsEngine(id);
      this.flashMessages.success(`The ${engineType} Secrets Engine at ${path} has been disabled.`);
      this.router.transitionTo('vault.cluster.secrets.backends');
    } catch (err) {
      const { message } = yield this.api.parseError(err);
      this.flashMessages.danger(
        `There was an error disabling the ${engineType} Secrets Engine at ${path}: ${message}.`
      );
    } finally {
      this.engineToDisable = undefined;
    }
  }
}
