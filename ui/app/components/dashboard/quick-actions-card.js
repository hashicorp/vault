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

const getActionsByEngineType = (type) => {
  switch (type) {
    case 'kv':
      return [
        {
          actionTitle: 'Find KV secrets',
          actionType: 'find-kv',
        },
      ];
    case 'database':
      return [
        {
          actionTitle: 'Generate credentials for database',
          actionType: 'generate-credentials-db',
        },
      ];
    case 'pki':
      return [
        {
          actionTitle: 'Issue certificate',
          actionType: 'issue-certificate-pki',
        },
        {
          actionTitle: 'View certificate',
          actionType: 'view-certificate-pki',
        },
        {
          actionTitle: 'View issuer',
          actionType: 'view-issuer-pki',
        },
      ];
  }
};

const getActionParamByAction = (type) => {
  switch (type) {
    case 'find-kv':
      return {
        title: 'Secret Path',
        subText: 'Path of the secret you want to read, including the mount. E.g., secret/data/foo.',
        // elementType: 'input',
        buttonText: 'Read secrets',
        model: 'secret-v2',
        path: 'vault.cluster.secrets.backend.list',
        // ends in forward --> 'vault.cluster.secrets.backends.list', id
        // else --> 'vault.cluster.secrets.backends.show', id
      };
    case 'generate-credentials-db':
      return {
        title: 'Role to use',
        elementType: 'search-select',
        buttonText: 'Generate credentials',
        model: 'database/role',
        path: 'vault.cluster.secrets.backend.credentials',
      };
    case 'issue-certificate-pki':
      return {
        title: 'Role to use',
        elementType: 'search-select',
        placeholder: 'Type to find a role...',
        buttonText: 'Issue leaf certificate',
        model: 'pki/role',
        path: 'vault.cluster.secrets.backend.pki.roles.role.generate',
      };
    case 'view-certificate-pki':
      return {
        title: 'Certificate serial number',
        placeholder: '33:a3:...',
        elementType: 'search-select',
        buttonText: 'View certificate',
        model: 'pki/certificate/base',
        path: 'vault.cluster.secrets.backend.pki.certificates.certificate.details',
      };
    case 'view-issuer-pki':
      return {
        title: 'Issuer',
        placeholder: 'Type issuer name or ID',
        elementType: 'search-select',
        buttonText: 'View issuer',
        model: 'pki/issuer',
        nameKey: 'issuerName',
        path: 'vault.cluster.secrets.backend.pki.issuers.issuer.details',
      };
  }
};

const PKI_ENGINE_TYPE = 'pki';
const KV_ENGINE_TYPE = 'kv';
const DB_ENGINE_TYPE = 'database';
export default class DashboardQuickActionsCard extends Component {
  @service router;

  @tracked selectedEngine;
  @tracked selectedActions = [];
  @tracked selectedAction;
  @tracked actionParamField;
  @tracked selectedEngineName;
  @tracked value;

  constructor() {
    super(...arguments);

    if (!this.selectedEngine) {
      this.setSelectedAction();
    }
  }

  get filteredSecretEngines() {
    return this.args.secretsEngines.filter(
      (secretEngine) =>
        secretEngine.type === PKI_ENGINE_TYPE ||
        secretEngine.type === KV_ENGINE_TYPE ||
        secretEngine.type === DB_ENGINE_TYPE
    );
  }

  get secretsEnginesOptions() {
    return this.filteredSecretEngines.map((engine) => {
      const { id, type } = engine;
      return { name: id, type, id };
    });
  }

  @action
  handleSearchEngineSelect([selection]) {
    if (selection) {
      this.selectedEngine = selection.type;
      this.selectedEngineName = selection.id;
      this.selectedActions = getActionsByEngineType(this.selectedEngine);
      this.setSelectedAction(this.selectedActions?.firstObject.actionType);
    } else {
      // no selection, clear tracked properties
      this.selectedEngine = '';
      this.selectedEngineName = '';
      this.setSelectedAction();
    }
  }

  @action
  setSelectedAction(selectedAction) {
    if (selectedAction) {
      this.selectedAction = selectedAction;
      this.actionParamField = getActionParamByAction(selectedAction);
    } else {
      this.selectedActions = getActionsByEngineType('kv');
      this.actionParamField = getActionParamByAction('find-kv');
    }
  }

  @action
  handleSearch(val) {
    if (Array.isArray(val)) {
      this.value = val[0];
    } else {
      this.value = val;
    }
  }

  @action
  navigateToPage() {
    this.router.transitionTo(this.actionParamField.path, this.selectedEngineName, this.value);
  }
}
