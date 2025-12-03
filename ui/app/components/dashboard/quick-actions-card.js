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

const QUICK_ACTION_ENGINES = ['pki', 'database'];

export default class DashboardQuickActionsCard extends Component {
  @service router;
  @service api;
  @service flashMessages;

  @tracked selectedEngine;
  @tracked selectedAction;
  @tracked paramValue;
  @tracked searchSelectOptions = [];

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
          buttonText: 'Generate credentials',
          route: 'vault.cluster.secrets.backend.credentials',
        };
      case 'Issue certificate':
        return {
          title: 'Role to use',
          placeholder: 'Type to find a role',
          buttonText: 'Issue leaf certificate',
          route: 'vault.cluster.secrets.backend.pki.roles.role.generate',
        };
      case 'View certificate':
        return {
          title: 'Certificate serial number',
          placeholder: '33:a3:...',
          buttonText: 'View certificate',
          route: 'vault.cluster.secrets.backend.pki.certificates.certificate.details',
        };
      case 'View issuer':
        return {
          title: 'Issuer',
          placeholder: 'Type issuer name or ID',
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

  get filteredSecretEngines() {
    return this.args.secretsEngines?.filter(
      (engine) => (engine.type === 'kv' && engine.version == 2) || QUICK_ACTION_ENGINES.includes(engine.type)
    );
  }

  get mountOptions() {
    return this.filteredSecretEngines?.map((engine) => {
      const { id, type } = engine;

      return { name: id, type, id };
    });
  }

  fetchOptions = task(
    waitFor(async () => {
      this.searchSelectOptions = [];
      const action = this.selectedAction;
      const api = this.api.secrets;
      const catchError = (e) => (e.response.status === 404 ? { keys: [] } : Promise.reject(e));

      try {
        // kv-suggestion-input fetches secrets internally -- this handles the remaining action types
        if (action && action !== 'Find KV secrets') {
          // fetch database roles, pki roles, pki certificates or pki issuers
          const methods = {
            'Generate credentials for database': ['databaseListStaticRoles', 'databaseListRoles'],
            'Issue certificate': ['pkiListRoles'],
            'View certificate': ['pkiListCerts'],
            'View issuer': ['pkiListIssuers'],
          }[this.selectedAction];
          const responses = await Promise.all(
            methods.map((method) => api[method](this.selectedEngine.id, true).catch(catchError))
          );
          responses.forEach((response) => {
            const options =
              action === 'View issuer'
                ? this.api.keyInfoToArray(response).map(({ id, issuer_name }) => ({
                    name: issuer_name || id,
                    id,
                  }))
                : response.keys.map((id) => ({ id }));
            this.searchSelectOptions.push(...options);
          });
        }
      } catch (e) {
        const { message } = await this.api.parseError(e);
        this.flashMessages.danger(`Error fetching options for selected action: ${message}`);
      }
    })
  );

  @action
  handleSearchEngineSelect([selection]) {
    this.selectedEngine = selection;
    // reset tracked properties
    this.selectedAction = null;
    this.paramValue = null;
  }

  @action
  setSelectedAction(selectedAction) {
    this.selectedAction = selectedAction;
    this.paramValue = null;
    this.fetchOptions.perform();
  }

  @action
  handleActionSelect(val) {
    if (Array.isArray(val)) {
      this.paramValue = val[0];
    } else {
      this.paramValue = val;
    }
  }

  @action
  navigateToPage() {
    let route = this.searchSelectParams.route;
    // If search-select falls back to stringInput, paramValue is a string not object
    let param = this.paramValue.id || this.paramValue;

    // kv has a special use case where if the paramValue ends in a '/' you should
    // link to different route
    if (this.selectedEngine.type === 'kv') {
      const path = this.paramValue.path || this.paramValue;
      route = pathIsDirectory(path)
        ? 'vault.cluster.secrets.backend.kv.list-directory'
        : 'vault.cluster.secrets.backend.kv.secret.index';
      param = path;
    }

    this.router.transitionTo(route, this.selectedEngine.id, param);
  }
}
