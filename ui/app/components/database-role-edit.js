/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { waitFor } from '@ember/test-waiters';
import { task } from 'ember-concurrency';
import errorMessage from 'vault/utils/error-message';

/**
 * @module DatabaseRoleEdit component is used to configure a database role.
 * See secret-edit-layout which uses options for backend to determine which component to render.
 * @example
 * <DatabaseRoleEdit
    @model={{this.model.database.role}}
    @tab="edit"
    @model="edit"
    @initialKey=this.initialKey
    />
 *
 * @param {object} model - The database role model.
 * @param {string} tab - The tab to render.
 * @param {string} mode - The mode to render. Either 'create' or 'edit'.
 * @param {string} [initialKey] - The initial key to set for the database role.
 */
const LIST_ROOT_ROUTE = 'vault.cluster.secrets.backend.list-root';
const SHOW_ROUTE = 'vault.cluster.secrets.backend.show';
export default class DatabaseRoleEdit extends Component {
  @service router;
  @service flashMessages;
  @service store;
  @tracked modelValidations;
  @tracked invalidFormAlert;
  @tracked errorMessage = '';

  constructor() {
    super(...arguments);
    if (this.args.initialKey) {
      this.args.model.database = [this.args.initialKey];
    }
  }

  isValid() {
    const { isValid, state } = this.args.model.validate();
    this.modelValidations = isValid ? null : state;
    this.invalidFormAlert = 'There was an error submitting this form.';
    return isValid;
  }

  resetErrors() {
    this.flashMessages.clearMessages();
    this.errorMessage = this.invalidFormAlert = '';
    this.modelValidations = null;
  }

  get warningMessages() {
    const warnings = {};
    const { canCreateDynamic, canCreateStatic, type } = this.args.model;
    if (
      (type === 'dynamic' && canCreateDynamic === false) ||
      (type === 'static' && canCreateStatic === false)
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

  handleCreateEditRole = task(
    waitFor(async (evt) => {
      evt.preventDefault();
      this.resetErrors();
      const { mode, model } = this.args;
      if (!this.isValid()) return;
      if (mode === 'create') {
        model.id = model.name;
        const path = model.type === 'static' ? 'static-roles' : 'roles';
        model.path = path;
      }
      try {
        await model.save();
        this.router.transitionTo(SHOW_ROUTE, `role/${model.name}`);
      } catch (e) {
        this.errorMessage = errorMessage(e);
        this.flashMessages.danger(
          this.errorMessage || 'Could not save the role. Please check Vault logs for more information.'
        );
      }
    })
  );

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
