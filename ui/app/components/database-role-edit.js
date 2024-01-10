/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

const LIST_ROOT_ROUTE = 'vault.cluster.secrets.backend.list-root';
const SHOW_ROUTE = 'vault.cluster.secrets.backend.show';

export default class DatabaseRoleEdit extends Component {
  @service router;
  @service flashMessages;
  @service store;

  constructor() {
    super(...arguments);
    if (this.args.initialKey) {
      this.args.model.database = [this.args.initialKey];
    }
  }

  @tracked loading = false;

  get warningMessages() {
    const warnings = {};
    if (
      (this.args.model.type === 'dynamic' && this.args.model.canCreateDynamic === false) ||
      (this.args.model.type === 'static' && this.args.model.canCreateStatic === false)
    ) {
      warnings.type = `You don't have permissions to create this type of role.`;
    }
    return warnings;
  }

  get databaseType() {
    const backend = this.args.model?.backend;
    const dbs = this.args.model?.database || [];
    if (!backend || dbs.length === 0) {
      return null;
    }
    return this.store
      .queryRecord('database/connection', { id: dbs[0], backend })
      .then((record) => record.plugin_name)
      .catch(() => null);
  }

  @action
  generateCreds(roleId, roleType = '') {
    this.router.transitionTo('vault.cluster.secrets.backend.credentials', roleId, {
      queryParams: { roleType },
    });
  }

  @action
  delete() {
    const secret = this.args.model;
    const backend = secret.backend;
    return secret
      .destroyRecord()
      .then(() => {
        try {
          this.router.transitionTo(LIST_ROOT_ROUTE, backend, { queryParams: { tab: 'role' } });
        } catch (e) {
          console.debug(e); // eslint-disable-line
        }
      })
      .catch((e) => {
        this.flashMessages.danger(e.errors?.join('. '));
      });
  }

  @action
  handleCreateEditRole(evt) {
    evt.preventDefault();
    this.loading = true;

    const mode = this.args.mode;
    const roleSecret = this.args.model;
    const secretId = roleSecret.name;
    if (mode === 'create') {
      roleSecret.set('id', secretId);
      const path = roleSecret.type === 'static' ? 'static-roles' : 'roles';
      roleSecret.set('path', path);
    }
    return roleSecret
      .save()
      .then(() => {
        try {
          this.router.transitionTo(SHOW_ROUTE, `role/${secretId}`);
        } catch (e) {
          console.debug(e); // eslint-disable-line
        }
      })
      .catch((e) => {
        const errorMessage = e.errors?.join('. ') || e.message;
        this.flashMessages.danger(
          errorMessage || 'Could not save the role. Please check Vault logs for more information.'
        );
        this.loading = false;
      });
  }
  @action
  rotateRoleCred(id) {
    const backend = this.args.model?.backend;
    const adapter = this.store.adapterFor('database/credential');
    return adapter
      .rotateRoleCredentials(backend, id)
      .then(() => {
        this.flashMessages.success(`Success: Credentials for ${id} role were rotated`);
      })
      .catch((e) => {
        this.flashMessages.danger(e.errors);
      });
  }
}
