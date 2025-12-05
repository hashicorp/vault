/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { restartableTask } from 'ember-concurrency';

import type LdapLibraryModel from 'vault/models/ldap/library';
import type SecretsEngineResource from 'vault/resources/secrets/engine';
import type RouterService from '@ember/routing/router-service';
import type Store from '@ember-data/store';
import type { Breadcrumb } from 'vault/vault/app-types';
import LdapRoleModel from 'vault/models/ldap/role';
import { LdapLibraryAccountStatus } from 'vault/vault/adapters/ldap/library';

interface Args {
  roles: Array<LdapRoleModel>;
  promptConfig: boolean;
  secretsEngine: SecretsEngineResource;
  breadcrumbs: Array<Breadcrumb>;
}

interface Option {
  id: string;
  name: string;
  type: string;
}

export default class LdapLibrariesPageComponent extends Component<Args> {
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly store: Store;

  @tracked selectedRole: LdapRoleModel | undefined;
  @tracked librariesStatus: Array<LdapLibraryAccountStatus> = [];
  @tracked allLibraries: Array<LdapLibraryModel> = [];
  @tracked librariesError: string | null = null;

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    this.fetchLibraries.perform();
    this.fetchLibrariesStatus.perform();
  }

  get roleOptions() {
    const options = this.args.roles
      // hierarchical roles are not selectable
      .filter((r: LdapRoleModel) => !r.name.endsWith('/'))
      // *hack alert* - type is set as id so it renders beside name in search select
      // this is to avoid more changes to search select and is okay here because
      // we use the type and name to select the item below, not the id
      .map((r: LdapRoleModel) => ({ id: r.type, name: r.name, type: r.type }));
    return options;
  }

  @action
  async selectRole([option]: Array<Option>) {
    if (option) {
      const { name, type } = option;
      const model = this.args.roles.find((role) => role.name === name && role.type === type);
      this.selectedRole = model;
    }
  }

  @action
  generateCredentials() {
    const { type, name } = this.selectedRole as LdapRoleModel;
    this.router.transitionTo('vault.cluster.secrets.backend.ldap.roles.role.credentials', type, name);
  }

  fetchLibraries = restartableTask(async () => {
    const backend = this.args.secretsEngine.id;
    const allLibraries: Array<LdapLibraryModel> = [];

    try {
      this.librariesError = null; // Clear any previous errors
      await this.discoverAllLibrariesRecursively(backend, '', allLibraries);
      this.allLibraries = allLibraries;
    } catch (error) {
      // Hierarchical discovery failed - display inline error
      this.librariesError = 'Unable to load complete library information. Please try refreshing the page.';
      this.allLibraries = [];
    }
  });

  fetchLibrariesStatus = restartableTask(async () => {
    // Wait for fetchLibraries task to complete before proceeding
    await this.fetchLibraries.last;

    const allStatuses: Array<LdapLibraryAccountStatus> = [];

    for (const library of this.allLibraries) {
      try {
        const statuses = await library.fetchStatus();
        allStatuses.push(...statuses);
      } catch (error) {
        // suppressing error
      }
    }

    this.librariesStatus = allStatuses;
  });

  private async discoverAllLibrariesRecursively(
    backend: string,
    currentPath: string,
    allLibraries: Array<LdapLibraryModel>
  ): Promise<void> {
    const queryParams: { backend: string; path_to_library?: string } = { backend };
    if (currentPath) {
      queryParams.path_to_library = currentPath;
    }

    const items = await this.store.query('ldap/library', queryParams);
    const libraryItems = items.toArray() as LdapLibraryModel[];

    for (const item of libraryItems) {
      if (item.name.endsWith('/')) {
        // This is a directory - recursively explore it
        const nextPath = currentPath ? `${currentPath}${item.name}` : item.name;
        await this.discoverAllLibrariesRecursively(backend, nextPath, allLibraries);
      } else {
        // This is an actual library
        allLibraries.push(item);
      }
    }
  }
}
