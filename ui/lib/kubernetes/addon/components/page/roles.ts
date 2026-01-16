/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { getOwner } from '@ember/owner';
import { tracked } from '@glimmer/tracking';
import keys from 'core/utils/keys';

import type FlashMessageService from 'ember-cli-flash/services/flash-messages';
import type RouterService from '@ember/routing/router-service';
import type ApiService from 'vault/services/api';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type SecretsEngineResource from 'vault/resources/secrets/engine';
import type { Breadcrumb, EngineOwner } from 'vault/app-types';
import { HTMLElementEvent } from 'vault/forms';

/**
 * @module Roles
 * RolesPage component is a child component to show list of roles.
 * It also handles the filtering actions of roles.
 *
 * @param {array} roles - array of roles
 * @param {boolean} promptConfig - whether or not to display config cta
 * @param {string} filterValue - value of queryParam pageFilter
 * @param {array} breadcrumbs - breadcrumbs as an array of objects that contain label and route
 */

interface Args {
  roles: Array<string>;
  promptConfig: boolean;
  secretsEngine: SecretsEngineResource;
  filterValue: string;
  breadcrumbs: Array<Breadcrumb>;
}

export default class RolesPageComponent extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPath;

  @tracked query;
  @tracked roleToDelete = null;

  constructor(owner: EngineOwner, args: Args) {
    super(owner, args);
    this.query = this.args.filterValue;
  }

  get mountPoint() {
    return (getOwner(this) as EngineOwner).mountPoint;
  }

  navigate(pageFilter?: string) {
    const route = `${this.mountPoint}.roles.index`;
    const args = [route, { queryParams: { pageFilter: pageFilter || null } }];
    this.router.transitionTo(...args);
  }

  @action
  handleKeyDown(event: KeyboardEvent) {
    const isEscKeyPressed = keys.ESC.includes(event.key);
    if (isEscKeyPressed) {
      // On escape, transition to roles index route.
      this.navigate();
    }
    // ignore all other key events
  }

  @action handleInput(evt: HTMLElementEvent<HTMLInputElement>) {
    this.query = evt.target.value;
  }

  @action
  handleSearch(evt: HTMLElementEvent<HTMLInputElement>) {
    evt.preventDefault();
    this.navigate(this.query);
  }

  @action
  async onDelete(role: string) {
    try {
      await this.api.secrets.kubernetesDeleteRole(role, this.secretMountPath.currentPath);
      this.flashMessages.success(`Successfully deleted role ${role}`);
      this.router.refresh(`${this.mountPoint}.roles.index`);
    } catch (error) {
      const { message } = await this.api.parseError(
        error,
        'Error deleting role. Please try again or contact support'
      );
      this.flashMessages.danger(message);
    } finally {
      this.roleToDelete = null;
    }
  }
}
