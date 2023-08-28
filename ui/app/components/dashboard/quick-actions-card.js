/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { inject as service } from '@ember/service';

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

const QUICK_ACTION_ENGINES = ['pki', 'kv', 'database'];
const KV_VERSION_1 = 'kv version 1';
const KV_VERSION_2 = 'kv version 2';
export default class DashboardQuickActionsCard extends Component {
  @service router;

  @tracked selectedEngine;
  @tracked selectedAction;
  @tracked paramValue;

  get actionOptions() {
    switch (this.selectedEngine.type) {
      case `kv version ${this.selectedEngine?.version}`:
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
          title: 'Secret path',
          subText: 'Path of the secret you want to read, including the mount. E.g., secret/data/foo.',
          buttonText: 'Read secrets',
          // check kv version to figure out which model to use
          model: this.selectedEngine.version === 2 ? 'secret-v2' : 'secret',
          route:
            this.selectedEngine.version === 2
              ? 'vault.cluster.secrets.backend.kv.secret.details'
              : 'vault.cluster.secrets.backend.show',
        };
      case 'Generate credentials for database':
        return {
          title: 'Role to use',
          buttonText: 'Generate credentials',
          model: 'database/role',
          route: 'vault.cluster.secrets.backend.credentials',
        };
      case 'Issue certificate':
        return {
          title: 'Role to use',
          placeholder: 'Type to find a role',
          buttonText: 'Issue leaf certificate',
          model: 'pki/role',
          route: 'vault.cluster.secrets.backend.pki.roles.role.generate',
        };
      case 'View certificate':
        return {
          title: 'Certificate serial number',
          placeholder: '33:a3:...',
          buttonText: 'View certificate',
          model: 'pki/certificate/base',
          route: 'vault.cluster.secrets.backend.pki.certificates.certificate.details',
        };
      case 'View issuer':
        return {
          title: 'Issuer',
          placeholder: 'Type issuer name or ID',
          buttonText: 'View issuer',
          model: 'pki/issuer',
          nameKey: 'issuerName',
          route: 'vault.cluster.secrets.backend.pki.issuers.issuer.details',
        };
      default:
        return {
          placeholder: 'Please select an action above',
          buttonText: 'Select an action',
          model: '',
        };
    }
  }

  get filteredSecretEngines() {
    return this.args.secretsEngines.filter((engine) => QUICK_ACTION_ENGINES.includes(engine.type));
  }

  get mountOptions() {
    return this.filteredSecretEngines.map((engine) => {
      let { id, type, version } = engine;
      if (type === 'kv') type = `kv version ${version}`;

      return { name: id, type, id, version };
    });
  }

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
    let searchSelectParamRoute = this.searchSelectParams.route;

    // kv has a special use case where if the paramValue ends in a '/' you should
    // link to different route
    if (this.selectedEngine.type === KV_VERSION_1) {
      searchSelectParamRoute =
        this.paramValue && this.paramValue?.endsWith('/')
          ? 'vault.cluster.secrets.backend.list'
          : 'vault.cluster.secrets.backend.show';
    }

    if (this.selectedEngine.type === KV_VERSION_2) {
      searchSelectParamRoute =
        this.paramValue && this.paramValue?.endsWith('/')
          ? 'vault.cluster.secrets.backend.kv.list-directory'
          : 'vault.cluster.secrets.backend.kv.secret.details';
    }

    this.router.transitionTo(searchSelectParamRoute, this.selectedEngine.id, this.paramValue);
  }
}
