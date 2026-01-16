/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { getOwner } from '@ember/owner';

import type { LdapLibrary } from 'vault/secrets/ldap';
import type FlashMessageService from 'vault/services/flash-messages';
import type { Breadcrumb, CapabilitiesMap, EngineOwner } from 'vault/vault/app-types';
import type RouterService from '@ember/routing/router-service';
import type SecretsEngineResource from 'vault/resources/secrets/engine';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type ApiService from 'vault/services/api';

interface Args {
  libraries: Array<LdapLibrary>;
  capabilities: CapabilitiesMap;
  promptConfig: boolean;
  secretsEngine: SecretsEngineResource;
  breadcrumbs: Array<Breadcrumb>;
}

export default class LdapLibrariesPageComponent extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly secretMountPath: SecretMountPath;
  @service declare readonly api: ApiService;

  @tracked filterValue = '';
  @tracked libraryToDelete: LdapLibrary | null = null;

  isHierarchical = (name: string) => name.endsWith('/');

  linkParams = (library: LdapLibrary) => {
    const route = this.isHierarchical(library.name) ? 'libraries.subdirectory' : 'libraries.library.details';
    return [route, library.completeLibraryName];
  };

  getEncodedLibraryName = (library: LdapLibrary) => {
    return library.completeLibraryName;
  };

  get mountPoint(): string {
    const owner = getOwner(this) as EngineOwner;
    return owner.mountPoint;
  }

  get filteredLibraries() {
    const { libraries } = this.args;
    return this.filterValue
      ? libraries.filter((library) => library.name.toLowerCase().includes(this.filterValue.toLowerCase()))
      : libraries;
  }

  @action
  async onDelete(library: LdapLibrary) {
    try {
      const { completeLibraryName } = library;
      await this.api.secrets.ldapLibraryDelete(completeLibraryName, this.secretMountPath.currentPath);
      this.router.transitionTo('vault.cluster.secrets.backend.ldap.libraries');
      this.flashMessages.success(`Successfully deleted library ${completeLibraryName}.`);
    } catch (error) {
      const { message } = await this.api.parseError(error);
      this.flashMessages.danger(`Error deleting library \n ${message}`);
    }
  }
}
