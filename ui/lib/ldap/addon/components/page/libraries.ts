/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import { getOwner } from '@ember/application';
import errorMessage from 'vault/utils/error-message';

import type LdapLibraryModel from 'vault/models/ldap/library';
import type SecretEngineModel from 'vault/models/secret-engine';
import type FlashMessageService from 'vault/services/flash-messages';
import type { Breadcrumb, EngineOwner } from 'vault/vault/app-types';

interface Args {
  libraries: Array<LdapLibraryModel>;
  promptConfig: boolean;
  backendModel: SecretEngineModel;
  breadcrumbs: Array<Breadcrumb>;
}

export default class LdapLibrariesPageComponent extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;

  @tracked filterValue = '';

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
      const message = `Successfully deleted library ${model.name}.`;
      await model.destroyRecord();
      this.args.libraries.removeObject(model);
      this.flashMessages.success(message);
    } catch (error) {
      this.flashMessages.danger(`Error deleting library \n ${errorMessage(error)}`);
    }
  }
}
