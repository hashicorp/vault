/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { restartableTask } from 'ember-concurrency';
import {
  LdapLibraryListListEnum,
  LdapLibraryListLibraryPathListEnum,
} from '@hashicorp/vault-client-typescript';

import type {
  LdapRole,
  LdapLibrary,
  LdapLibraryAccountStatus,
  LdapLibraryAccountStatusResponse,
} from 'vault/secrets/ldap';
import type SecretsEngineResource from 'vault/resources/secrets/engine';
import type RouterService from '@ember/routing/router-service';
import type { Breadcrumb, CapabilitiesMap } from 'vault/vault/app-types';
import type CapabilitiesService from 'vault/services/capabilities';
import type ApiService from 'vault/services/api';
import type SecretMountPath from 'vault/services/secret-mount-path';

interface Args {
  roles: Array<LdapRole>;
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
  @service declare readonly api: ApiService;
  @service declare readonly capabilities: CapabilitiesService;
  @service declare readonly secretMountPath: SecretMountPath;

  @tracked selectedRole: LdapRole | undefined;
  @tracked librariesStatus: Array<LdapLibraryAccountStatus> = [];
  @tracked allLibraries: Array<LdapLibrary> = [];
  @tracked librariesError: string | null = null;
  @tracked declare checkInCapabilities: CapabilitiesMap;

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    this.fetchLibraries.perform();
    this.fetchLibrariesStatus.perform();
  }

  get roleOptions() {
    const options = this.args.roles
      // hierarchical roles are not selectable
      .filter((r: LdapRole) => !r.name.endsWith('/'))
      // *hack alert* - type is set as id so it renders beside name in search select
      // this is to avoid more changes to search select and is okay here because
      // we use the type and name to select the item below, not the id
      .map((r: LdapRole) => ({ id: r.type, name: r.name, type: r.type }));
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
    const { type, name } = this.selectedRole as LdapRole;
    this.router.transitionTo('vault.cluster.secrets.backend.ldap.roles.role.credentials', type, name);
  }

  fetchLibraries = restartableTask(async () => {
    const backend = this.args.secretsEngine.id;
    const allLibraries: Array<LdapLibrary> = [];

    try {
      this.librariesError = null; // Clear any previous errors
      await this.discoverAllLibrariesRecursively(backend, '', allLibraries);
      this.allLibraries = allLibraries;
      // fetch capabilities for all libraries
      const paths = this.allLibraries.map(({ completeLibraryName: name }) =>
        this.capabilities.pathFor('ldapLibraryCheckIn', { backend, name })
      );
      this.checkInCapabilities = await this.capabilities.fetch(paths);
    } catch (error) {
      const { status } = await this.api.parseError(error);
      // Only set error message if status is not a 404 which just means an empty response.
      if (status !== 404) {
        // Hierarchical discovery failed - display inline error
        this.librariesError = 'Unable to load complete library information. Please try refreshing the page.';
      }
      this.allLibraries = [];
    }
  });

  fetchLibrariesStatus = restartableTask(async () => {
    // Wait for fetchLibraries task to complete before proceeding
    await this.fetchLibraries.last;

    const allStatuses: Array<LdapLibraryAccountStatus> = [];
    const requests = this.allLibraries.map((library) =>
      this.api.secrets.ldapLibraryCheckStatus(library.completeLibraryName, this.secretMountPath.currentPath)
    );
    const results = await Promise.allSettled(requests);

    for (const result of results) {
      // ignore failures and only extract statuses from successful requests
      if (result.status === 'fulfilled') {
        const index = results.indexOf(result);
        const response = result.value.data as LdapLibraryAccountStatusResponse;

        for (const key in response) {
          const status = response[key] as LdapLibraryAccountStatusResponse[string];
          allStatuses.push({
            ...status,
            account: key,
            library: this.allLibraries[index]?.completeLibraryName as string,
          });
        }
      }
    }

    this.librariesStatus = allStatuses;
  });

  private async discoverAllLibrariesRecursively(
    backend: string,
    pathToLibrary: string,
    allLibraries: Array<LdapLibrary>
  ): Promise<void> {
    const { currentPath } = this.secretMountPath;
    const { keys } = pathToLibrary
      ? await this.api.secrets.ldapLibraryListLibraryPath(
          pathToLibrary,
          currentPath,
          LdapLibraryListLibraryPathListEnum.TRUE
        )
      : await this.api.secrets.ldapLibraryList(currentPath, LdapLibraryListListEnum.TRUE);

    const libraries =
      keys?.map((name) => {
        // if path is provided combine with name for completeLibraryName
        const completeLibraryName = pathToLibrary ? `${pathToLibrary}${name}` : name;
        return { name, completeLibraryName } as LdapLibrary;
      }) || [];

    for (const library of libraries) {
      if (library.name.endsWith('/')) {
        // This is a directory - recursively explore it
        await this.discoverAllLibrariesRecursively(backend, library.completeLibraryName, allLibraries);
      } else {
        // This is an actual library
        allLibraries.push(library);
      }
    }
  }
}
