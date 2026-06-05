/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { service } from '@ember/service';
import {
  SecretsApiTransformListRolesListEnum,
  SecretsApiTransformListTemplatesListEnum,
} from '@hashicorp/vault-client-typescript';

/**
 * @module TransformationEdit
 * `TransformationEdit` is a component that allows you to create/edit or view a transformation.
 *
 * @example
 * ```js
 *   <TransformationEdit @form={{this.form}} @capabilities={{this.capabilities}} @mode={{this.mode}} />
 * ```
 * @param {object} form - TransformationForm instance with data and formFields.
 * @param {object} capabilities - Object with canDelete, canUpdate, canRead capabilities.
 * @param {string} mode - Is either show, create or edit.
 */

export default class TransformationEditComponent extends Component {
  @service flashMessages;
  @service router;
  @service api;

  @tracked errorMessage = '';
  @tracked modelValidations;
  @tracked roles = [];
  @tracked templates = [];
  @tracked isDeleteModalActive = false;
  @tracked isEditModalActive = false;

  // Non-tracked: used only to diff added/removed roles after save
  initialAllowedRoles = [];

  constructor() {
    super(...arguments);
    this.initialAllowedRoles = [...(this.args.form.data.allowed_roles ?? [])];
    this.fetchRoles();
    this.fetchTemplates();
  }

  get visibleFormFields() {
    const type = this.args.form.data.type;
    return this.args.form.formFields.filter((field) => {
      switch (field.name) {
        case 'tweak_source':
          return type === 'fpe';
        case 'masking_character':
          return type === 'masking';
        case 'template':
          return type !== 'tokenization';
        case 'mapping_mode':
        case 'convergent':
        case 'max_ttl':
        case 'stores':
          return type === 'tokenization';
        default:
          return true;
      }
    });
  }

  async fetchRoles() {
    try {
      const resp = await this.api.secrets.transformListRoles(
        this.args.form.data.backend,
        SecretsApiTransformListRolesListEnum.TRUE
      );
      this.roles = (resp.keys ?? []).map((key) => ({ id: key }));
    } catch {
      // swallow errors, SearchSelect will fall back to string-list
    }
  }

