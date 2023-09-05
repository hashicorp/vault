/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { inject as service } from '@ember/service';
import { pathIsDirectory } from 'kv/utils/kv-breadcrumbs';
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

  @tracked selectedEngine;
  @tracked selectedAction;
  @tracked paramValue;

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
          title: 'Secret path',
          subText: 'Path of the secret you want to read.',
          buttonText: 'Read secrets',
          model: 'kv/metadata',
          route: 'vault.cluster.secrets.backend.kv.secret.details',
          nameKey: 'path',
          queryObject: { pathToSecret: '', backend: this.selectedEngine.id },
          objectKeys: ['path', 'id'],
        };
      case 'Generate credentials for database':
        return {
          title: 'Role to use',
          buttonText: 'Generate credentials',
          model: 'database/role',
          route: 'vault.cluster.secrets.backend.credentials',
          queryObject: { backend: this.selectedEngine.id },
        };
      case 'Issue certificate':
        return {
          title: 'Role to use',
          placeholder: 'Type to find a role',
          buttonText: 'Issue leaf certificate',
          model: 'pki/role',
          route: 'vault.cluster.secrets.backend.pki.roles.role.generate',
          queryObject: { backend: this.selectedEngine.id },
        };
      case 'View certificate':
        return {
          title: 'Certificate serial number',
          placeholder: '33:a3:...',
          buttonText: 'View certificate',
          model: 'pki/certificate/base',
          route: 'vault.cluster.secrets.backend.pki.certificates.certificate.details',
          queryObject: { backend: this.selectedEngine.id },
        };
      case 'View issuer':
        return {
          title: 'Issuer',
          placeholder: 'Type issuer name or ID',
          buttonText: 'View issuer',
          model: 'pki/issuer',
          route: 'vault.cluster.secrets.backend.pki.issuers.issuer.details',
          nameKey: 'issuerName',
          queryObject: { backend: this.selectedEngine.id },
          objectKeys: ['id', 'issuerName'],
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
    let route = this.searchSelectParams.route;
    let param = this.paramValue.id;

    // kv has a special use case where if the paramValue ends in a '/' you should
    // link to different route
    if (this.selectedEngine.type === 'kv') {
      route = pathIsDirectory(this.paramValue?.path)
        ? 'vault.cluster.secrets.backend.kv.list-directory'
        : 'vault.cluster.secrets.backend.kv.secret.details';
      param = this.paramValue?.path;
    }

    this.router.transitionTo(route, this.selectedEngine.id, param);
  }
}
