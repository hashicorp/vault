/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

import type RouterService from '@ember/routing/router-service';
import type { Breadcrumb, EngineOwner } from 'vault/vault/app-types';
import type SecretsEngineResource from 'vault/resources/secrets/engine';

/**
 * @module Overview
 * OverviewPage component is a child component to overview kubernetes secrets engine.
 *
 * @param {boolean} promptConfig - whether or not to display config cta
 * @param {object} secretsEngine -  SecretsEngine resource that contains kubernetes configuration
 * @param {array} roles - array of roles
 * @param {array} breadcrumbs - breadcrumbs as an array of objects that contain label and route
 */

interface Args {
  promptConfig: boolean;
  secretsEngine: SecretsEngineResource;
  roles: string[];
  breadcrumbs: Array<Breadcrumb>;
}

export default class OverviewPageComponent extends Component<Args> {
  @service('app-router') declare readonly router: RouterService;

  @tracked selectedRole = '';
  @tracked roleOptions: Array<{ name: string; id: string }> = [];

  constructor(owner: EngineOwner, args: Args) {
    super(owner, args);
    this.roleOptions = this.args.roles.map((role) => {
      return { name: role, id: role };
    });
  }

  @action
  selectRole([roleName]: [string]) {
    this.selectedRole = roleName;
  }

  @action
  generateCredential() {
    this.router.transitionTo(
      'vault.cluster.secrets.backend.kubernetes.roles.role.credentials',
      this.selectedRole
    );
  }
}
