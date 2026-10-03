/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';

import type SecretsEngineResource from 'vault/resources/secrets/engine';
import type RouterService from '@ember/routing/router-service';
import type FlashMessageService from 'vault/services/flash-messages';
import type ApiService from 'vault/services/api';

/**
 * @module LdapHeader handles the ldap page header.
 *
 * @example
 * <SecretEngine::LdapHeader
 *    @model={{this.model}}
 *    />
 *
 * @param {object} secretsEngine - A model contains a ldap secret engine resource.
 * @param {object} config - A model contains the configuration of the ldap secret engine.
 */

interface Args {
  secretsEngine: SecretsEngineResource;
  config: Record<string, unknown>;
}

export default class LdapHeader extends Component<Args> {
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly api: ApiService;
  @service declare readonly flashMessages: FlashMessageService;

  @tracked engineToDisable: SecretsEngineResource | undefined = undefined;

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