  async fetchTemplates() {
    try {
      const resp = await this.api.secrets.transformListTemplates(
        this.args.form.data.backend,
        SecretsApiTransformListTemplatesListEnum.TRUE
      );
      this.templates = (resp.keys ?? []).map((key) => ({ id: key }));
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
      },
      { label: this.title },
      { label: this.args.mode === 'create' ? 'transformation' : name },
    ];
  }

  get title() {
    if (this.args.mode === 'create') {
      return 'Create transformation';
    } else if (this.args.mode === 'edit') {
      return 'Edit transformation';
    } else {
      return 'Transformation';
    }
  }

  get subtitle() {
    if (this.args.mode === 'show') {
      return this.args.form?.data?.name;
    }
    return '';
  }

  get cliCommand() {
    const { type, allowed_roles, tweak_source, name } = this.args.form?.data ?? {};
    if (!name) return '';

    const rolesArr = allowed_roles ?? [];
    const wildCardRole = rolesArr.find((role) => role.includes('*'));
    let role = '<choose a role>';
    if (rolesArr.length === 1 && !wildCardRole) {
      role = rolesArr[0];
    }

    let tweak = '';
    if (type === 'fpe' && tweak_source === 'supplied') {
      tweak = 'tweak=<enter your tweak>';
    }

    return `${role} value=<enter your value here> ${tweak} transformation=${name}`;
  }

  isWildcard(roleName) {
    return typeof roleName === 'string' && roleName.includes('*');
  }

  async syncRoleForTransformation(roleName, transformationName, backend, syncAction) {
    if (this.isWildcard(roleName)) return;

    let currentTransformations;
    try {
      const resp = await this.api.secrets.transformReadRole(roleName, backend);
      const data = resp?.data || resp || {};
      currentTransformations = data.transformations || [];
    } catch (readErr) {
      const { status } = await this.api.parseError(readErr);
      if (status === 403) {
        this.flashMessages.info(
          `The transformation was saved, but the role "${roleName}" could not be updated due to a lack of permissions.`,
          { sticky: true, priority: 300 }
        );
        return;
      }
      // Role not found (404) or other non-403 error
      if (syncAction === 'ADD') {
        // Auto-create the role with this transformation
        try {
          await this.api.secrets.transformWriteRole(roleName, backend, {
            transformations: [transformationName],
          });
        } catch (createErr) {
          const { message } = await this.api.parseError(createErr);
          this.flashMessages.info(
            `The transformation was saved, but the role "${roleName}" could not be created: ${message}`,
            { sticky: true, priority: 300 }
          );
        }
      }
      // For REMOVE: role doesn't exist, nothing to do
      return;
    }

    let updatedTransformations;
    if (syncAction === 'ADD') {
      updatedTransformations = currentTransformations.includes(transformationName)
        ? currentTransformations
        : [...currentTransformations, transformationName];
    } else {
      updatedTransformations = currentTransformations.filter((t) => t !== transformationName);
    }

    try {
      await this.api.secrets.transformWriteRole(roleName, backend, {
        transformations: updatedTransformations,
      });
    } catch (writeErr) {
      const { status, message } = await this.api.parseError(writeErr);
      const detail = status === 403 ? `due to a lack of permissions` : message;
      this.flashMessages.info(
        `The transformation was saved, but the role "${roleName}" could not be updated: ${detail}`,
        { sticky: true, priority: 300 }
      );
    }
  }

  // Diffs current vs initial allowed_roles and syncs each changed role in parallel.
  async handleRoleSync(transformationName, backend) {
    const currentRoles = this.args.form.data.allowed_roles ?? [];
    const initialRoles = this.initialAllowedRoles;

    const addedRoles = currentRoles.filter((r) => !this.isWildcard(r) && !initialRoles.includes(r));
    const removedRoles = initialRoles.filter((r) => !this.isWildcard(r) && !currentRoles.includes(r));

    await Promise.all([
      ...addedRoles.map((r) => this.syncRoleForTransformation(r, transformationName, backend, 'ADD')),
      ...removedRoles.map((r) => this.syncRoleForTransformation(r, transformationName, backend, 'REMOVE')),
    ]);
  }

  transition(route = 'show') {
    this.errorMessage = '';
    this.modelValidations = null;
    const { name } = this.args.form.data;
    if (route === 'list') {
      const { backend } = this.args.form.data;
      this.router.transitionTo('vault.cluster.secrets.backend.list-root', backend, {
        queryParams: { tab: 'transformations' },
      });
    } else if (route === 'edit') {
      this.router.transitionTo('vault.cluster.secrets.backend.edit', name);
    } else {
      this.router.transitionTo('vault.cluster.secrets.backend.show', name);
    }
  }

  @action async createOrUpdate(event) {
    event.preventDefault();

    const { isValid, state, invalidFormMessage, data } = this.args.form.toJSON();
    this.modelValidations = isValid ? null : state;
    this.errorMessage = invalidFormMessage;
    if (!isValid) return;

    const {
      name,
      backend,
      type,
      tweak_source,
      masking_character,
      template,
      allowed_roles,
      deletion_allowed,
      mapping_mode,
      convergent,
      max_ttl,
      stores,
    } = data;

    const templateValue = Array.isArray(template) ? template[0] : template;

    try {
      await this.api.secrets.transformWriteTransformation(name, backend, {
        type,
        tweak_source,
        masking_character,
        template: templateValue,
        allowed_roles,
        deletion_allowed,
        mapping_mode,
        convergent,
        max_ttl,
        stores,
      });
      this.flashMessages.success('Transformation saved.');
      await this.handleRoleSync(name, backend);
      this.transition();
    } catch (e) {
      const { message } = await this.api.parseError(e);
      this.errorMessage = message;
    }
  }

  @action async onDelete() {
    const { name, backend } = this.args.form.data;
    try {
      await this.api.secrets.transformDeleteTransformation(name, backend);
      this.flashMessages.success('Transformation deleted.');
      this.transition('list');
    } catch (e) {
      const { message } = await this.api.parseError(e);
      this.flashMessages.danger(message);
    }
  }

  @action confirmEdit() {
    this.isEditModalActive = false;
    this.transition('edit');
  }
}
