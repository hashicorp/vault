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
 * <DashboardQuickActionsCard @secretsEngines={{@model.secretsEngines}} />
 * ```
 * @param {array}
 */

const ENGINE_TYPE = ['pki', 'kv', 'database'];

export default class DashboardQuickActionsCard extends Component {
  @service router;

  @tracked selectedEngine;
  @tracked selectedAction;
  @tracked selectedEngineName;
  @tracked paramValue;

  // 2nd input field
  get actionOptions() {
    switch (this.selectedEngine) {
      case 'kv':
        return ['Find KV secrets'];
      case 'database':
        return ['Generate credentials for database'];
      case 'pki':
        return ['Issue certificate', 'View certificate', 'View issuer'];
      default:
        return ['Find KV secrets'];
    }
  }

  // setting the search select object
  get searchSelectParams() {
    switch (this.selectedAction) {
      case 'Find KV secrets':
        return {
          title: 'Secret Path',
          subText: 'Path of the secret you want to read, including the mount. E.g., secret/data/foo.',
          buttonText: 'Read secrets',
          model: 'secret-v2',
          path: 'vault.cluster.secrets.backend.list',
          // ends in forward --> 'vault.cluster.secrets.backends.list', id
          // else --> 'vault.cluster.secrets.backends.show', id
        };
      case 'Generate credentials for database':
        return {
          title: 'Role to use',
          buttonText: 'Generate credentials',
          model: 'database/role',
          path: 'vault.cluster.secrets.backend.credentials',
        };
      case 'Issue certificate':
        return {
          title: 'Role to use',
          placeholder: 'Type to find a role...',
          buttonText: 'Issue leaf certificate',
          model: 'pki/role',
          path: 'vault.cluster.secrets.backend.pki.roles.role.generate',
        };
      case 'View certificate':
        return {
          title: 'Certificate serial number',
          placeholder: '33:a3:...',
          buttonText: 'View certificate',
          model: 'pki/certificate/base',
          path: 'vault.cluster.secrets.backend.pki.certificates.certificate.details',
        };
      case 'View issuer':
        return {
          title: 'Issuer',
          placeholder: 'Type issuer name or ID',
          buttonText: 'View issuer',
          model: 'pki/issuer',
          nameKey: 'issuerName',
          path: 'vault.cluster.secrets.backend.pki.issuers.issuer.details',
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
    return this.args.secretsEngines.filter((engine) => ENGINE_TYPE.includes(engine.type));
  }

  get mountOptions() {
    return this.filteredSecretEngines.map((engine) => {
      const { id, type } = engine;
      return { name: id, type, id };
    });
  }

  @action
  handleSearchEngineSelect([selection]) {
    this.selectedEngine = selection?.type;
    this.selectedEngineName = selection?.id;
    // reset tracked properties
    this.selectedAction = null;
  }

  @action
  setSelectedAction(selectedAction) {
    this.selectedAction = selectedAction;
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
    this.router.transitionTo(this.searchSelectParams.path, this.selectedEngineName, this.paramValue);
  }
}
