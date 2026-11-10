/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { waitFor } from '@ember/test-waiters';
import type RoleAwsForm from 'vault/forms/aws-role';
import type ApiService from 'vault/services/api';
import type RouterService from '@ember/routing/router-service';
import type { Capabilities } from 'vault/app-types';

interface Args {
  form: RoleAwsForm;
  mode: 'create' | 'edit' | 'show';
  root: object;
  capabilities: Partial<Capabilities> & { canEdit?: boolean; canGenerate?: boolean };
}

/** Renders AWS roles and persists changes through the API service. */
export default class RoleAwsEdit extends Component<Args> {
  @service declare readonly api: ApiService;
  @service declare readonly router: RouterService;
  @tracked errorMessage = '';
  @tracked isSaving = false;

  get form() {
    return this.args.form;
  }

  get title() {
    return { create: 'Create an AWS Role', edit: 'Edit AWS Role', show: 'AWS Role' }[this.args.mode];
  }

  get subtitle() {
    return this.args.mode === 'create' ? undefined : this.form.data.name;
  }

  get breadcrumbs() {
    return [
      { label: 'Vault', text: 'Vault', icon: 'vault', path: 'vault.cluster.dashboard' },
      { text: 'Secrets engines', path: 'vault.cluster.secrets.backends' },
      this.args.root,
      { label: this.title, text: this.title },
    ];
  }

  @action
  @waitFor
  async save(event: Event) {
    event.preventDefault();
    if (this.isSaving) return;
    const { isValid, invalidFormMessage, data } = this.form.toJSON();
    if (!isValid) {
      this.errorMessage = invalidFormMessage;
      return;
    }
    const { name, backend } = this.form.data;
    if (!name || !backend) {
      this.errorMessage = 'Role name and secrets engine path are required.';
      return;
    }
    this.errorMessage = '';
    this.isSaving = true;
    try {
      await this.api.secrets.awsWriteRole(name, backend, data);
      this.router.transitionTo('vault.cluster.secrets.backend.show', name);
    } catch (error) {
      const { message } = await this.api.parseError(error);
      this.errorMessage = message;
    } finally {
      this.isSaving = false;
    }
  }

  @action
  @waitFor
  async delete() {
    const { name, backend } = this.form.data;
    if (!name || !backend) {
      this.errorMessage = 'Role name and secrets engine path are required.';
      return;
    }
    try {
      await this.api.secrets.awsDeleteRole(name, backend);
      this.router.transitionTo('vault.cluster.secrets.backend.list-root');
    } catch (error) {
      const { message } = await this.api.parseError(error);
      this.errorMessage = message;
    }
  }
}
