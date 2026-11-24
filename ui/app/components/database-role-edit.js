/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */
import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { waitFor } from '@ember/test-waiters';
import { task } from 'ember-concurrency';
import DatabaseConnectionForm, { normalizeDatabaseConnection } from 'vault/forms/database/connection';
import { getRoleFieldAttrs, getRoleSettingAttrs, getRoleShowFields } from 'vault/forms/database/role';
import { addToArray } from 'vault/helpers/add-to-array';
import { removeFromArray } from 'vault/helpers/remove-from-array';

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
 * @param {object} model - The database role. A DatabaseRoleForm in create and edit mode, plain data in show mode.
 * @param {object} capabilities - Capabilities for the role: canDelete, canEditRole, canGenerateCredentials, canGetCredentials, canRotateRoleCredentials, canCreateDynamic, canCreateStatic.
 * @param {string} tab - The tab to render.
 * @param {string} mode - The mode to render. Either 'create' or 'edit'.
 * @param {string} [initialKey] - The initial key to set for the database role.
 */
const LIST_ROOT_ROUTE = 'vault.cluster.secrets.backend.list-root';
const SHOW_ROUTE = 'vault.cluster.secrets.backend.show';
export default class DatabaseRoleEdit extends Component {
  @service api;
  @service controlGroup;
  @service router;
  @service flashMessages;
  @service version;
  @tracked modelValidations;
  @tracked invalidFormAlert;
  @tracked errorMessage = '';
  @tracked saveIssuerWarning = '';

  constructor() {
    super(...arguments);
    if (this.args.initialKey) {
      this.args.model.database = [this.args.initialKey];
    }
  }

  breadcrumbs = [
    { label: 'Vault', text: 'Vault', icon: 'vault', path: 'vault.cluster.dashboard' },
    { text: 'Secrets engines', path: 'vault.cluster.secrets.backends' },
    this.args.root,
    { label: this.title, text: this.title },
  ];

  isValid() {
    const { isValid, state } = this.args.model.toJSON();
    this.modelValidations = isValid ? null : state;
    this.invalidFormAlert = isValid ? '' : 'There was an error submitting this form.';
    return isValid;
  }

  // Control group responses are also 403s, so they are checked first and start the approval flow.
  async errorMessageFor(e) {
    const { message, response } = await this.api.parseError(e);
    if (response?.isControlGroupError) {
      this.controlGroup.saveTokenFromError(response);
      return this.controlGroup.logFromError(response).content;
    }
    return message;
  }

  resetErrors() {
    this.flashMessages.clearMessages();
    this.errorMessage = this.invalidFormAlert = '';
    this.modelValidations = null;
  }

  get title() {
    if (this.args?.mode === 'create') {
      return 'Create role';
    } else if (this.args?.mode === 'edit') {
      return 'Edit role';
    } else {
      return this.args?.model?.name;
    }
  }

  // main fields on the create / edit form
  get fieldAttrs() {
    return getRoleFieldAttrs();
  }

  // role settings and statements, which of them render is decided by DatabaseRoleSettingForm
  get roleSettingAttrs() {
    return getRoleSettingAttrs(this.version.isEnterprise);
  }

  // fields displayed on the show view
  get showFields() {
    const { type, database } = this.args.model;
    return getRoleShowFields(type, database);
  }

  get warningMessages() {
    const warnings = {};
    const { canCreateDynamic, canCreateStatic } = this.args.capabilities || {};
    const { type } = this.args.model;
    if (
      (type === 'dynamic' && canCreateDynamic === false) ||
      (type === 'static' && canCreateStatic === false)
    ) {
      warnings.type = `You don't have permissions to create this type of role.`;
    }
    return warnings;
  }

  // Resolves to the selected connection's params. The settings only render once this resolves (see the template)
  // because the "Rotate immediately" toggle reads its initial state from the form when it is constructed.
  // If the connection cannot be read it resolves to an empty object so the settings still render.
  get databaseParams() {
    const backend = this.args.model?.backend;
    const dbs = this.args.model?.database || [];
    if (!backend || dbs.length === 0) {
      return null;
    }
    return this.api.secrets
      .databaseReadConnectionConfiguration(dbs[0], backend)
      .then(({ data }) => {
        const { plugin_name, skip_static_role_import_rotation } = normalizeDatabaseConnection(
          data,
          dbs[0],
          backend
        );
        this.applyConnectionDefault(skip_static_role_import_rotation);
        return { plugin_name, skip_static_role_import_rotation };
      })
      .catch(() => ({}));
  }

