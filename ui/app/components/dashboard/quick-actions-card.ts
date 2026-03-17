/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import { pathIsDirectory } from 'kv/utils/kv-breadcrumbs';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import {
  SecretsApiDatabaseListStaticRolesListEnum,
  SecretsApiDatabaseListRolesListEnum,
  SecretsApiPkiListRolesListEnum,
  SecretsApiPkiListCertsListEnum,
  SecretsApiPkiListIssuersListEnum,
} from '@hashicorp/vault-client-typescript';

import type RouterService from '@ember/routing/router-service';
import type ApiService from 'vault/services/api';
import type FlashMessageService from 'vault/services/flash-messages';
import type SecretsEngineResource from 'vault/resources/secrets/engine';
import type { StandardListResponse, PkiListIssuersResponse } from '@hashicorp/vault-client-typescript';

/**
 * @module DashboardQuickActionsCard
 * DashboardQuickActionsCard component allows users to see a list of secrets engines filtered by
 * kv, pki and database and perform certain actions based on the type of secret engine selected
 *
 * @example
 * ```js
 *   <Dashboard::QuickActionsCard @secretsEngines={{@model.secretsEngines}} />
 * ```
 */

interface Args {
  secretsEngines: SecretsEngineResource[];
}

export default class DashboardQuickActionsCard extends Component<Args> {
  @service declare readonly router: RouterService;
  @service declare readonly api: ApiService;
  @service declare readonly flashMessages: FlashMessageService;

  @tracked declare selectedEngine: SecretsEngineResource;
  @tracked selectedAction: string | null = null;
  @tracked paramValue: string | null = null;
  @tracked searchSelectOptions: string[] = [];

  QUICK_ACTION_ENGINES = ['pki', 'database'];

  get actionOptions() {
    switch (this.selectedEngine.type) {
      case 'kv':
        return ['Find KV secrets'];
      case 'database':
        return ['Generate credentials for database'];
      case 'pki':
        return ['Issue certificate', 'View certificate', 'View issuer'];
      default:
        return [];
    }
  }

  get searchSelectParams() {
    switch (this.selectedAction) {
      case 'Find KV secrets':
        return {
          isKV: true,
          buttonText: 'Read secrets',
          route: 'vault.cluster.secrets.backend.kv.secret.index',
        };
      case 'Generate credentials for database':
        return {
          title: 'Role to use',
          placeholder: 'Select a role',
          searchPlaceholder: 'Type to find a role',
          inputPlaceholder: 'Enter role name',
          buttonText: 'Generate credentials',
          route: 'vault.cluster.secrets.backend.credentials',
        };
      case 'Issue certificate':
        return {
          title: 'Role to use',
          placeholder: 'Select a role',
          searchPlaceholder: 'Type to find a role',
          inputPlaceholder: 'Enter role name',
          buttonText: 'Issue leaf certificate',
          route: 'vault.cluster.secrets.backend.pki.roles.role.generate',
        };
      case 'View certificate':
        return {
          title: 'Certificate serial number',
          placeholder: 'Select certificate',
          searchPlaceholder: '33:a3:...',
          inputPlaceholder: 'Enter certificate serial number',
          buttonText: 'View certificate',
          route: 'vault.cluster.secrets.backend.pki.certificates.certificate.details',
        };
      case 'View issuer':
        return {
          title: 'Issuer',
          placeholder: 'Select issuer',
          searchPlaceholder: 'Type issuer name or ID',
          inputPlaceholder: 'Enter issuer name or ID',
          buttonText: 'View issuer',
          route: 'vault.cluster.secrets.backend.pki.issuers.issuer.details',
          shouldRenderName: true,
        };
      default:
        return {
          placeholder: 'Please select an action above',
          buttonText: 'Select an action',
        };
    }
  }

  get filteredSecretsEngines() {
    return this.args.secretsEngines?.filter(
      (engine) =>
        (engine.type === 'kv' && engine.version == 2) || this.QUICK_ACTION_ENGINES.includes(engine.type)
    );
  }

  getListRequests() {
    const { secrets } = this.api;
    const { id } = this.selectedEngine;
    const action = this.selectedAction;

    if (action === 'Generate credentials for database') {
      return [
        secrets.databaseListStaticRoles(id, SecretsApiDatabaseListStaticRolesListEnum.TRUE),
        secrets.databaseListRoles(id, SecretsApiDatabaseListRolesListEnum.TRUE),
      ];
    } else if (action === 'Issue certificate') {
      return [secrets.pkiListRoles(id, SecretsApiPkiListRolesListEnum.TRUE)];
    } else if (action === 'View certificate') {
      return [secrets.pkiListCerts(id, SecretsApiPkiListCertsListEnum.TRUE)];
    } else if (action === 'View issuer') {
      return [secrets.pkiListIssuers(id, SecretsApiPkiListIssuersListEnum.TRUE)];
    }
    return [];
  }

  fetchOptions = task(
    waitFor(async () => {
      this.searchSelectOptions = [];
      const action = this.selectedAction;
      // kv-suggestion-input fetches secrets internally -- this handles the remaining action types
      if (action && action !== 'Find KV secrets') {
        // fetch database roles, pki roles, pki certificates or pki issuers
        const requests = this.getListRequests();
        const results = await Promise.allSettled(requests);

        for (const result of results) {
          // ignore failures and only extract data from successful requests
          // if there are no options the user will have the opportunity to enter a value via text input
          if (result.status === 'fulfilled') {
            if (action === 'View issuer') {
              const response = result.value as PkiListIssuersResponse;
              const options = this.api
                .keyInfoToArray(response)
                .map(({ id, issuer_name }) => issuer_name || id) as string[];
              this.searchSelectOptions.push(...options);
            } else {
              const { keys = [] } = result.value as StandardListResponse;
              this.searchSelectOptions.push(...keys);
            }
          }
        }
      }
    })
  );

  @action
  handleSearchEngineSelect(secretsEngine: SecretsEngineResource) {
    this.selectedEngine = secretsEngine;
    // reset tracked properties
    this.selectedAction = null;
    this.paramValue = null;
  }

  @action
  setSelectedAction(selectedAction: string) {
    this.selectedAction = selectedAction;
    this.paramValue = null;
    this.fetchOptions.perform();
  }

  @action
  navigateToPage() {
    let route = this.searchSelectParams.route;
    // kv has a special use case where if the paramValue ends in a '/' you should
    // link to different route
    if (this.selectedEngine.type === 'kv') {
      route = pathIsDirectory(this.paramValue)
        ? 'vault.cluster.secrets.backend.kv.list-directory'
        : 'vault.cluster.secrets.backend.kv.secret.index';
    }

    this.router.transitionTo(route, this.selectedEngine.id, this.paramValue);
  }
}
