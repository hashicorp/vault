/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

/**
 * @module DashboardQuickActionsCard
 * DashboardQuickActionsCard component allows users to see a list of secrets engines filtered by
 * kv, pki and database and perform certain actions based on the type of secret engine selected
 *
 * @example
 * ```js
 * <DashboardQuickActionsCard />
 * ```
 * @param {array}
 */

const getQuickActions = (type) => {
  switch (type) {
    case 'kv':
      return [
        {
          actionTitle: 'Find KV secrets',
          actionType: 'find-kv',
          path: 'vault.cluster.secrets',
        },
      ];
    case 'database':
      return [
        {
          actionTitle: 'Generate credentials for database',
          actionType: 'generate-credentials-db',
          path: 'vault.cluster.database',
        },
      ];
    case 'pki':
      return [
        {
          actionTitle: 'Generate certificate',
          actionType: 'generate-certificate-pki',
          path: 'vault.cluster.pki',
        },
        {
          actionTitle: 'View certificate',
          actionType: 'view-certificate-pki',
          path: 'vault.cluster.pki',
        },
        {
          actionTitle: 'View issuer',
          actionType: 'view-issuer-pki',
          path: 'vault.cluster.pki',
        },
      ];
  }
};

const getActionMethod = (type) => {
  switch (type) {
    case 'find-kv':
      return {
        title: 'Secret Path',
        subText: 'Path of the secret you want to read, including the mount. E.g., secret/data/foo.',
        elementType: 'input',
      };
    case 'generate-credentials-db':
      return { title: 'Role to use', elementType: 'select' };
    case 'generate-certificate-pki':
      return { title: 'Role to use', elementType: 'select' };
    case 'view-certificate-pki':
      return { title: 'Certificate serial number', placeholder: '33:a3:...', elementType: 'search-select' };
    case 'view-issuer-pki':
      return { title: 'Issuer', placeholder: 'Type issuer name or ID', elementType: 'search-select' };
  }
};

export default class DashboardQuickActionsCard extends Component {
  @tracked selectedEngine;
  @tracked selectedQuickActions = [];
  @tracked selectedAction = null;
  @tracked actionType = '';

  get filteredSecretEngines() {
    return this.args.secretsEngines.filter(
      (secretEngine) =>
        (secretEngine.shouldIncludeInList && secretEngine.type === 'pki') ||
        secretEngine.type === 'kv' ||
        secretEngine.type === 'database'
    );
  }

  get secretsEnginesOptions() {
    return this.filteredSecretEngines.map((filteredSecretEngine) => ({
      name: filteredSecretEngine.path,
      id: filteredSecretEngine.type,
    }));
  }

  @action
  onSearchEngineSelect([selectedSearchEngines]) {
    this.selectedEngine = selectedSearchEngines;
    this.selectedQuickActions = getQuickActions(selectedSearchEngines);
    if (!this.selectedAction) {
      this.selectedAction = this.setSelectedAction(getQuickActions(selectedSearchEngines)?.[0]?.actionType);
    }
  }

  @action
  setSelectedAction(selectedAction) {
    this.selectedAction = selectedAction;
    this.actionType = getActionMethod(selectedAction);
  }
}
