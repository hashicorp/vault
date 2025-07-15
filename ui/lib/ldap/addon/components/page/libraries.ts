/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { getOwner } from '@ember/owner';
import errorMessage from 'vault/utils/error-message';

import type LdapLibraryModel from 'vault/models/ldap/library';
import type SecretEngineModel from 'vault/models/secret-engine';
import type FlashMessageService from 'vault/services/flash-messages';
import type { Breadcrumb, EngineOwner } from 'vault/vault/app-types';
import type RouterService from '@ember/routing/router-service';

interface Args {
  libraries: Array<LdapLibraryModel>;
  promptConfig: boolean;
  backendModel: SecretEngineModel;
  breadcrumbs: Array<Breadcrumb>;
}

export default class LdapLibrariesPageComponent extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;
  @service('app-router') declare readonly router: RouterService;

  @tracked filterValue = '';
  @tracked libraryToDelete: LdapLibraryModel | null = null;

  isHierarchical = (name: string) => name.endsWith('/');

  linkParams = (library: LdapLibraryModel) => {
    const route = this.isHierarchical(library.name) ? 'libraries.subdirectory' : 'libraries.library.details';
    return [route, library.completeLibraryName];
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
  async onDelete(model: LdapLibraryModel) {
    try {
      const message = `Successfully deleted library ${model.completeLibraryName}.`;
      await model.destroyRecord();
      this.router.transitionTo('vault.cluster.secrets.backend.ldap.libraries');
      this.flashMessages.success(message);
    } catch (error) {
      this.flashMessages.danger(`Error deleting library \n ${errorMessage(error)}`);
    }
  }
}
