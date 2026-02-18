/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import Component from '@glimmer/component';
import { task } from 'ember-concurrency';

import type { PATH_MAP } from 'vault/utils/constants/capabilities';
import type ApiService from 'vault/services/api';
import type CapabilitiesService from 'vault/services/capabilities';
import type FlashMessageService from 'vault/services/flash-messages';
import type RouterService from '@ember/routing/router-service';
import type SecretsEngineResource from 'vault/resources/secrets/engine';

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

const ROUTE_PATH_MAP = {
  'vault.cluster.secrets.backend.pki.certificates.index': ['pkiCertificates'],
  'vault.cluster.secrets.backend.pki.roles.index': ['pkiRoles'],
  'vault.cluster.secrets.backend.pki.tidy.index': ['pkiTidy', 'pkiTidyStatus', 'pkiConfigAutoTidy'],
} satisfies Record<string, readonly (keyof typeof PATH_MAP)[]>;

export default class PkiPageHeader extends Component<Args> {
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly api: ApiService;
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly capabilities: CapabilitiesService;

  @tracked engineToDisable = undefined;

  get breadcrumbs() {
    return [
      { label: 'Vault', route: 'vault', icon: 'vault', linkExternal: true },
      { label: 'Secrets engines', route: 'secrets', linkExternal: true },
      { label: this.args?.backend?.id },
    ];
  }

  // PKI does not make capability requests for these routes
  // so manually pass the relevant paths for each route.
  get policyPaths() {
    const backend = this.args?.backend?.id;
    const { currentRouteName } = this.router;
    const paths = ROUTE_PATH_MAP[currentRouteName as keyof typeof ROUTE_PATH_MAP];
    if (paths) {
      return this.capabilities.pathsForList(paths, { backend });
    }
    return null;
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
