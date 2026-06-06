/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { SecretsApiTransformListTransformationsListEnum } from '@hashicorp/vault-client-typescript';

/**
 * @module TransformRoleEdit
 * `TransformRoleEdit` is a component that allows you to create/edit or view a transform role.
 *
 * @example
 * ```js
 *   <TransformRoleEdit @form={{this.form}} @capabilities={{this.capabilities}} @mode={{this.mode}} />
 * ```
 * @param {object} form - RoleForm instance with data and formFields.
 * @param {object} capabilities - Object with canDelete, canUpdate, canRead capabilities.
 * @param {string} mode - Is either show, create or edit.
 */

export default class TransformRoleEditComponent extends Component {
  @service flashMessages;
  @service router;
  @service api;

  @tracked errorMessage = '';
  @tracked modelValidations;
  @tracked transformations = [];

  // Non-tracked: used only to diff added/removed transformations after save
  initialTransformations = [];

  constructor() {
    super(...arguments);
    this.initialTransformations = [...(this.args.form.data.transformations ?? [])];
    this.fetchTransformations();
  }

  async fetchTransformations() {
    try {
      const resp = await this.api.secrets.transformListTransformations(
        this.args.form.data.backend,
        SecretsApiTransformListTransformationsListEnum.TRUE
      );
      this.transformations = (resp.keys ?? []).map((key) => ({ id: key }));
    } catch {
      // swallow errors, SearchSelect will fall back to string-list
    }
  }

  get breadcrumbs() {
    const backend = this.args.form?.data?.backend;
    const name = this.args.form?.data?.name;
    return [
      { label: 'Vault', route: 'vault.cluster.dashboard', icon: 'vault' },
      { label: 'Secrets engines', route: 'vault.cluster.secrets.backends' },
      {
        label: backend,
        route: 'vault.cluster.secrets.backend.list-root',
        model: backend,
        query: { tab: 'role' },
      },
      { label: this.title },
      { label: this.args.mode === 'create' ? 'role' : name },
    ];
  }

  get title() {
    if (this.args.mode === 'create') {
      return 'Create role';
    } else if (this.args.mode === 'edit') {
      return 'Edit role';
    } else {
      return 'Role';
    }
  }

  get subtitle() {
    if (this.args.mode === 'show') {
      return this.args.form?.data?.name;
    }
    return '';
  }

  // Reads a transformation, updates its allowed_roles (add/remove this role), then saves.
  async syncTransformationForRole(transformationName, roleName, backend, syncAction) {
    let currentAllowedRoles;
    try {
      const resp = await this.api.secrets.transformReadTransformation(transformationName, backend);
      const data = resp?.data || resp || {};
      currentAllowedRoles = data.allowed_roles || [];
    } catch {
      // If the transformation can't be read, skip it
      return { transformationName, syncAction, errorStatus: null, skipped: true };
    }

    let updatedRoles;
    if (syncAction === 'ADD') {
      updatedRoles = currentAllowedRoles.includes(roleName)
        ? currentAllowedRoles
        : [...currentAllowedRoles, roleName];
    } else {
      updatedRoles = currentAllowedRoles.filter((r) => r !== roleName);
    }

    try {
      await this.api.secrets.transformWriteTransformation(transformationName, backend, {
        allowed_roles: updatedRoles,
      });
      return { transformationName, syncAction, errorStatus: null };
    } catch (writeErr) {
      const { status } = await this.api.parseError(writeErr);
      return { transformationName, syncAction, errorStatus: status };
    }
  }

  // Diffs current vs initial transformations, syncs allowed_roles on each
  // affected transformation, then shows a single contextual flash if any failed.
  async handleTransformationSync(roleName, backend, type = 'update') {
    const currentTransformations = this.args.form.data.transformations ?? [];
    const initialTransformations = this.initialTransformations;

    let syncOps;
    if (type === 'create') {
      syncOps = currentTransformations.map((t) => ({ id: t, syncAction: 'ADD' }));
    } else {
      const added = currentTransformations.filter((t) => !initialTransformations.includes(t));
      const removed = initialTransformations.filter((t) => !currentTransformations.includes(t));
      syncOps = [
        ...added.map((t) => ({ id: t, syncAction: 'ADD' })),
        ...removed.map((t) => ({ id: t, syncAction: 'REMOVE' })),
      ];
    }

    if (syncOps.length === 0) return;

    const results = await Promise.all(
      syncOps.map(({ id, syncAction }) => this.syncTransformationForRole(id, roleName, backend, syncAction))
    );

    const errors = results.filter((r) => r.errorStatus === 403);
    if (errors.length === 0) return;

    const errorAdding = errors.some((r) => r.syncAction === 'ADD');
    const errorRemoving = errors.some((r) => r.syncAction === 'REMOVE');

    let message;
    if (type === 'create') {
      message =
        'Transformations have been attached to this role, but the role was not added to those transformations\u2019 allowed_roles due to a lack of permissions.';
    } else if (errorAdding && errorRemoving) {
      message =
        'This role was edited to both add and remove transformations; however, this role was not added or removed from those transformations\u2019 allowed_roles due to a lack of permissions.';
    } else if (errorAdding) {
      message =
        'This role was edited to include new transformations, but this role was not added to those transformations\u2019 allowed_roles due to a lack of permissions.';
    } else {
      message =
        'This role was edited to remove transformations, but this role was not removed from those transformations\u2019 allowed_roles due to a lack of permissions.';
    }

    this.flashMessages.info(message, { sticky: true, priority: 300 });
  }

  // Removes this role from all of its transformations' allowed_roles on delete.
  async cleanupTransformationsOnDelete(roleName, backend) {
    const transformations = this.args.form.data.transformations ?? [];
    if (transformations.length === 0) return;

    await Promise.all(
      transformations.map((t) => this.syncTransformationForRole(t, roleName, backend, 'REMOVE'))
    );
  }

  transition(route = 'show') {
    this.errorMessage = '';
    this.modelValidations = null;
    const { backend, name } = this.args.form.data;
    if (route === 'list') {
      this.router.transitionTo('vault.cluster.secrets.backend.list-root', backend, {
        queryParams: { tab: 'role' },
      });
      return;
    }
    this.router.transitionTo('vault.cluster.secrets.backend.show', `role/${name}`);
  }

  @action async createOrUpdate(event) {
    event.preventDefault();

    const { isValid, state, invalidFormMessage, data } = this.args.form.toJSON();
    this.modelValidations = isValid ? null : state;
    this.errorMessage = invalidFormMessage;
    if (!isValid) return;

    const { name, transformations, backend } = data;

    const isCreate = this.args.mode === 'create';
    try {
      await this.api.secrets.transformWriteRole(name, backend, { transformations });
      this.flashMessages.success('Role saved.');
      await this.handleTransformationSync(name, backend, isCreate ? 'create' : 'update');
      this.transition();
    } catch (e) {
      const { message } = await this.api.parseError(e);
      this.errorMessage = message;
    }
  }

  @action async onDelete() {
    const { name, backend } = this.args.form.data;
    try {
      await this.api.secrets.transformDeleteRole(name, backend);
      this.flashMessages.success('Role deleted.');
      await this.cleanupTransformationsOnDelete(name, backend);
      this.transition('list');
    } catch (e) {
      const { message } = await this.api.parseError(e);
      this.flashMessages.danger(message);
    }
  }
}