  // A new static role rotates immediately unless its connection says otherwise. The default is set on the form
  // before the settings render so the "Rotate immediately" toggle starts in the right state. It is replaced when
  // the connection changes, but never once the user has chosen a value themselves.
  appliedConnectionDefault;
  applyConnectionDefault(skipImport) {
    const { mode, model } = this.args;
    if (mode !== 'create') return;
    const current = model.skip_import_rotation;
    if (current === undefined || current === this.appliedConnectionDefault) {
      model.skip_import_rotation = skipImport;
      this.appliedConnectionDefault = skipImport;
    }
  }

  // Adds or removes this role from the allowed_roles of its database connection.
  // A 403 means the user can only manage the role and not the connection, so it is ignored and the role is still saved.
  // A control group 403 is rethrown so the role is not saved until the connection update is authorized.
  async updateAllowedRoles(type) {
    const { name, backend, database } = this.args.model;
    const connectionName = database[0];
    try {
      const { data } = await this.api.secrets.databaseReadConnectionConfiguration(connectionName, backend);
      const connection = normalizeDatabaseConnection(data, connectionName, backend);
      const roles = [...(connection.allowed_roles || [])];
      const allowed_roles = type === 'add' ? addToArray(roles, name) : removeFromArray(roles, name);
      const form = new DatabaseConnectionForm({ ...connection, allowed_roles }, { isNew: false });
      // Uses api.request instead of api.secrets.databaseConfigureConnection because the generated request
      // serializer (DatabaseConfigureConnectionRequestToJSON) only keeps the generic connection params and would
      // drop plugin-specific ones such as connection_url, username and password.
      await this.api.request.post(`/${backend}/config/${connectionName}`, form.toJSON().data);
    } catch (e) {
      const { message, status, response } = await this.api.parseError(e);
      if (response?.isControlGroupError) throw e;
      if (status === 403) return;
      throw new Error(`Could not update allowed roles for selected database: ${message}`);
    }
  }

  @action continueSubmitForm() {
    this.saveIssuerWarning = '';
    this.saveRole.perform();
  }

  @action
  generateCreds(roleId, roleType = '') {
    this.router.transitionTo('vault.cluster.secrets.backend.credentials', roleId, {
      queryParams: { roleType },
    });
  }

  @action
  async delete() {
    const { name, backend, type } = this.args.model;
    try {
      await this.updateAllowedRoles('remove');
      if (type === 'static') {
        await this.api.secrets.databaseDeleteStaticRole(name, backend);
      } else {
        await this.api.secrets.databaseDeleteRole(name, backend);
      }
    } catch (e) {
      this.flashMessages.danger(await this.errorMessageFor(e));
      return;
    }
    try {
      this.router.transitionTo(LIST_ROOT_ROUTE, backend, { queryParams: { tab: 'role' } });
    } catch (e) {
      console.debug(e); // eslint-disable-line
    }
  }

  @action
  async handleCreateEditRole(evt) {
    evt.preventDefault();
    this.resetErrors();
    const { mode, model } = this.args;
    if (!this.isValid()) return;

    // if we're creating and rotating a static role immediately, warn user
    if (!model.skip_import_rotation && model.type === 'static' && mode === 'create') {
      this.saveIssuerWarning =
        "You have enabled 'Rotate immediately' for this static role. Vault will update the password immediately after you save. NOTE: This will disrupt any active use of this role outside of Vault.";
      return;
    }
    await this.saveRole.perform();
  }

  saveRole = task(
    waitFor(async () => {
      const { mode, model } = this.args;
      const { name, backend, type } = model;
      const { data } = model.toJSON();

      try {
        if (mode === 'create') {
          await this.updateAllowedRoles('add');
        }
        if (type === 'static') {
          await this.api.secrets.databaseWriteStaticRole(name, backend, data);
        } else {
          await this.api.secrets.databaseWriteRole(name, backend, data);
        }
        this.router.transitionTo(SHOW_ROUTE, `role/${name}`);
      } catch (e) {
        this.errorMessage = await this.errorMessageFor(e);
        this.flashMessages.danger(
          this.errorMessage || 'Could not save the role. Please check Vault logs for more information.'
        );
      }
    })
  );

  @action
  async rotateRoleCred(id) {
    const backend = this.args.model?.backend;
    try {
      await this.api.secrets.databaseRotateStaticRoleCredentials(id, backend);
      this.flashMessages.success(`Success: Credentials for ${id} role were rotated`);
    } catch (e) {
      const { message } = await this.api.parseError(e);
      this.flashMessages.danger(message);
    }
  }
}
