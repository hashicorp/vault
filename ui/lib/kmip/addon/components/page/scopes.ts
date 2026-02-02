/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { getOwner } from '@ember/owner';
import { task } from 'ember-concurrency';

import type RouterService from '@ember/routing/router-service';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type ApiService from 'vault/services/api';
import type { CapabilitiesMap, EngineOwner } from 'vault/app-types';
import type SecretsEngineResource from 'vault/resources/secrets/engine';
import FlashMessageService from 'vault/services/flash-messages';

interface Args {
  scopes: string[];
  capabilities: CapabilitiesMap;
}

export default class KmipScopesPageComponent extends Component<Args> {
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly secretMountPath: SecretMountPath;
  @service declare readonly api: ApiService;
  @service declare readonly flashMessages: FlashMessageService;
  @tracked engineToDisable: SecretsEngineResource | undefined = undefined;

  @tracked scopeToDelete: string | null = null;

  get mountPoint() {
    return (getOwner(this) as EngineOwner).mountPoint;
  }

  get paginationQueryParams() {
    return (page: number) => ({ page });
  }

  @action
  onFilterChange(pageFilter: string) {
    this.router.transitionTo({ queryParams: { pageFilter } });
  }

  @action
  async deleteScope() {
    try {
      await this.api.secrets.kmipDeleteScope(this.scopeToDelete as string, this.secretMountPath.currentPath);
      this.flashMessages.success(`Successfully deleted scope ${this.scopeToDelete}`);
      this.scopeToDelete = null;
      this.router.refresh();
    } catch (error) {
      const { message } = await this.api.parseError(error);
      this.flashMessages.danger(`Error deleting scope ${this.scopeToDelete}: ${message}`);
    }
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
