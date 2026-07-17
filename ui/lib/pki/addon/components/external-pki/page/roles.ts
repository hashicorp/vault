/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import Component from '@glimmer/component';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { SetupSteps } from '../implementation-select';

import type { RolesIndexRouteModel } from 'pki/routes/external/roles';
import type RouterService from '@ember/routing/router-service';
import type { HTMLElementEvent } from 'vault/forms';

interface Args {
  model: RolesIndexRouteModel;
}

export default class ExternalPkiPageRolesComponent extends Component<Args> {
  @service('app-router') declare readonly router: RouterService;

  @tracked searchInput = '';

  roleConfig = SetupSteps.ROLE_CONFIG;
  tableColumns = [
    {
      key: 'name',
      label: 'Role name',
      isSortable: true,
      customTableItem: true,
    },
  ];

  get rolesList() {
    const filteredRoles = this.args.model.roles.filter((n) => n.includes(this.searchInput));
    return filteredRoles?.map((r) => ({ name: r }));
  }

  @action
  handleSearch(e: HTMLElementEvent<HTMLInputElement>) {
    this.searchInput = e.target.value;
  }

  @action
  refresh() {
    this.router.refresh(this.router.currentRoute?.parent?.name);
  }
}
